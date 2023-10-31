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

func getSQLAuditResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := OS_REGION_NAME
	// getSQLAudit: Query the RDS SQL audit
	var (
		getSQLAuditHttpUrl = "v3/{project_id}/instances/{instance_id}/auditlog-policy"
		getSQLAuditProduct = "rds"
	)
	getSQLAuditClient, err := cfg.NewServiceClient(getSQLAuditProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RDS client: %s", err)
	}

	getSQLAuditPath := getSQLAuditClient.Endpoint + getSQLAuditHttpUrl
	getSQLAuditPath = strings.ReplaceAll(getSQLAuditPath, "{project_id}", getSQLAuditClient.ProjectID)
	getSQLAuditPath = strings.ReplaceAll(getSQLAuditPath, "{instance_id}", state.Primary.ID)

	getSQLAuditOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getSQLAuditResp, err := getSQLAuditClient.Request("GET", getSQLAuditPath, &getSQLAuditOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving RDS SQL audit: %s", err)
	}

	getSQLAuditRespBody, err := utils.FlattenResponse(getSQLAuditResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving RDS SQL audit: %s", err)
	}

	keepDays := utils.PathSearch("keep_days", getSQLAuditRespBody, 0).(float64)
	if keepDays == 0 {
		return nil, fmt.Errorf("error retrieving RDS SQL audit: %s", err)
	}

	return getSQLAuditRespBody, nil
}

func TestAccSQLAudit_basic(t *testing.T) {
	var obj interface{}
	name := acceptance.RandomAccResourceName()
	rName := "flexibleengine_rds_sql_audit.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getSQLAuditResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testSQLAudit_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"flexibleengine_rds_instance_v3.test", "id"),
					resource.TestCheckResourceAttr(rName, "keep_days", "5"),
				),
			},
			{
				Config: testSQLAudit_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"flexibleengine_rds_instance_v3.test", "id"),
					resource.TestCheckResourceAttr(rName, "keep_days", "9"),
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

func testAccRdsSqlAudit_base(name string) string {
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

func testSQLAudit_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_rds_sql_audit" "test" {
  instance_id = flexibleengine_rds_instance_v3.test.id
  keep_days   = "5"
}
`, testAccRdsSqlAudit_base(name))
}

func testSQLAudit_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_rds_sql_audit" "test" {
  instance_id = flexibleengine_rds_instance_v3.test.id
  keep_days   = "9"
}
`, testAccRdsSqlAudit_base(name))
}
