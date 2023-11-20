/*
Package acceptance includes all test cases of resources and data sources that
imported directly from huaweicloud.
*/
package acceptance

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	hwacceptance "github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"

	"github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/flexibleengine"
)

var (
	OS_DEPRECATED_ENVIRONMENT     = os.Getenv("OS_DEPRECATED_ENVIRONMENT")
	OS_AVAILABILITY_ZONE          = os.Getenv("OS_AVAILABILITY_ZONE")
	OS_REGION_NAME                = os.Getenv("OS_REGION_NAME")
	OS_ACCESS_KEY                 = os.Getenv("OS_ACCESS_KEY")
	OS_SECRET_KEY                 = os.Getenv("OS_SECRET_KEY")
	OS_PROJECT_ID                 = os.Getenv("OS_PROJECT_ID")
	OS_ENTERPRISE_PROJECT_ID_TEST = os.Getenv("OS_ENTERPRISE_PROJECT_ID_TEST")

	OS_VPC_ID     = os.Getenv("OS_VPC_ID")
	OS_NETWORK_ID = os.Getenv("OS_NETWORK_ID")
	OS_SUBNET_ID  = os.Getenv("OS_SUBNET_ID")

	OS_FLAVOR_ID    = os.Getenv("OS_FLAVOR_ID")
	OS_IMAGE_ID     = os.Getenv("OS_IMAGE_ID")
	OS_KEYPAIR_NAME = os.Getenv("OS_KEYPAIR_NAME")
	OS_FGS_BUCKET   = os.Getenv("OS_FGS_BUCKET")

	OS_DEST_REGION                 = os.Getenv("OS_DEST_REGION")
	OS_DEST_PROJECT_ID             = os.Getenv("OS_DEST_PROJECT_ID")
	OS_NEW_DEST_PROJECT_ID         = os.Getenv("OS_NEW_DEST_PROJECT_ID")
	OS_IMAGE_SHARE_SOURCE_IMAGE_ID = os.Getenv("OS_IMAGE_SHARE_SOURCE_IMAGE_ID")
	OS_SWR_SHARING_ACCOUNT         = os.Getenv("OS_SWR_SHARING_ACCOUNT")
	OS_DLI_FLINK_JAR_OBS_PATH      = os.Getenv("OS_DLI_FLINK_JAR_OBS_PATH")
	OS_WAF_ENABLE_FLAG             = os.Getenv("OS_WAF_ENABLE_FLAG")
	OS_SMS_SOURCE_SERVER           = os.Getenv("OS_SMS_SOURCE_SERVER")
	OS_IMS_BACKUP_ID               = os.Getenv("OS_IMS_BACKUP_ID")
)

// TestAccProviderFactories is a static map containing only the main provider instance
var TestAccProviderFactories map[string]func() (*schema.Provider, error)

// testAccProvider is the "main" provider instance
var testAccProvider *schema.Provider

func init() {
	testAccProvider = flexibleengine.Provider()
	// update TestAccProvider in huaweicloud acceptance package
	hwacceptance.TestAccProvider = testAccProvider

	TestAccProviderFactories = map[string]func() (*schema.Provider, error){
		"flexibleengine": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
}

func testAccPreCheckRequiredEnvVars(t *testing.T) {
	if OS_REGION_NAME == "" {
		t.Fatal("OS_REGION_NAME must be set for acceptance tests")
	}

	if OS_AVAILABILITY_ZONE == "" {
		t.Fatal("OS_AVAILABILITY_ZONE must be set for acceptance tests")
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

func testAccPreCheckAdminOnly(t *testing.T) {
	v := os.Getenv("OS_ADMIN")
	if v != "admin" {
		t.Skip("Skipping test because it requires the admin user")
	}
}

func testAccPreCheckOBS(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)

	if OS_ACCESS_KEY == "" || OS_SECRET_KEY == "" {
		t.Skip("OS_ACCESS_KEY and OS_SECRET_KEY must be set for OBS acceptance tests")
	}
}

func testAccPreCheckReplication(t *testing.T) {
	if OS_DEST_REGION == "" || OS_DEST_PROJECT_ID == "" {
		t.Skip("Skipping the replication policy acceptance tests.")
	}
}

func testAccPreCheckImsCopy(t *testing.T) {
	if OS_DEST_REGION == "" {
		t.Skip("OS_DEST_REGION must be set for IMS copy acceptance tests")
	}
}

func testAccPreCheckImsShare(t *testing.T) {
	if OS_DEST_PROJECT_ID == "" || OS_NEW_DEST_PROJECT_ID == "" {
		t.Skip("OS_DEST_PROJECT_ID and OS_NEW_DEST_PROJECT_ID must be set for IMS share acceptance tests")
	}
}

func testAccPreCheckImageShareAccepter(t *testing.T) {
	if OS_IMAGE_SHARE_SOURCE_IMAGE_ID == "" {
		t.Skip("OS_IMAGE_SHARE_SOURCE_IMAGE_ID must be set for IMS share acceptance tests")
	}
}

func testAccPreCheckSWRDomian(t *testing.T) {
	if OS_SWR_SHARING_ACCOUNT == "" {
		t.Skip("OS_SWR_SHARING_ACCOUNT must be set for swr domian tests, " +
			"the value of OS_SWR_SHARING_ACCOUNT should be another IAM user name")
	}
}

func testAccPreCheckDliJarPath(t *testing.T) {
	if OS_DLI_FLINK_JAR_OBS_PATH == "" {
		t.Skip("OS_DLI_FLINK_JAR_OBS_PATH must be set for DLI Flink Jar job acceptance tests.")
	}
}

func testAccPrecheckWafInstance(t *testing.T) {
	if OS_WAF_ENABLE_FLAG == "" {
		t.Skip("Jump the WAF acceptance tests.")
	}
}
func testAccPreCheckSms(t *testing.T) {
	if OS_SMS_SOURCE_SERVER == "" {
		t.Skip("OS_SMS_SOURCE_SERVER must be set for SMS acceptance tests")
	}
}

func testAccPreCheckImsBackupId(t *testing.T) {
	if OS_IMS_BACKUP_ID == "" {
		t.Skip("OS_IMS_BACKUP_ID must be set for IMS whole image with CBR backup id")
	}
}

func testAccPreCheckEpsID(t *testing.T) {
	// The environment variables in tests take OS_ENTERPRISE_PROJECT_ID_TEST instead of OS_ENTERPRISE_PROJECT_ID to
	// ensure that other data-resources that support enterprise projects query the default project without being
	// affected by this variable.
	if OS_ENTERPRISE_PROJECT_ID_TEST == "" {
		t.Skip("The environment variables does not support Enterprise Project ID for acc tests")
	}
}
