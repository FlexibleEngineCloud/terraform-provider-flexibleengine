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

func getDdsBackupResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := OS_REGION_NAME
	// getBackup: Query DDS backup
	var (
		getBackupHttpUrl = "v3/{project_id}/backups"
		getBackupProduct = "dds"
	)
	getBackupClient, err := cfg.NewServiceClient(getBackupProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DDS Client: %s", err)
	}

	getBackupPath := getBackupClient.Endpoint + getBackupHttpUrl
	getBackupPath = strings.ReplaceAll(getBackupPath, "{project_id}", getBackupClient.ProjectID)

	instanceId := state.Primary.Attributes["instance_id"]
	backupId := state.Primary.ID
	getBackupQueryParams := buildGetBackupQueryParams(instanceId, backupId)
	getBackupPath += getBackupQueryParams

	getBackupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getBackupResp, err := getBackupClient.Request("GET", getBackupPath, &getBackupOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DdsBackup: %s", err)
	}
	getBackupRespBody, err := utils.FlattenResponse(getBackupResp)
	if err != nil {
		return nil, err
	}
	backups := utils.PathSearch("backups", getBackupRespBody, make([]interface{}, 0)).([]interface{})
	if len(backups) == 0 {
		return nil, fmt.Errorf("error get backup by backup ID %s", backupId)
	}

	return backups[0], nil
}

func buildGetBackupQueryParams(instanceId, backupId string) string {
	res := ""
	if instanceId != "" {
		res = fmt.Sprintf("%s&instance_id=%v", res, instanceId)
	}
	if backupId != "" {
		res = fmt.Sprintf("%s&backup_id=%v", res, backupId)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func TestAccDdsBackup_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "flexibleengine_dds_backup.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDdsBackupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDdsBackup_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "this is a test dds instance"),
					resource.TestCheckResourceAttr(rName, "type", "Manual"),
					resource.TestCheckResourceAttr(rName, "status", "COMPLETED"),
					resource.TestCheckResourceAttr(rName, "datastore.0.type", "DDS-Community"),
					acceptance.TestCheckResourceAttrWithVariable(rName, "instance_name",
						"${flexibleengine_dds_instance_v3.instance.name}"),
					acceptance.TestCheckResourceAttrWithVariable(rName, "datastore.0.version",
						"${flexibleengine_dds_instance_v3.instance.datastore.0.version}"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccDdsBackupImportStateFunc(rName),
			},
		},
	})
}

func testAccResourceDdsBackup_base(rName string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_availability_zones" "test" {}

resource "flexibleengine_dds_instance_v3" "instance" {
  name              = "%s"
  availability_zone = data.flexibleengine_availability_zones.test.names[0]
  vpc_id            = flexibleengine_vpc_v1.test.id
  subnet_id         = flexibleengine_vpc_subnet_v1.test.id
  security_group_id = flexibleengine_networking_secgroup_v2.test.id
  password          = "Terraform@123"
  mode              = "Sharding"

  datastore {
    type           = "DDS-Community"
    version        = "3.4"
    storage_engine = "wiredTiger"
  }

  flavor {
    type      = "mongos"
    num       = 2
    spec_code = "dds.mongodb.s3.medium.4.mongos"
  }

  flavor {
    type      = "shard"
    num       = 2
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.s3.medium.4.shard"
  }

  flavor {
    type      = "config"
    num       = 1
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.s3.large.2.config"
  }

  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = "8"
  }

  tags = {
    foo   = "bar"
    owner = "terraform"
  }
}`, testBaseNetwork(rName), rName)
}

func testDdsBackup_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_dds_backup" "test" {
 instance_id = flexibleengine_dds_instance_v3.instance.id
 name        = "%s"
 description = "this is a test dds instance"
}
`, testAccResourceDdsBackup_base(name), name)
}

func testAccDdsBackupImportStateFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.ID == "" || rs.Primary.Attributes["instance_id"] == "" {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.ID), nil
	}
}
