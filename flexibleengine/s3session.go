package flexibleengine

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
)

var s3Session *session.Session
var s3Mutex = new(sync.Mutex)

type awsLogger struct{}

func (l awsLogger) Log(args ...interface{}) {
	tokens := make([]string, 0, len(args))
	for _, arg := range args {
		if token, ok := arg.(string); ok {
			tokens = append(tokens, token)
		}
	}
	log.Printf("[DEBUG] [aws-sdk-go] %s", strings.Join(tokens, " "))
}

func computeS3conn(c *Config, region string) (*s3.S3, error) {
	s3Mutex.Lock()
	defer s3Mutex.Unlock()

	var err error
	if s3Session == nil {
		log.Printf("[DEBUG] initialize Swift S3 session")
		s3Session, err = newS3Session(c, logging.IsDebugOrHigher())
		if err != nil {
			return nil, errwrap.Wrapf("Error creating Swift S3 session: {{err}}", err)
		}
	}

	endpoint := getOssEndpoint(c, region)
	awsS3Sess := s3Session.Copy(&aws.Config{Endpoint: aws.String(endpoint)})
	s3conn := s3.New(awsS3Sess)

	return s3conn, nil
}

func newS3Session(c *Config, osDebug bool) (*session.Session, error) {
	if c.AccessKey == "" || c.SecretKey == "" {
		return nil, fmt.Errorf("missing credentials for Swift S3 Provider, need access_key and secret_key values for provider")
	}

	// Setup S3 client/config information for Swift S3 buckets
	log.Println("[INFO] Building Swift S3 auth structure")
	creds, err := GetCredentials(c)
	if err != nil {
		return nil, err
	}
	// Call Get to check for credential provider. If nothing found, we'll get an
	// error, and we can present it nicely to the user
	cp, err := creds.Get()
	if err != nil {
		if sErr, ok := err.(awserr.Error); ok && sErr.Code() == "NoCredentialProviders" {
			return nil, fmt.Errorf("No valid credential sources found for S3 Provider")
		}

		return nil, fmt.Errorf("Error loading credentials for S3 Provider: %s", err)
	}

	log.Printf("[INFO] S3 Auth provider used: %q", cp.ProviderName)

	sConfig := &aws.Config{
		Credentials: creds,
		Region:      aws.String(c.Region),
		HTTPClient:  cleanhttp.DefaultClient(),
	}

	if osDebug {
		sConfig.LogLevel = aws.LogLevel(aws.LogDebugWithHTTPBody | aws.LogDebugWithRequestRetries | aws.LogDebugWithRequestErrors)
		sConfig.Logger = awsLogger{}
	}

	if c.Insecure {
		transport := sConfig.HTTPClient.Transport.(*http.Transport)
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	// Set up base session for S3
	return session.NewSession(sConfig)
}

func getOssEndpoint(c *Config, region string) string {
	if endpoint, ok := c.Endpoints["oss"]; ok {
		return endpoint
	}
	return fmt.Sprintf("https://oss.%s.%s/", region, c.Cloud)
}
