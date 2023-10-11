package acceptance

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/rds/v3/model"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getRdsDatabaseFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := OS_REGION_NAME
	// getMysqlDatabase: query RDS Mysql database
	var (
		getMysqlDatabaseHttpUrl = "v3/{project_id}/instances/{instance_id}/database/detail?page=1&limit=100"
		getMysqlDatabaseProduct = "rds"
	)
	getMysqlDatabaseClient, err := cfg.NewServiceClient(getMysqlDatabaseProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RDS client: %s", err)
	}

	// Split instance_id and database from resource id
	parts := strings.Split(state.Primary.ID, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, must be <instance_id>/<name>")
	}
	instanceId := parts[0]
	dbName := parts[1]

	getMysqlDatabasePath := getMysqlDatabaseClient.Endpoint + getMysqlDatabaseHttpUrl
	getMysqlDatabasePath = strings.ReplaceAll(getMysqlDatabasePath, "{project_id}", getMysqlDatabaseClient.ProjectID)
	getMysqlDatabasePath = strings.ReplaceAll(getMysqlDatabasePath, "{instance_id}", instanceId)

	getMysqlDatabaseResp, err := pagination.ListAllItems(
		getMysqlDatabaseClient,
		"page",
		getMysqlDatabasePath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return nil, fmt.Errorf("error retrieving MysqlDatabase")
	}

	getMysqlDatabaseRespJson, err := json.Marshal(getMysqlDatabaseResp)
	if err != nil {
		return nil, err
	}
	var getMysqlDatabaseRespBody interface{}
	err = json.Unmarshal(getMysqlDatabaseRespJson, &getMysqlDatabaseRespBody)
	if err != nil {
		return nil, err
	}

	database := utils.PathSearch(fmt.Sprintf("databases[?name=='%s']|[0]", dbName), getMysqlDatabaseRespBody, nil)
	if database != nil {
		return database, nil
	}

	return nil, fmt.Errorf("error get RDS Mysql database by instanceID %s and database %s", instanceId, dbName)
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
