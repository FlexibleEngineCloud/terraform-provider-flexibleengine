package acceptance

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/dli/v1/tables"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dli"
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
	databaseName, tableName := dli.ParseTableInfoFromId(state.Primary.ID)
	return tables.Get(client, databaseName, tableName)
}

func TestAccResourceDliTable_basic(t *testing.T) {
	var TableObj tables.CreateTableOpts
	resourceName := "flexibleengine_dli_table.test"
	name := acceptance.RandomAccResourceName()
	obsBucketName := acceptance.RandomAccResourceNameWithDash()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&TableObj,
		getDliTableResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckOBS(t)
		},
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDliTableResource_basic(name, obsBucketName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
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
