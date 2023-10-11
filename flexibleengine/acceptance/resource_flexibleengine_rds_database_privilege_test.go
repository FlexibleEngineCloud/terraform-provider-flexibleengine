package acceptance

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/pagination"
	model "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/rds/v3/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getRdsDatabasePrivilegeFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := OS_REGION_NAME
	// getMysqlDatabasePrivilege: query RDS Mysql database privilege
	var (
		getMysqlDatabasePrivilegeHttpUrl = "v3/{project_id}/instances/{instance_id}/database/db_user"
		getMysqlDatabasePrivilegeProduct = "rds"
	)
	getMysqlDatabasePrivilegeClient, err := cfg.NewServiceClient(getMysqlDatabasePrivilegeProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RDS client: %s", err)
	}

	// Split instance_id and database from resource id
	parts := strings.Split(state.Primary.ID, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, must be <instance_id>/<db_name>")
	}
	instanceId := parts[0]
	dbName := parts[1]

	getMysqlDatabasePrivilegePath := getMysqlDatabasePrivilegeClient.Endpoint + getMysqlDatabasePrivilegeHttpUrl
	getMysqlDatabasePrivilegePath = strings.ReplaceAll(getMysqlDatabasePrivilegePath, "{project_id}",
		getMysqlDatabasePrivilegeClient.ProjectID)
	getMysqlDatabasePrivilegePath = strings.ReplaceAll(getMysqlDatabasePrivilegePath, "{instance_id}", instanceId)

	getMysqlDatabasePrivilegeQueryParams := buildGetMysqlDatabasePrivilegeQueryParams(dbName)
	getMysqlDatabasePrivilegePath += getMysqlDatabasePrivilegeQueryParams

	getMysqlDatabasePrivilegeResp, err := pagination.ListAllItems(
		getMysqlDatabasePrivilegeClient,
		"page",
		getMysqlDatabasePrivilegePath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return nil, fmt.Errorf("error retrieving Mysql database privilege: %s", err)
	}

	getMysqlDatabasePrivilegeRespJson, err := json.Marshal(getMysqlDatabasePrivilegeResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Mysql database privilege: %s", err)
	}
	var getMysqlDatabasePrivilegeRespBody interface{}
	err = json.Unmarshal(getMysqlDatabasePrivilegeRespJson, &getMysqlDatabasePrivilegeRespBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Mysql database privilege: %s", err)
	}

	curJson := utils.PathSearch("users", getMysqlDatabasePrivilegeRespBody, make([]interface{}, 0))
	if len(curJson.([]interface{})) == 0 {
		return nil, fmt.Errorf("error get RDS Mysql database privilege")
	}

	return getMysqlDatabasePrivilegeRespBody, nil
}

func buildGetMysqlDatabasePrivilegeQueryParams(dbName string) string {
	return fmt.Sprintf("?db-name=%s&page=1&limit=100", dbName)
}

func TestAccRdsDatabasePrivilege_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "flexibleengine_rds_database_privilege.test"
	var users []model.UserWithPrivilege
	rc := acceptance.InitResourceCheck(
		resourceName,
		&users,
		getRdsDatabasePrivilegeFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsDatabasePrivilege_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "users.0.name",
						"flexibleengine_rds_account.test", "name"),
					resource.TestCheckResourceAttr(resourceName, "users.0.readonly", "false"),
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

func testAccRdsDatabasePrivilege_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_rds_account" "test" {
  instance_id = flexibleengine_rds_instance_v3.test.id
  name        = "%s"
  password    = "Test@12345678"
}

resource "flexibleengine_rds_database_privilege" "test" {
  instance_id = flexibleengine_rds_instance_v3.test.id
  db_name     = flexibleengine_rds_database.test.name

  users {
    name = flexibleengine_rds_account.test.name
  }
}
`, testRdsDatabase_basic(rName), rName)
}
