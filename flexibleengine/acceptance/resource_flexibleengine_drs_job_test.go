package acceptance

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/drs/v3/jobs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getDrsJobResourceFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.DrsV3Client(OS_REGION_NAME)
	if err != nil {
		return nil, fmtp.Errorf("error creating DRS client, err=%s", err)
	}
	return jobs.Get(client, jobs.QueryJobReq{Jobs: []string{state.Primary.ID}})
}

func TestAccResourceDrsJob_basic(t *testing.T) {
	var obj jobs.BatchCreateJobReq
	resourceName := "flexibleengine_drs_job.test"
	name := acceptance.RandomAccResourceName()
	dbName := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()
	pwd := "TestDrs@123"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDrsJobResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDrsJob_migrate_mysql(name, dbName, pwd),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "type", "migration"),
					resource.TestCheckResourceAttr(resourceName, "direction", "up"),
					resource.TestCheckResourceAttr(resourceName, "net_type", "eip"),
					resource.TestCheckResourceAttr(resourceName, "destination_db_readnoly", "true"),
					resource.TestCheckResourceAttr(resourceName, "migration_type", "FULL_INCR_TRANS"),
					resource.TestCheckResourceAttr(resourceName, "description", name),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.engine_type", "mysql"),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.ip", "192.168.0.58"),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.port", "3306"),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.user", "root"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.engine_type", "mysql"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.ip", "192.168.0.59"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.port", "3306"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.user", "root"),
					resource.TestCheckResourceAttrPair(resourceName, "destination_db.0.subnet_id",
						"flexibleengine_vpc_subnet_v1.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "destination_db.0.instance_id",
						"flexibleengine_rds_instance_v3.test2", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "destination_db.0.region",
						"flexibleengine_rds_instance_v3.test2", "region"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "public_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "private_ip"),
				),
			},
			{
				Config: testAccDrsJob_migrate_mysql(updateName, dbName, pwd),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "type", "migration"),
					resource.TestCheckResourceAttr(resourceName, "direction", "up"),
					resource.TestCheckResourceAttr(resourceName, "net_type", "eip"),
					resource.TestCheckResourceAttr(resourceName, "destination_db_readnoly", "true"),
					resource.TestCheckResourceAttr(resourceName, "migration_type", "FULL_INCR_TRANS"),
					resource.TestCheckResourceAttr(resourceName, "description", updateName),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.engine_type", "mysql"),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.ip", "192.168.0.58"),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.port", "3306"),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.user", "root"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.engine_type", "mysql"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.ip", "192.168.0.59"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.port", "3306"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.user", "root"),
					resource.TestCheckResourceAttrPair(resourceName, "destination_db.0.subnet_id",
						"flexibleengine_vpc_subnet_v1.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "destination_db.0.instance_id",
						"flexibleengine_rds_instance_v3.test2", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "destination_db.0.region",
						"flexibleengine_rds_instance_v3.test2", "region"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "public_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "private_ip"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"source_db.0.password", "destination_db.0.password",
					"expired_days", "migrate_definer", "force_destroy"},
			},
		},
	})
}

func testAccDrsNet_base(name string) string {
	return fmt.Sprintf(`
data "flexibleengine_availability_zones" "test" {}

resource "flexibleengine_vpc_v1" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "test" {
  name          = "%s"
  cidr          = "192.168.0.0/24"
  gateway_ip    = "192.168.0.1"
  primary_dns   = "100.125.1.250"
  secondary_dns = "100.125.21.250"
  vpc_id        = flexibleengine_vpc_v1.test.id
}

resource "flexibleengine_networking_secgroup_v2" "test" {
  name = "%s"
}

resource "flexibleengine_networking_secgroup_rule_v2" "test" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 3306
  port_range_max    = 3306
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = flexibleengine_networking_secgroup_v2.test.id
}
`, name, name, name)
}

func testAccDrsJob_mysql(index int, name, pwd, ip string) string {
	return fmt.Sprintf(`
resource "flexibleengine_rds_instance_v3" "test%d" {
  name                = "%s%d"
  flavor              = "rds.mysql.s3.large.2.ha"
  security_group_id   = flexibleengine_networking_secgroup_v2.test.id
  subnet_id           = flexibleengine_vpc_subnet_v1.test.id
  vpc_id              = flexibleengine_vpc_v1.test.id
  fixed_ip            = "%s"
  ha_replication_mode = "semisync"

  availability_zone = [
    data.flexibleengine_availability_zones.test.names[0],
    data.flexibleengine_availability_zones.test.names[1],
  ]

  db {
    password = "%s"
    type     = "MySQL"
    version  = "5.7"
    port     = 3306
  }

  volume {
    type = "ULTRAHIGH"
    size = 40
  }

  depends_on = [
    flexibleengine_networking_secgroup_rule_v2.test
  ]
}
`, index, name, index, ip, pwd)
}

func testAccDrsJob_migrate_mysql(name, dbName, pwd string) string {
	netConfig := testAccDrsNet_base(name)
	sourceDb := testAccDrsJob_mysql(1, dbName, pwd, "192.168.0.58")
	destDb := testAccDrsJob_mysql(2, dbName, pwd, "192.168.0.59")

	return fmt.Sprintf(`
%s
%s
%s

resource "flexibleengine_drs_job" "test" {
  name           = "%s"
  type           = "migration"
  engine_type    = "mysql"
  direction      = "up"
  net_type       = "eip"
  migration_type = "FULL_INCR_TRANS"
  description    = "%s"
  force_destroy  = true

  source_db {
    engine_type = "mysql"
    ip          = flexibleengine_rds_instance_v3.test1.fixed_ip
    port        = 3306
    user        = "root"
    password    = "%s"
  }


  destination_db {
    region      = flexibleengine_rds_instance_v3.test2.region
    ip          = flexibleengine_rds_instance_v3.test2.fixed_ip
    port        = 3306
    engine_type = "mysql"
    user        = "root"
    password    = "%s"
    instance_id = flexibleengine_rds_instance_v3.test2.id
    subnet_id   = flexibleengine_rds_instance_v3.test2.subnet_id
  }

  lifecycle {
    ignore_changes = [
      source_db.0.password, destination_db.0.password, force_destroy,
    ]
  }
}
`, netConfig, sourceDb, destDb, name, name, pwd, pwd)
}
