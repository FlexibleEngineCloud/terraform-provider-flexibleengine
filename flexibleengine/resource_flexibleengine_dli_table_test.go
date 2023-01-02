package flexibleengine

import (
	"fmt"
	"strings"
	"testing"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/dli/v1/tables"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getDliTableResourceFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.DliV1Client(OS_REGION_NAME)
	if err != nil {
		return nil, fmtp.Errorf("error creating Dli v1 client, err=%s", err)
	}
	databaseName, tableName := ParseTableInfoFromId(state.Primary.ID)
	return tables.Get(client, databaseName, tableName)
}

func TestAccResourceDliTable_basic(t *testing.T) {
	resourceName := "flexibleengine_dli_table.test"
	name := acceptance.RandomAccResourceName()
	obsBucketName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckS3(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDliTableV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDliTableResource_basic(name, obsBucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDliTableV1Exists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "database_name", name),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "data_location", tables.TableTypeOBS),
					resource.TestCheckResourceAttr(resourceName, "description", "dli table test"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDliTableResource_basic(name string, obsBucketName string) string {

	return fmt.Sprintf(`
resource "flexibleengine_obs_bucket" "test" {
  bucket = "%s"
  acl    = "private"
}


resource "flexibleengine_obs_bucket_object" "test" {
  bucket       = flexibleengine_obs_bucket.test.bucket
  key          = "user/data/user.csv"
  content      = "Jason,Tokyo"
  content_type = "text/plain"
}

resource "flexibleengine_dli_database" "test" {
  name        = "%s"
  description = "For terraform acc test"
}

resource "flexibleengine_dli_table" "test" {
  database_name   = flexibleengine_dli_database.test.name
  name            = "%s"
  data_location   = "OBS"
  description     = "dli table test"
  data_format     = "csv"
  bucket_location = "obs://${flexibleengine_obs_bucket_object.test.bucket}/user/data"

  columns {
    name = "name"
    type        = "string"
    description = "person name"
  }

  columns {
    name = "addrss"
    type        = "string"
    description = "home address"
  }

}
`, obsBucketName, name, name)
}

func testAccCheckDliTableV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.DliV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating dli client, err=%s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_dli_table" {
			continue
		}

		result, err := fetchDliTableV1ByTableNameOnTest(rs.Primary.Attributes["name"], rs.Primary.Attributes["database_name"], client)
		if err == nil && result != nil {
			return fmt.Errorf("dli table still exists: %s,%+v,%+v", rs.Primary.ID, err, result)
		}
	}

	return nil
}

func testAccCheckDliTableV1Exists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		client, err := config.DliV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating dli client, err=%s", err)
		}

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Error checking flexibleengine_dli_table.test exist, err=not found this resource")
		}
		_, err = fetchDliTableV1ByTableNameOnTest(rs.Primary.Attributes["name"], rs.Primary.Attributes["database_name"], client)
		if err != nil {
			if strings.Contains(err.Error(), "Error finding the resource by list api") {
				return fmt.Errorf("flexibleengine_dli_table is not exist")
			}
			return fmt.Errorf("error checking flexibleengine_dli_table.test exist, err=%s", err)
		}
		return nil
	}
}

func fetchDliTableV1ByTableNameOnTest(tableName, databaseName string, client *golangsdk.ServiceClient) (*tables.Table, error) {
	return tables.Get(client, databaseName, tableName)
}
