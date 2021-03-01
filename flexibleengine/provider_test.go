package flexibleengine

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/pathorcontents"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var (
	OS_DEPRECATED_ENVIRONMENT = os.Getenv("OS_DEPRECATED_ENVIRONMENT")
	OS_EXTGW_ID               = os.Getenv("OS_EXTGW_ID")
	OS_FLAVOR_ID              = os.Getenv("OS_FLAVOR_ID")
	OS_FLAVOR_NAME            = os.Getenv("OS_FLAVOR_NAME")
	OS_IMAGE_ID               = os.Getenv("OS_IMAGE_ID")
	OS_IMAGE_NAME             = os.Getenv("OS_IMAGE_NAME")
	OS_NETWORK_ID             = os.Getenv("OS_NETWORK_ID")
	OS_POOL_NAME              = os.Getenv("OS_POOL_NAME")
	OS_REGION_NAME            = os.Getenv("OS_REGION_NAME")
	OS_ACCESS_KEY             = os.Getenv("OS_ACCESS_KEY")
	OS_SECRET_KEY             = os.Getenv("OS_SECRET_KEY")
	OS_AVAILABILITY_ZONE      = os.Getenv("OS_AVAILABILITY_ZONE")
	OS_VPC_ID                 = os.Getenv("OS_VPC_ID")
	OS_SUBNET_ID              = os.Getenv("OS_SUBNET_ID")
	OS_KEYPAIR_NAME           = os.Getenv("OS_KEYPAIR_NAME")
	OS_BMS_FLAVOR_NAME        = os.Getenv("OS_BMS_FLAVOR_NAME")
	OS_MRS_ENVIRONMENT        = os.Getenv("OS_MRS_ENVIRONMENT")
	OS_SDRS_ENVIRONMENT       = os.Getenv("OS_SDRS_ENVIRONMENT")
	OS_TENANT_NAME            = getTenantName()
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"flexibleengine": testAccProvider,
	}
}

func getTenantName() string {
	tn := os.Getenv("OS_TENANT_NAME")
	if tn == "" {
		tn = os.Getenv("OS_PROJECT_NAME")
	}
	return tn
}

func testAccPreCheckRequiredEnvVars(t *testing.T) {
	v := os.Getenv("OS_AUTH_URL")
	if v == "" {
		t.Fatal("OS_AUTH_URL must be set for acceptance tests")
	}

	if OS_IMAGE_ID == "" && OS_IMAGE_NAME == "" {
		t.Fatal("OS_IMAGE_ID or OS_IMAGE_NAME must be set for acceptance tests")
	}

	if OS_POOL_NAME == "" {
		t.Fatal("OS_POOL_NAME must be set for acceptance tests")
	}

	if OS_REGION_NAME == "" {
		t.Fatal("OS_REGION_NAME must be set for acceptance tests")
	}

	if OS_REGION_NAME == "" {
		t.Fatal("OS_REGION_NAME must be set for acceptance tests")
	}

	if OS_FLAVOR_ID == "" && OS_FLAVOR_NAME == "" {
		t.Fatal("OS_FLAVOR_ID or OS_FLAVOR_NAME must be set for acceptance tests")
	}

	if OS_NETWORK_ID == "" {
		t.Fatal("OS_NETWORK_ID must be set for acceptance tests")
	}

	if OS_EXTGW_ID == "" {
		t.Fatal("OS_EXTGW_ID must be set for acceptance tests")
	}
}

func testAccPreCheck(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)

	// Do not run the test if this is a deprecated testing environment.
	if OS_DEPRECATED_ENVIRONMENT != "" {
		t.Skip("This environment only runs deprecated tests")
	}
}

func testAccPreCheckDeprecated(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)

	if OS_DEPRECATED_ENVIRONMENT == "" {
		t.Skip("This environment does not support deprecated tests")
	}
}

func testAccPreCheckMrs(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)

	if OS_MRS_ENVIRONMENT == "" {
		t.Skip("This environment does not support MRS tests")
	}
}

func testAccPreCheckSdrs(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)

	if OS_SDRS_ENVIRONMENT == "" {
		t.Skip("This environment does not support SDRS tests")
	}
}

func testAccPreCheckAdminOnly(t *testing.T) {
	v := os.Getenv("OS_ADMIN")
	if v != "admin" {
		t.Skip("Skipping test because it requires the admin user")
	}
}

func testAccPreCheckS3(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)

	if OS_ACCESS_KEY == "" || OS_SECRET_KEY == "" {
		t.Skip("OS_ACCESS_KEY and OS_SECRET_KEY must be set for OBS/S3 acceptance tests")
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

// Steps for configuring FlexibleEngine with SSL validation are here:
// https://github.com/hashicorp/terraform/pull/6279#issuecomment-219020144
func TestAccProvider_caCertFile(t *testing.T) {
	if os.Getenv("TF_ACC") == "" || os.Getenv("OS_SSL_TESTS") == "" {
		t.Skip("TF_ACC or OS_SSL_TESTS not set, skipping FlexibleEngine SSL test.")
	}
	if os.Getenv("OS_CACERT") == "" {
		t.Skip("OS_CACERT is not set; skipping FlexibleEngine CA test.")
	}

	p := Provider()

	caFile, err := envVarFile("OS_CACERT")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(caFile)

	raw := map[string]interface{}{
		"cacert_file": caFile,
	}

	err = p.Configure(terraform.NewResourceConfigRaw(raw))
	if err != nil {
		t.Fatalf("Unexpected err when specifying FlexibleEngine CA by file: %s", err)
	}
}

func TestAccProvider_caCertString(t *testing.T) {
	if os.Getenv("TF_ACC") == "" || os.Getenv("OS_SSL_TESTS") == "" {
		t.Skip("TF_ACC or OS_SSL_TESTS not set, skipping FlexibleEngine SSL test.")
	}
	if os.Getenv("OS_CACERT") == "" {
		t.Skip("OS_CACERT is not set; skipping FlexibleEngine CA test.")
	}

	p := Provider()

	caContents, err := envVarContents("OS_CACERT")
	if err != nil {
		t.Fatal(err)
	}
	raw := map[string]interface{}{
		"cacert_file": caContents,
	}

	err = p.Configure(terraform.NewResourceConfigRaw(raw))
	if err != nil {
		t.Fatalf("Unexpected err when specifying FlexibleEngine CA by string: %s", err)
	}
}

func TestAccProvider_clientCertFile(t *testing.T) {
	if os.Getenv("TF_ACC") == "" || os.Getenv("OS_SSL_TESTS") == "" {
		t.Skip("TF_ACC or OS_SSL_TESTS not set, skipping FlexibleEngine SSL test.")
	}
	if os.Getenv("OS_CERT") == "" || os.Getenv("OS_KEY") == "" {
		t.Skip("OS_CERT or OS_KEY is not set; skipping FlexibleEngine client SSL auth test.")
	}

	p := Provider()

	certFile, err := envVarFile("OS_CERT")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(certFile)
	keyFile, err := envVarFile("OS_KEY")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(keyFile)

	raw := map[string]interface{}{
		"cert": certFile,
		"key":  keyFile,
	}

	err = p.Configure(terraform.NewResourceConfigRaw(raw))
	if err != nil {
		t.Fatalf("Unexpected err when specifying FlexibleEngine Client keypair by file: %s", err)
	}
}

func TestAccProvider_clientCertString(t *testing.T) {
	if os.Getenv("TF_ACC") == "" || os.Getenv("OS_SSL_TESTS") == "" {
		t.Skip("TF_ACC or OS_SSL_TESTS not set, skipping FlexibleEngine SSL test.")
	}
	if os.Getenv("OS_CERT") == "" || os.Getenv("OS_KEY") == "" {
		t.Skip("OS_CERT or OS_KEY is not set; skipping FlexibleEngine client SSL auth test.")
	}

	p := Provider()

	certContents, err := envVarContents("OS_CERT")
	if err != nil {
		t.Fatal(err)
	}
	keyContents, err := envVarContents("OS_KEY")
	if err != nil {
		t.Fatal(err)
	}

	raw := map[string]interface{}{
		"cert": certContents,
		"key":  keyContents,
	}

	err = p.Configure(terraform.NewResourceConfigRaw(raw))
	if err != nil {
		t.Fatalf("Unexpected err when specifying FlexibleEngine Client keypair by contents: %s", err)
	}
}

func envVarContents(varName string) (string, error) {
	contents, _, err := pathorcontents.Read(os.Getenv(varName))
	if err != nil {
		return "", fmt.Errorf("Error reading %s: %s", varName, err)
	}
	return contents, nil
}

func envVarFile(varName string) (string, error) {
	contents, err := envVarContents(varName)
	if err != nil {
		return "", err
	}

	tmpFile, err := ioutil.TempFile("", varName)
	if err != nil {
		return "", fmt.Errorf("Error creating temp file: %s", err)
	}
	if _, err := tmpFile.Write([]byte(contents)); err != nil {
		_ = os.Remove(tmpFile.Name())
		return "", fmt.Errorf("Error writing temp file: %s", err)
	}
	if err := tmpFile.Close(); err != nil {
		_ = os.Remove(tmpFile.Name())
		return "", fmt.Errorf("Error closing temp file: %s", err)
	}
	return tmpFile.Name(), nil
}

func testAccBmsKeyPairPreCheck(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)
	if OS_KEYPAIR_NAME == "" {
		t.Skip("Provide the key pair name")
	}
}

func testAccBmsFlavorPreCheck(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)
	if OS_BMS_FLAVOR_NAME == "" {
		t.Skip("Provide the bms name starting with 'physical'")
	}
}

func testAccCCEKeyPairPreCheck(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)
	if OS_KEYPAIR_NAME == "" {
		t.Skip("OS_KEYPAIR_NAME must be set for acceptance tests")
	}
}

func testAccIdentityV3AgencyPreCheck(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)
	if OS_TENANT_NAME == "" {
		t.Skip("OS_TENANT_NAME must be set for acceptance tests")
	}
}
