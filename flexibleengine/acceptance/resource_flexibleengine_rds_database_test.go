package acceptance

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/rds/v3/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/rds"
)

func getRdsDatabaseFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.HcRdsV3Client(OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating RDS client: %s", err)
	}

	// Split instance_id and database from resource id
	parts := strings.SplitN(state.Primary.ID, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, must be <instance_id>/<database_name>")
	}
	instanceId := parts[0]
	dbName := parts[1]
	return rds.QueryDatabases(client, instanceId, dbName)
}

func TestAccRdsDatabase_basic(t *testing.T) {
	var database model.DatabaseForCreation
	rName := acceptance.RandomAccResourceName()
	resourceName := "flexibleengine_rds_database.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&database,
		getRdsDatabaseFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testRdsDatabase_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "character_set", "utf8"),
				),
			},
		},
	})
}

func testRdsDatabase_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_rds_database" "test" {
  instance_id   = flexibleengine_rds_instance_v3.test.id
  name          = "%s"
  character_set = "utf8"
}
`, testRdsAccount_base(rName), rName)
}
