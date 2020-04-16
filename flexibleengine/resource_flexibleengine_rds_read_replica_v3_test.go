package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/rds/v3/instances"
)

func TestAccRdsReplicaInstanceV3_basic(t *testing.T) {
	var resourceMap string = "flexibleengine_rds_read_replica_v3.replica_instance"
	var replica instances.RdsInstanceResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRdsReplicaInstanceV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsReplicaInstanceV3_basic(acctest.RandString(10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsReplicaInstanceV3Exists(resourceMap, &replica),
					resource.TestCheckResourceAttr(resourceMap, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceMap, "type", "Replica"),
					resource.TestCheckResourceAttr(resourceMap, "volume.0.type", "ULTRAHIGH"),
					resource.TestCheckResourceAttr(resourceMap, "volume.0.size", "100"),
					// port of read replica is not same with port of rds instance
					//resource.TestCheckResourceAttr(resourceMap, "db.0.port", "8635"),
				),
			},
		},
	})
}

func testAccCheckRdsReplicaInstanceV3Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.rdsV3Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine rds client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_rds_read_replica_v3" {
			continue
		}

		id := rs.Primary.ID
		instance, err := getRdsInstanceByID(client, id)
		if err != nil {
			return err
		}
		if instance.Id != "" {
			return fmt.Errorf("flexibleengine_rds_read_replica_v3 (%s) still exists", id)
		}
	}

	return nil
}

func testAccCheckRdsReplicaInstanceV3Exists(n string, instance *instances.RdsInstanceResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s. ", n)
		}

		id := rs.Primary.ID
		if id == "" {
			return fmt.Errorf("No ID is set. ")
		}

		config := testAccProvider.Meta().(*Config)
		client, err := config.rdsV3Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine rds client: %s", err)
		}

		found, err := getRdsInstanceByID(client, id)
		if err != nil {
			return fmt.Errorf("Error checking %s exist, err=%s", n, err)
		}
		if found.Id == "" {
			return fmt.Errorf("resource %s does not exist", n)
		}

		instance = found
		return nil
	}
}

func testAccRdsReplicaInstanceV3_basic(val string) string {
	return fmt.Sprintf(`
resource "flexibleengine_networking_secgroup_v2" "secgroup" {
  name = "acctest_sg"
  description = "security group for acceptance test"
}

resource "flexibleengine_rds_instance_v3" "instance" {
  name = "rds_instance_%s"
  flavor = "rds.pg.c2.large"
  availability_zone = ["%s"]
  security_group_id = flexibleengine_networking_secgroup_v2.secgroup.id
  vpc_id = "%s"
  subnet_id = "%s"

  db {
    password = "Huangwei!120521"
    type = "PostgreSQL"
    version = "9.5.5"
    port = "8635"
  }
  volume {
    type = "ULTRAHIGH"
    size = 100
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days = 1
  }
}

resource "flexibleengine_rds_read_replica_v3" "replica_instance" {
  name = "replica_instance_%s"
  flavor = "rds.pg.c2.large.rr"
  replica_of_id = flexibleengine_rds_instance_v3.instance.id
  availability_zone = "%s"

  volume {
    type = "ULTRAHIGH"
  }
}

	`, val, OS_AVAILABILITY_ZONE, OS_VPC_ID, OS_NETWORK_ID, val, OS_AVAILABILITY_ZONE)
}
