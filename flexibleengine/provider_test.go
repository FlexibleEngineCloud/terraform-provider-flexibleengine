package flexibleengine

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/pathorcontents"
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
	OS_DELEGATED_DOMAIN_NAME  = os.Getenv("OS_DELEGATED_DOMAIN_NAME")
	OS_TENANT_NAME            = getTenantName()
)

// testAccProviders is a static map containing only the main provider instance.
//
// Deprecated: Terraform Plugin SDK version 2 uses TestCase.ProviderFactories
// but supports this value in TestCase.Providers for backwards compatibility.
// In the future Providers: testAccProviders will be changed to
// ProviderFactories: testAccProviderFactories
var testAccProviders map[string]*schema.Provider

// testAccProviderFactories is a static map containing only the main provider instance
var TestAccProviderFactories map[string]func() (*schema.Provider, error)

// testAccProvider is the "main" provider instance
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()

	testAccProviders = map[string]*schema.Provider{
		"flexibleengine": testAccProvider,
	}

	TestAccProviderFactories = map[string]func() (*schema.Provider, error){
		"flexibleengine": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
}

// ServiceFunc the HuaweiCloud resource query functions.
type ServiceFunc func(*Config, *terraform.ResourceState) (interface{}, error)

// resourceCheck resource check object, only used in the package.
type resourceCheck struct {
	resourceName    string
	resourceObject  interface{}
	getResourceFunc ServiceFunc
	resourceType    string
}

const (
	resourceTypeCode   = "resource"
	dataSourceTypeCode = "dataSource"

	checkAttrRegexpStr = `^\$\{([^\}]+)\}$`
)

/*
initDataSourceCheck build a 'resourceCheck' object. Only used to check datasource attributes.
  Parameters:
    resourceName:    The resource name is used to check in the terraform.State.e.g. : huaweicloud_waf_domain.domain_1.
  Return:
    *resourceCheck: resourceCheck object
*/
func initDataSourceCheck(sourceName string) *resourceCheck {
	return &resourceCheck{
		resourceName: sourceName,
		resourceType: dataSourceTypeCode,
	}
}

/*
initResourceCheck build a 'resourceCheck' object. The common test methods are provided in 'resourceCheck'.
  Parameters:
    resourceName:    The resource name is used to check in the terraform.State.e.g. : huaweicloud_waf_domain.domain_1.
    resourceObject:  Resource object, used to check whether the resource exists in HuaweiCloud.
    getResourceFunc: The function used to get the resource object.
  Return:
    *resourceCheck: resourceCheck object
*/
func initResourceCheck(resourceName string, resourceObject interface{}, getResourceFunc ServiceFunc) *resourceCheck {
	return &resourceCheck{
		resourceName:    resourceName,
		resourceObject:  resourceObject,
		getResourceFunc: getResourceFunc,
		resourceType:    resourceTypeCode,
	}
}

func parseVariableToName(varStr string) (string, string, error) {
	var resName, keyName string
	// Check the format of the variable.
	match, _ := regexp.MatchString(checkAttrRegexpStr, varStr)
	if !match {
		return resName, keyName, fmt.Errorf("The type of 'variable' is error, "+
			"expected ${resourceType.name.field} got %s", varStr)
	}

	reg, err := regexp.Compile(checkAttrRegexpStr)
	if err != nil {
		return resName, keyName, fmt.Errorf("The acceptance function is wrong.")
	}
	mArr := reg.FindStringSubmatch(varStr)
	if len(mArr) != 2 {
		return resName, keyName, fmt.Errorf("The type of 'variable' is error, "+
			"expected ${resourceType.name.field} got %s", varStr)
	}

	// Get resName and keyName from variable.
	strs := strings.Split(mArr[1], ".")
	for i, s := range strs {
		if strings.Contains(s, "huaweicloud_") {
			resName = strings.Join(strs[0:i+2], ".")
			keyName = strings.Join(strs[i+2:], ".")
			break
		}
	}
	return resName, keyName, nil
}

/*
testCheckResourceAttrWithVariable validates the variable in state for the given name/key combination.
  Parameters:
    resourceName: The resource name is used to check in the terraform.State.
    key:          The field name of the resource.
    variable:     The variable name of the value to be checked.

    variable such like ${huaweicloud_waf_certificate.certificate_1.id}
    or ${data.huaweicloud_waf_policies.policies_2.policies.0.id}
*/
func testCheckResourceAttrWithVariable(resourceName, key, varStr string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resName, keyName, err := parseVariableToName(varStr)
		if err != nil {
			return err
		}

		if strings.EqualFold(resourceName, resName) {
			return fmt.Errorf("Meaningless verification. " +
				"The referenced resource cannot be the current resource.")
		}

		// Get the value based on resName and keyName from the state.
		rs, ok := s.RootModule().Resources[resName]
		if !ok {
			return fmt.Errorf("Can't find %s in state : %v.", resName, ok)
		}
		value := rs.Primary.Attributes[keyName]

		return resource.TestCheckResourceAttr(resourceName, key, value)(s)
	}
}

// CheckResourceDestroy check whether resources destroied in HuaweiCloud.
func (rc *resourceCheck) CheckResourceDestroy() resource.TestCheckFunc {
	if strings.Compare(rc.resourceType, dataSourceTypeCode) == 0 {
		fmt.Errorf("Error, you built a resourceCheck with 'initDataSourceCheck', " +
			"it cannot run CheckResourceDestroy().")
		return nil
	}
	return func(s *terraform.State) error {
		strs := strings.Split(rc.resourceName, ".")
		var resourceType string
		for _, str := range strs {
			if strings.Contains(str, "huaweicloud_") {
				resourceType = strings.Trim(str, " ")
				break
			}
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != resourceType {
				continue
			}

			conf := testAccProvider.Meta().(*Config)
			if rc.getResourceFunc != nil {
				if _, err := rc.getResourceFunc(conf, rs); err == nil {
					return fmt.Errorf("failed to destroy resource. The resource of %s : %s still exists.",
						resourceType, rs.Primary.ID)
				}
			} else {
				return fmt.Errorf("The 'getResourceFunc' is nil, please set it during initialization.")
			}
		}
		return nil
	}
}

// CheckResourceExists check whether resources exist in HuaweiCloud.
func (rc *resourceCheck) CheckResourceExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rc.resourceName]
		if !ok {
			return fmt.Errorf("Can not found the resource or data source in state: %s", rc.resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No id set for the resource or data source: %s", rc.resourceName)
		}
		if strings.EqualFold(rc.resourceType, dataSourceTypeCode) {
			return nil
		}

		if rc.getResourceFunc != nil {
			conf := testAccProvider.Meta().(*Config)
			r, err := rc.getResourceFunc(conf, rs)
			if err != nil {
				return fmt.Errorf("checking resource %s %s exists error: %s ",
					rc.resourceName, rs.Primary.ID, err)
			}
			if rc.resourceObject != nil {
				b, err := json.Marshal(r)
				if err != nil {
					return fmt.Errorf("marshaling resource %s %s error: %s ",
						rc.resourceName, rs.Primary.ID, err)
				}
				json.Unmarshal(b, rc.resourceObject)
			} else {
				log.Printf("[WARN] The 'resourceObject' is nil, please set it during initialization.")
			}
		} else {
			return fmt.Errorf("The 'getResourceFunc' is nil, please set it.")
		}

		return nil
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
	if OS_REGION_NAME == "" {
		t.Fatal("OS_REGION_NAME must be set for acceptance tests")
	}

	v := os.Getenv("OS_AUTH_URL")
	if v == "" {
		t.Fatal("OS_AUTH_URL must be set for acceptance tests")
	}

	if OS_IMAGE_ID == "" && OS_IMAGE_NAME == "" {
		t.Fatal("OS_IMAGE_ID or OS_IMAGE_NAME must be set for acceptance tests")
	}

	if OS_FLAVOR_ID == "" && OS_FLAVOR_NAME == "" {
		t.Fatal("OS_FLAVOR_ID or OS_FLAVOR_NAME must be set for acceptance tests")
	}

	if OS_NETWORK_ID == "" {
		t.Fatal("OS_NETWORK_ID must be set for acceptance tests")
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

func testAccPreCheckFloatingIP(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)

	if OS_POOL_NAME == "" {
		t.Fatal("OS_POOL_NAME must be set for floating network tests")
	}

	if OS_EXTGW_ID == "" {
		t.Fatal("OS_EXTGW_ID must be set for floating network tests")
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
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
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

	diags := p.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("Unexpected err when specifying FlexibleEngine CA by file: %s", diags[0].Summary)
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

	diags := p.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("Unexpected err when specifying FlexibleEngine CA by string: %s", diags[0].Summary)
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

	diags := p.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("Unexpected err when specifying FlexibleEngine Client keypair by file: %s", diags[0].Summary)
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

	diags := p.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("Unexpected err when specifying FlexibleEngine Client keypair by contents: %s", diags[0].Summary)
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
