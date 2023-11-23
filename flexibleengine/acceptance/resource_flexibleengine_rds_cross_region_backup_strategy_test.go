package acceptance

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getBackupStrategyResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := OS_REGION_NAME
	// getBackupStrategy: Query the RDS cross region backup strategy
	var (
		getBackupStrategyHttpUrl = "v3/{project_id}/instances/{instance_id}/backups/offsite-policy"
		getBackupStrategyProduct = "rds"
	)
	getBackupStrategyClient, err := cfg.NewServiceClient(getBackupStrategyProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RDS client: %s", err)
	}

	getBackupStrategyPath := getBackupStrategyClient.Endpoint + getBackupStrategyHttpUrl
	getBackupStrategyPath = strings.ReplaceAll(getBackupStrategyPath, "{project_id}", getBackupStrategyClient.ProjectID)
	getBackupStrategyPath = strings.ReplaceAll(getBackupStrategyPath, "{instance_id}", state.Primary.ID)

	getBackupStrategyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getBackupStrategyResp, err := getBackupStrategyClient.Request("GET", getBackupStrategyPath, &getBackupStrategyOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving RDS cross region backup strategy: %s", err)
	}

	getBackupStrategyRespBody, err := utils.FlattenResponse(getBackupStrategyResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving RDS cross region backup strategy: %s", err)
	}

	policyPara := utils.PathSearch("policy_para", getBackupStrategyRespBody, nil)
	if policyPara == nil {
		return nil, fmt.Errorf("error retrieving RDS cross region backup strategy: %s", err)
	}

	backupStrategies := policyPara.([]interface{})
	if len(backupStrategies) == 0 || utils.PathSearch("keep_days", backupStrategies[0], 0).(float64) == 0 {
		return nil, fmt.Errorf("error retrieving RDS cross region backup strategy: %s", err)
	}

	return getBackupStrategyRespBody, nil
}

func TestAccBackupStrategy_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "flexibleengine_rds_cross_region_backup_strategy.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getBackupStrategyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			acceptance.TestAccPreCheckReplication(t)
		},
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testBackupStrategy_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"flexibleengine_rds_instance_v3.test", "id"),
					resource.TestCheckResourceAttr(rName, "backup_type", "auto"),
					resource.TestCheckResourceAttr(rName, "keep_days", "5"),
					resource.TestCheckResourceAttr(rName, "destination_region", OS_DEST_REGION),
					resource.TestCheckResourceAttr(rName, "destination_project_id", OS_DEST_PROJECT_ID),
				),
			},
			{
				Config: testBackupStrategy_basic_update1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"flexibleengine_rds_instance_v3.test", "id"),
					resource.TestCheckResourceAttr(rName, "backup_type", "all"),
					resource.TestCheckResourceAttr(rName, "keep_days", "8"),
					resource.TestCheckResourceAttr(rName, "destination_region", OS_DEST_REGION),
					resource.TestCheckResourceAttr(rName, "destination_project_id", OS_DEST_PROJECT_ID),
				),
			},
			{
				Config: testBackupStrategy_basic_update2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"flexibleengine_rds_instance_v3.test", "id"),
					resource.TestCheckResourceAttr(rName, "backup_type", "auto"),
					resource.TestCheckResourceAttr(rName, "keep_days", "10"),
					resource.TestCheckResourceAttr(rName, "destination_region", OS_DEST_REGION),
					resource.TestCheckResourceAttr(rName, "destination_project_id", OS_DEST_PROJECT_ID),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccRdsCrossRegionBackupStrategy(name string) string {
	return fmt.Sprintf(`
%[1]s

data "flexibleengine_availability_zones" "test" {}

data "flexibleengine_networking_secgroup_v2" "test" {
  name = "default"
}

data "flexibleengine_rds_flavors_v3" "test" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "single"
}

resource "flexibleengine_rds_instance_v3" "test" {
  name              = "%[2]s"
  flavor            = data.flexibleengine_rds_flavors_v3.test.flavors[2].name
  availability_zone = [data.flexibleengine_availability_zones.test.names[0]]
  security_group_id = data.flexibleengine_networking_secgroup_v2.test.id
  subnet_id         = flexibleengine_vpc_subnet_v1.test.id
  vpc_id            = flexibleengine_vpc_v1.test.id
  time_zone         = "UTC+08:00"

  db {
    password = "FlexibleEngine!120521"
    type     = "MySQL"
    version  = "8.0"
    port     = 8630
  }

  volume {
    type = "COMMON"
    size = 50
  }

  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }

  lifecycle {
    ignore_changes = [
      backup_strategy,
    ]
  }
}
`, testVpc(name), name)
}

func testBackupStrategy_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_rds_cross_region_backup_strategy" "test" {
  instance_id            = flexibleengine_rds_instance_v3.test.id
  backup_type            = "auto"
  keep_days              = "5"
  destination_region     = "%s"
  destination_project_id = "%s"
}
`, testAccRdsCrossRegionBackupStrategy(name), OS_DEST_REGION, OS_DEST_PROJECT_ID)
}

func testBackupStrategy_basic_update1(name string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_rds_cross_region_backup_strategy" "test" {
  instance_id            = flexibleengine_rds_instance_v3.test.id
  backup_type            = "all"
  keep_days              = "8"
  destination_region     = "%s"
  destination_project_id = "%s"
}
`, testAccRdsCrossRegionBackupStrategy(name), OS_DEST_REGION, OS_DEST_PROJECT_ID)
}

func testBackupStrategy_basic_update2(name string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_rds_cross_region_backup_strategy" "test" {
  instance_id            = flexibleengine_rds_instance_v3.test.id
  backup_type            = "auto"
  keep_days              = "10"
  destination_region     = "%s"
  destination_project_id = "%s"
}
`, testAccRdsCrossRegionBackupStrategy(name), OS_DEST_REGION, OS_DEST_PROJECT_ID)
}
