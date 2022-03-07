package flexibleengine

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/chnsz/golangsdk"
	huaweisdk "github.com/chnsz/golangsdk/openstack"
	"github.com/chnsz/golangsdk/openstack/identity/v3/domains"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	huaweiconfig "github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/pathorcontents"
)

type Config struct {
	huaweiconfig.Config
}

func (c *Config) LoadAndValidate() error {
	if c.MaxRetries < 0 {
		return fmt.Errorf("max_retries should be a positive value")
	}

	err := fmt.Errorf("Must config token or aksk or username password to be authorized")

	if c.Token != "" {
		err = buildClientByToken(c)
	} else if c.Password != "" {
		if c.Username == "" && c.UserID == "" {
			err = fmt.Errorf("\"password\": one of `user_name, user_id` must be specified")
		} else {
			err = buildClientByPassword(c)
		}
	} else if c.AccessKey != "" && c.SecretKey != "" {
		err = buildClientByAKSK(c)
	}

	if err != nil {
		return err
	}

	if c.HwClient != nil && c.HwClient.ProjectID != "" {
		c.RegionProjectIDMap[c.Region] = c.HwClient.ProjectID
	}

	// set DomainID for IAM resource
	if c.DomainID == "" {
		if domainID, err := c.getDomainID(); err == nil {
			c.DomainID = domainID
		} else {
			log.Printf("[WARN] get domain id failed: %s", err)
		}
	}

	return nil
}

func generateTLSConfig(c *Config) (*tls.Config, error) {
	config := &tls.Config{}
	if c.CACertFile != "" {
		caCert, _, err := pathorcontents.Read(c.CACertFile)
		if err != nil {
			return nil, fmt.Errorf("Error reading CA Cert: %s", err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM([]byte(caCert))
		config.RootCAs = caCertPool
	}

	if c.Insecure {
		config.InsecureSkipVerify = true
	}

	if c.ClientCertFile != "" && c.ClientKeyFile != "" {
		clientCert, _, err := pathorcontents.Read(c.ClientCertFile)
		if err != nil {
			return nil, fmt.Errorf("Error reading Client Cert: %s", err)
		}
		clientKey, _, err := pathorcontents.Read(c.ClientKeyFile)
		if err != nil {
			return nil, fmt.Errorf("Error reading Client Key: %s", err)
		}

		cert, err := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
		if err != nil {
			return nil, err
		}

		config.Certificates = []tls.Certificate{cert}
		config.BuildNameToCertificate()
	}

	return config, nil
}

func retryBackoffFunc(ctx context.Context, respErr *golangsdk.ErrUnexpectedResponseCode, e error, retries uint) error {
	minutes := int(math.Pow(2, float64(retries)))
	if minutes > 30 { // won't wait more than 30 minutes
		minutes = 30
	}

	log.Printf("[WARN] Received StatusTooManyRequests response code, try to sleep %d minutes", minutes)
	sleep := time.Duration(minutes) * time.Minute

	if ctx != nil {
		select {
		case <-time.After(sleep):
		case <-ctx.Done():
			return e
		}
	} else {
		time.Sleep(sleep)
	}

	return nil
}

func genClient(c *Config, ao golangsdk.AuthOptionsProvider) (*golangsdk.ProviderClient, error) {
	client, err := huaweisdk.NewClient(ao.GetIdentityEndpoint())
	if err != nil {
		return nil, err
	}

	// Set UserAgent
	client.UserAgent.Prepend("terraform-provider-flexibleengine")

	config, err := generateTLSConfig(c)
	if err != nil {
		return nil, err
	}
	transport := &http.Transport{Proxy: http.ProxyFromEnvironment, TLSClientConfig: config}

	client.HTTPClient = http.Client{
		Transport: &huaweiconfig.LogRoundTripper{
			Rt:         transport,
			OsDebug:    logging.IsDebugOrHigher(),
			MaxRetries: c.MaxRetries,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if client.AKSKAuthOptions.AccessKey != "" {
				golangsdk.ReSign(req, golangsdk.SignOptions{
					AccessKey: client.AKSKAuthOptions.AccessKey,
					SecretKey: client.AKSKAuthOptions.SecretKey,
				})
			}
			return nil
		},
	}

	if c.MaxRetries > 0 {
		client.MaxBackoffRetries = uint(c.MaxRetries)
		client.RetryBackoffFunc = retryBackoffFunc
	}

	// Validate authentication normally.
	err = huaweisdk.Authenticate(client, ao)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func buildClientByToken(c *Config) error {
	var pao, dao golangsdk.AuthOptions

	pao = golangsdk.AuthOptions{
		DomainID:   c.DomainID,
		DomainName: c.DomainName,
		TenantID:   c.TenantID,
		TenantName: c.TenantName,
	}

	dao = golangsdk.AuthOptions{
		DomainID:   c.DomainID,
		DomainName: c.DomainName,
	}

	for _, ao := range []*golangsdk.AuthOptions{&pao, &dao} {
		ao.IdentityEndpoint = c.IdentityEndpoint
		ao.TokenID = c.Token

	}
	return genClients(c, pao, dao)
}

func buildClientByAKSK(c *Config) error {
	var pao, dao golangsdk.AKSKAuthOptions

	pao = golangsdk.AKSKAuthOptions{
		ProjectName: c.TenantName,
		ProjectId:   c.TenantID,
	}

	dao = golangsdk.AKSKAuthOptions{
		DomainID: c.DomainID,
		Domain:   c.DomainName,
	}

	for _, ao := range []*golangsdk.AKSKAuthOptions{&pao, &dao} {
		ao.IdentityEndpoint = c.IdentityEndpoint
		ao.AccessKey = c.AccessKey
		ao.SecretKey = c.SecretKey
	}
	return genClients(c, pao, dao)
}

func buildClientByPassword(c *Config) error {
	var pao, dao golangsdk.AuthOptions

	pao = golangsdk.AuthOptions{
		DomainID:   c.DomainID,
		DomainName: c.DomainName,
		TenantID:   c.TenantID,
		TenantName: c.TenantName,
	}

	dao = golangsdk.AuthOptions{
		DomainID:   c.DomainID,
		DomainName: c.DomainName,
	}

	for _, ao := range []*golangsdk.AuthOptions{&pao, &dao} {
		ao.IdentityEndpoint = c.IdentityEndpoint
		ao.Password = c.Password
		ao.Username = c.Username
		ao.UserID = c.UserID
	}
	return genClients(c, pao, dao)
}

func genClients(c *Config, pao, dao golangsdk.AuthOptionsProvider) error {
	client, err := genClient(c, pao)
	if err != nil {
		return err
	}
	c.HwClient = client

	client, err = genClient(c, dao)
	if err == nil {
		c.DomainClient = client
	}
	return err
}

func (c *Config) determineRegion(region string) string {
	// If a resource-level region was not specified, and a provider-level region was set,
	// use the provider-level region.
	if region == "" && c.Region != "" {
		region = c.Region
	}

	log.Printf("[DEBUG] FlexibleEngine Region is: %s", region)
	return region
}

func (c *Config) getDomainID() (string, error) {
	identityClient, err := c.identityV3Client(c.Region)
	if err != nil {
		return "", fmt.Errorf("Error creating FlexibleEngine identity client: %s", err)
	}

	identityClient.ResourceBase = identityClient.Endpoint + "auth/"

	// the List request does not support query options
	allPages, err := domains.List(identityClient, nil).AllPages()
	if err != nil {
		return "", fmt.Errorf("List domains failed, err=%s", err)
	}

	all, err := domains.ExtractDomains(allPages)
	if err != nil {
		return "", fmt.Errorf("Extract domains failed, err=%s", err)
	}

	if len(all) == 0 {
		return "", fmt.Errorf("domain was not found")
	}

	if c.DomainName != "" && c.DomainName != all[0].Name {
		return "", fmt.Errorf("domain %s was not found, got %s", c.DomainName, all[0].Name)
	}

	return all[0].ID, nil
}

// client for ecs v1
func (c *Config) computeV1Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewComputeV1(c.HwClient, golangsdk.EndpointOpts{
		Region:       region,
		Availability: c.getHwEndpointType(),
	})
}

// client for nova v2 and bms Services i.e. flavor, nic, keypair.
func (c *Config) computeV2Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewComputeV2(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

//bmsClient used to access the v2.1 bms Services i.e. servers, tags.
func (c *Config) bmsClient(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewBMSV2(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

// dnsV2Client is a global client for DNS service, the endpoint is https://dns.{cloud}/v2/
func (c *Config) dnsV2Client(_ string) (*golangsdk.ServiceClient, error) {
	sc := new(golangsdk.ServiceClient)
	sc.ProviderClient = c.HwClient
	sc.Endpoint = fmt.Sprintf("https://dns.%s", defaultCloud)
	sc.ResourceBase = fmt.Sprintf("%s/v2/", sc.Endpoint)

	return sc, nil
}

func (c *Config) identityV3Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewIdentityV3(c.DomainClient, golangsdk.EndpointOpts{
		//Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) imageV2Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewImageServiceV2(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) networkingV1Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewNetworkV1(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) networkingV2Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewNetworkV2(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) otcV1Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewElbV1(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	}, "elb")
}

func (c *Config) elbV2Client(region string) (*golangsdk.ServiceClient, error) {
	sc, err := c.sdkClient(region, "network")
	if err == nil {
		sc.Endpoint = strings.Replace(sc.Endpoint, "vpc", "elb", 1)
		sc.ResourceBase = sc.Endpoint + fmt.Sprintf("v2.0/%s/", c.HwClient.ProjectID)
	}

	return sc, err
}

func (c *Config) vpcepV1Client(region string) (*golangsdk.ServiceClient, error) {
	sc, err := c.sdkClient(region, "network")
	if err == nil {
		sc.Endpoint = strings.Replace(sc.Endpoint, "vpc", "vpcep", 1)
		sc.ResourceBase = sc.Endpoint + fmt.Sprintf("v1/%s/", c.HwClient.ProjectID)
	}

	return sc, err
}

func (c *Config) autoscalingV1Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewAutoScalingService(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) SmnV2Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewSMNV2(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) MlsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewMLSV1(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) MrsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewMapReduceV1(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) natV2Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewNatV2(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) drsV2Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewDRSServiceV2(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

// sfsV1Client used to sfs-turbo resource
func (c *Config) sfsV1Client(region string) (*golangsdk.ServiceClient, error) {
	sc, err := huaweisdk.NewSharedFileSystemV2(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})

	if err == nil {
		sc.Endpoint = strings.Replace(sc.Endpoint, "sfs", "sfs-turbo", 1)
		sc.Endpoint = strings.Replace(sc.Endpoint, "/v2/", "/v1/", 1)
		sc.ResourceBase = sc.Endpoint
	}
	return sc, err
}

func (c *Config) orchestrationV1Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewOrchestrationV1(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) dwsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewDWSClient(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) csbsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewCSBSService(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) vbsV2Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewVBSV2(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) antiddosV1Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewAntiDDoSV1(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) ctsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewCTSService(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) dcsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewDCSServiceV1(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) cceV3Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewCCEV3(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) kmsKeyV1Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewKMSV1(c.HwClient, golangsdk.EndpointOpts{
		Region:       region,
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) rdsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewRDSV1(c.HwClient, golangsdk.EndpointOpts{
		Region:       region,
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) rdsV3Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewRDSV3(c.HwClient, golangsdk.EndpointOpts{
		Region:       region,
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) ddsV3Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewDDSV3(c.HwClient, golangsdk.EndpointOpts{
		Region:       region,
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) sdrsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewSDRSV1(c.HwClient, golangsdk.EndpointOpts{
		Region:       region,
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) sdkClient(region, serviceType string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewSDKClient(
		c.HwClient,
		golangsdk.EndpointOpts{
			Region:       c.determineRegion(region),
			Availability: c.getHwEndpointType(),
		},
		serviceType)
}

func (c *Config) getHwEndpointType() golangsdk.Availability {
	return golangsdk.AvailabilityPublic
}
