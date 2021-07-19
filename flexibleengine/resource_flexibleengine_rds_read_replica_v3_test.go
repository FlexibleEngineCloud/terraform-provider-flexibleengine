package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/golangsdk/openstack/rds/v3/instances"
)

func TestAccRdsReplicaInstanceV3_basic(t *testing.T) {
	var replica instances.RdsInstanceResponse
	resourceName := "flexibleengine_rds_read_replica_v3.replica_instance"
	resourceType := "flexibleengine_rds_read_replica_v3"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRdsInstanceV3Destroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsReplicaInstanceV3_basic(acctest.RandString(10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(resourceName, &replica),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "type", "Replica"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.type", "ULTRAHIGH"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "100"),
					// port of read replica is not same with port of rds instance
					//resource.TestCheckResourceAttr(resourceName, "db.0.port", "8635"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.func", "readonly"),
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

func testAccRdsReplicaInstanceV3_basic(val string) string {
	return fmt.Sprintf(`
resource "flexibleengine_networking_secgroup_v2" "secgroup" {
  name        = "sg-acc-%s"
  description = "security group for acceptance test"
}

resource "flexibleengine_rds_instance_v3" "instance" {
  name              = "rds_instance_%s"
  flavor            = "rds.pg.c2.large"
  availability_zone = ["%s"]
  security_group_id = flexibleengine_networking_secgroup_v2.secgroup.id
  vpc_id            = "%s"
  subnet_id         = "%s"

  db {
    password = "Huangwei!120521"
    type     = "PostgreSQL"
    version  = "11"
    port     = "8635"
  }
  volume {
    type = "ULTRAHIGH"
    size = 100
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }
  tags = {
    key  = "value"
    func = "readwrite"
  }
}

resource "flexibleengine_rds_read_replica_v3" "replica_instance" {
  name              = "replica_instance_%s"
  flavor            = "rds.pg.c2.large.rr"
  replica_of_id     = flexibleengine_rds_instance_v3.instance.id
  availability_zone = "%s"

  volume {
    type = "ULTRAHIGH"
  }
  tags = {
    key  = "value"
    func = "readonly"
  }
}

	`, val, val, OS_AVAILABILITY_ZONE, OS_VPC_ID, OS_NETWORK_ID, val, OS_AVAILABILITY_ZONE)
}
