package flexibleengine

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/golangsdk"
)

const (
	greenCode  = "\033[0m\033[1;32m"
	yellowCode = "\033[0m\033[1;33m"
	resetCode  = "\033[0m\033[1;31m"
)

func green(str interface{}) string {
	return fmt.Sprintf("%s%#v%s", greenCode, str, resetCode)
}

func yellow(str interface{}) string {
	return fmt.Sprintf("%s%#v%s", yellowCode, str, resetCode)
}

func testAccPreCheckServiceEndpoints(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("TF_ACC not set, skipping FlexibleEngine service endpoints test.")
	}

	projectID := os.Getenv("OS_PROJECT_ID")
	if projectID == "" {
		t.Fatalf(yellow("OS_PROJECT_ID must be set for service endpoint acceptance test"))
	}
}

func testCheckServiceURL(t *testing.T, expected, actual, service string) {
	if actual != expected {
		t.Fatalf("%s endpoint: expected %s but got %s", service, green(expected), yellow(actual))
	}
	t.Logf("%s endpoint:\t %s", service, actual)
}

func TestAccServiceEndpoints_Global(t *testing.T) {
	testAccPreCheckServiceEndpoints(t)

	testProvider := Provider()
	raw := make(map[string]interface{})
	diags := testProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("Unexpected error when configure FlexibleEngine provider: %s", diags[0].Summary)
	}

	var expectedURL, actualURL string
	var serviceClient *golangsdk.ServiceClient
	var err error
	config := testProvider.Meta().(*Config)

	// test the endpoint of identity service
	serviceClient, err = config.identityV3Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating FlexibleEngine identity client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://iam.%s.%s/v3/", OS_REGION_NAME, defaultCloud)
	actualURL = serviceClient.ResourceBaseURL()
	testCheckServiceURL(t, expectedURL, actualURL, "Identity v3")
}

func TestAccServiceEndpoints_Management(t *testing.T) {
	testAccPreCheckServiceEndpoints(t)

	testProvider := Provider()
	raw := make(map[string]interface{})
	diags := testProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("Unexpected error when configure FlexibleEngine provider: %s", diags[0].Summary)
	}

	var expectedURL, actualURL string
	var serviceClient *golangsdk.ServiceClient
	var err error
	config := testProvider.Meta().(*Config)

	// test the endpoint of CTS service
	serviceClient, err = config.ctsV1Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating FlexibleEngine CTS client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://cts.%s.%s/v1.0/%s/", OS_REGION_NAME, defaultCloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	testCheckServiceURL(t, expectedURL, actualURL, "CTS")

	// test the endpoint of CES service
	serviceClient, err = config.CesV1Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating FlexibleEngine CES client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://ces.%s.%s/V1.0/%s/", OS_REGION_NAME, defaultCloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	testCheckServiceURL(t, expectedURL, actualURL, "CES")
}

// TestAccServiceEndpoints_Compute test for endpoints of the clients used in ecs
// include computeV1Client,computeV2Client,bmsClient,autoscalingV1Client,imageV2Client,
// cceV3Client
func TestAccServiceEndpoints_Compute(t *testing.T) {
	testAccPreCheckServiceEndpoints(t)

	testProvider := Provider()
	raw := make(map[string]interface{})
	diags := testProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("Unexpected error when configure FlexibleEngine provider: %s", diags[0].Summary)
	}

	config := testProvider.Meta().(*Config)
	var expectedURL, actualURL string
	var serviceClient *golangsdk.ServiceClient
	var err error

	// test for computeV1Client
	serviceClient, err = config.computeV1Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating FlexibleEngine ecs v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://ecs.%s.%s/v1/%s/", OS_REGION_NAME, defaultCloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	testCheckServiceURL(t, expectedURL, actualURL, "ECS v1")

	// test for computeV2Client
	serviceClient, err = nil, nil
	serviceClient, err = config.computeV2Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating FlexibleEngine ecs v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://ecs.%s.%s/v2/%s/", OS_REGION_NAME, defaultCloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	testCheckServiceURL(t, expectedURL, actualURL, "ecs v2")

	// test for bmsClient
	serviceClient, err = nil, nil
	serviceClient, err = config.bmsClient(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating FlexibleEngine bms v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://ecs.%s.%s/v2.1/%s/", OS_REGION_NAME, defaultCloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	testCheckServiceURL(t, expectedURL, actualURL, "bms v2.1")

	// test for autoscalingV1Client
	serviceClient, err = nil, nil
	serviceClient, err = config.autoscalingV1Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating FlexibleEngine autoscaling v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://as.%s.%s/autoscaling-api/v1/%s/", OS_REGION_NAME, defaultCloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	testCheckServiceURL(t, expectedURL, actualURL, "autoscaling v1")

	// test for imageV2Client
	serviceClient, err = nil, nil
	serviceClient, err = config.imageV2Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating FlexibleEngine image v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://ims.%s.%s/v2/", OS_REGION_NAME, defaultCloud)
	actualURL = serviceClient.ResourceBaseURL()
	testCheckServiceURL(t, expectedURL, actualURL, "image v2")

	// test for cceV3Client
	serviceClient, err = nil, nil
	serviceClient, err = config.cceV3Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating FlexibleEngine cce v3 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://cce.%s.%s/api/v3/projects/%s/", OS_REGION_NAME, defaultCloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	testCheckServiceURL(t, expectedURL, actualURL, "cce v3")
}

// TestAccServiceEndpoints_Storage test for the endpoints of the clients used in storage
func TestAccServiceEndpoints_Storage(t *testing.T) {

	testAccPreCheckServiceEndpoints(t)

	testProvider := Provider()
	raw := make(map[string]interface{})
	diags := testProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("Unexpected error when configure FlexibleEngine provider: %s", diags[0].Summary)
	}

	config := testProvider.Meta().(*Config)
	var expectedURL, actualURL string
	var serviceClient *golangsdk.ServiceClient
	var err error

	// test for blockStorageV2Client
	serviceClient, err = nil, nil
	serviceClient, err = config.blockStorageV2Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating FlexibleEngine blockStorage v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://evs.%s.%s/v2/%s/", OS_REGION_NAME, defaultCloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	testCheckServiceURL(t, expectedURL, actualURL, "blockStorage v2")

	// test for	sfsV2Client
	serviceClient, err = nil, nil
	serviceClient, err = config.sfsV2Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating FlexibleEngine sfsV2 v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://sfs.%s.%s/v2/%s/", OS_REGION_NAME, defaultCloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	testCheckServiceURL(t, expectedURL, actualURL, "sfsV2 v2")

	// test for sfsV1Client
	serviceClient, err = nil, nil
	serviceClient, err = config.sfsV1Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating FlexibleEngine sfsV1 v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://sfs-turbo.%s.%s/v1/%s/", OS_REGION_NAME, defaultCloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	testCheckServiceURL(t, expectedURL, actualURL, "sfs turbo")

	// test for csbsV1Client
	serviceClient, err = nil, nil
	serviceClient, err = config.csbsV1Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating FlexibleEngine csbsV1 v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://csbs.%s.%s/v1/%s/", OS_REGION_NAME, defaultCloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	testCheckServiceURL(t, expectedURL, actualURL, "csbs v1")

	// test for vbsV2Client
	serviceClient, err = nil, nil
	serviceClient, err = config.vbsV2Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating FlexibleEngine vbsV2 v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://vbs.%s.%s/v2/%s/", OS_REGION_NAME, defaultCloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	testCheckServiceURL(t, expectedURL, actualURL, "vbs v2")
}

// TestAccServiceEndpoints_Network test for the endpoints of the clients used in network
func TestAccServiceEndpoints_Network(t *testing.T) {

	testAccPreCheckServiceEndpoints(t)

	testProvider := Provider()
	raw := make(map[string]interface{})
	diags := testProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("Unexpected error when configure FlexibleEngine provider: %s", diags[0].Summary)
	}

	config := testProvider.Meta().(*Config)
	var expectedURL, actualURL string
	var serviceClient *golangsdk.ServiceClient
	var err error

	// test endpoint of network v1 service
	serviceClient, err = config.networkingV1Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating FlexibleEngine networking v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://vpc.%s.%s/v1/", OS_REGION_NAME, defaultCloud)
	actualURL = serviceClient.ResourceBaseURL()
	testCheckServiceURL(t, expectedURL, actualURL, "vpc v1")

	// test endpoint of network v2 service
	serviceClient, err = nil, nil
	serviceClient, err = config.networkingV2Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating FlexibleEngine networking v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://vpc.%s.%s/v2.0/", OS_REGION_NAME, defaultCloud)
	actualURL = serviceClient.ResourceBaseURL()
	testCheckServiceURL(t, expectedURL, actualURL, "networking v2.0")

	// test endpoint of nat gateway
	serviceClient, err = nil, nil
	serviceClient, err = config.natV2Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating FlexibleEngine nat gateway client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://nat.%s.%s/v2.0/", OS_REGION_NAME, defaultCloud)
	actualURL = serviceClient.ResourceBaseURL()
	testCheckServiceURL(t, expectedURL, actualURL, "nat gateway v2")

	// test endpoint of elb/otc v1.0
	serviceClient, err = nil, nil
	serviceClient, err = config.otcV1Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating FlexibleEngine ELB/otc v1.0 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://elb.%s.%s/v1.0/%s/", OS_REGION_NAME, defaultCloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	testCheckServiceURL(t, expectedURL, actualURL, "elb/otc v1.0")

	// test endpoint of elb v2.0
	serviceClient, err = nil, nil
	serviceClient, err = config.elbV2Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating FlexibleEngine ELB v2.0 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://elb.%s.%s/v2.0/%s/", OS_REGION_NAME, defaultCloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	testCheckServiceURL(t, expectedURL, actualURL, "elb v2.0")

	// test the endpoint of DNS service
	serviceClient, err = config.dnsV2Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating FlexibleEngine DNS client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://dns.%s/v2/", defaultCloud)
	actualURL = serviceClient.ResourceBaseURL()
	testCheckServiceURL(t, expectedURL, actualURL, "dns")

	// test the endpoint of VPC endpoint
	serviceClient, err = config.vpcepV1Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating FlexibleEngine VPC endpoint client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://vpcep.%s.%s/v1/%s/", OS_REGION_NAME, defaultCloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	testCheckServiceURL(t, expectedURL, actualURL, "VPCEP")
}

func TestAccServiceEndpoints_Database(t *testing.T) {
	testAccPreCheckServiceEndpoints(t)

	testProvider := Provider()
	raw := make(map[string]interface{})
	diags := testProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("Unexpected error when configure FlexibleEngine provider: %s", diags[0].Summary)
	}

	var expectedURL, actualURL string
	var serviceClient *golangsdk.ServiceClient
	var err error
	config := testProvider.Meta().(*Config)

	// test the endpoint of RDS v1 service
	serviceClient, err = config.rdsV1Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating FlexibleEngine rds v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://rds.%s.%s/rds/v1/%s/", OS_REGION_NAME, defaultCloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	testCheckServiceURL(t, expectedURL, actualURL, "RDS v1")

	// test the endpoint of RDS v3 service
	serviceClient, err = config.rdsV3Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating FlexibleEngine RDS v3 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://rds.%s.%s/v3/%s/", OS_REGION_NAME, defaultCloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	testCheckServiceURL(t, expectedURL, actualURL, "RDS v3")

	// test the endpoint of DDS v3 service
	serviceClient, err = config.ddsV3Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating FlexibleEngine DDS v3 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://dds.%s.%s/v3/%s/", OS_REGION_NAME, defaultCloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	testCheckServiceURL(t, expectedURL, actualURL, "DDS v3")
}

func TestAccServiceEndpoints_Security(t *testing.T) {
	testAccPreCheckServiceEndpoints(t)

	testProvider := Provider()
	raw := make(map[string]interface{})
	diags := testProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("Unexpected error when configure FlexibleEngine provider: %s", diags[0].Summary)
	}

	var expectedURL, actualURL string
	var serviceClient *golangsdk.ServiceClient
	var err error
	config := testProvider.Meta().(*Config)

	// test the endpoint of anti-ddos service
	serviceClient, err = config.antiddosV1Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating FlexibleEngine anti-ddos client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://antiddos.%s.%s/v1/%s/", OS_REGION_NAME, defaultCloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	testCheckServiceURL(t, expectedURL, actualURL, "anti-ddos")

	// test the endpoint of KMS service
	serviceClient, err = config.kmsKeyV1Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating FlexibleEngine KMS client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://kms.%s.%s/v1.0/", OS_REGION_NAME, defaultCloud)
	actualURL = serviceClient.ResourceBaseURL()
	testCheckServiceURL(t, expectedURL, actualURL, "KMS")
}

func TestAccServiceEndpoints_Application(t *testing.T) {
	testAccPreCheckServiceEndpoints(t)

	testProvider := Provider()
	raw := make(map[string]interface{})
	diags := testProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("Unexpected error when configure FlexibleEngine provider: %s", diags[0].Summary)
	}

	var expectedURL, actualURL string
	var serviceClient *golangsdk.ServiceClient
	var err error
	config := testProvider.Meta().(*Config)

	// test the endpoint of DCS v1 service
	serviceClient, err = config.dcsV1Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating FlexibleEngine dcs v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://dcs.%s.%s/v1.0/%s/", OS_REGION_NAME, defaultCloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	testCheckServiceURL(t, expectedURL, actualURL, "DCS v1")

}

func TestAccServiceEndpoints_EnterpriseIntelligence(t *testing.T) {
	testAccPreCheckServiceEndpoints(t)

	testProvider := Provider()
	raw := make(map[string]interface{})
	diags := testProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("Unexpected error when configure FlexibleEngine provider: %s", diags[0].Summary)
	}

	var expectedURL, actualURL string
	var serviceClient *golangsdk.ServiceClient
	var err error
	config := testProvider.Meta().(*Config)

	// test the endpoint of MRS service
	serviceClient, err = config.MrsV1Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating FlexibleEngine MRS client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://mrs.%s.%s/v1.1/%s/", OS_REGION_NAME, defaultCloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	testCheckServiceURL(t, expectedURL, actualURL, "MRS")

	serviceClient, err = config.MlsV1Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating FlexibleEngine MLS client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://mls.%s.%s/v1.0/%s/", OS_REGION_NAME, defaultCloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	testCheckServiceURL(t, expectedURL, actualURL, "MLS")

	// test the endpoint of SMN service
	serviceClient, err = config.SmnV2Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating FlexibleEngine SMN client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://smn.%s.%s/v2/%s/notifications/", OS_REGION_NAME, defaultCloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	testCheckServiceURL(t, expectedURL, actualURL, "SMN v2")

	// test the endpoint of DWS service
	serviceClient, err = config.dwsV1Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating FlexibleEngine DWS client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://dws.%s.%s/v1.0/%s/", OS_REGION_NAME, defaultCloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	testCheckServiceURL(t, expectedURL, actualURL, "DWS")
}

func TestAccServiceEndpoints_Others(t *testing.T) {
	testAccPreCheckServiceEndpoints(t)

	testProvider := Provider()
	raw := make(map[string]interface{})
	diags := testProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("Unexpected error when configure FlexibleEngine provider: %s", diags[0].Summary)
	}

	var expectedURL, actualURL string
	var serviceClient *golangsdk.ServiceClient
	var err error
	config := testProvider.Meta().(*Config)

	// test the endpoint of DRS service
	serviceClient, err = config.drsV2Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating FlexibleEngine DRS client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://evs.%s.%s/v2/%s/", OS_REGION_NAME, defaultCloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	testCheckServiceURL(t, expectedURL, actualURL, "DRS")

	// test the endpoint of SDRS service
	serviceClient, err = config.sdrsV1Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating FlexibleEngine SDRS client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://sdrs.%s.%s/v1/%s/", OS_REGION_NAME, defaultCloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	testCheckServiceURL(t, expectedURL, actualURL, "SDRS")

	// test the endpoint of RTS service
	serviceClient, err = config.orchestrationV1Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating FlexibleEngine RTS client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://rts.%s.%s/v1/%s/", OS_REGION_NAME, defaultCloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	testCheckServiceURL(t, expectedURL, actualURL, "RTS")
}
