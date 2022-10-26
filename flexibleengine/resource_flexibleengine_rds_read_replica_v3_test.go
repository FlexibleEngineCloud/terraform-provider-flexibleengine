package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/rds/v3/instances"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccRdsReplicaInstanceV3_basic(t *testing.T) {
	var replica instances.RdsInstanceResponse
	resourceName := "flexibleengine_rds_read_replica_v3.replica_instance"
	resourceType := "flexibleengine_rds_read_replica_v3"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRdsInstanceV3Destroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsReplicaInstanceV3_basic(acctest.RandString(4)),
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

func testAccCheckRdsInstanceV3Destroy(rsType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conf := testAccProvider.Meta().(*Config)
		client, err := conf.RdsV3Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating flexibleengine rds client: %s", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != rsType {
				continue
			}

			id := rs.Primary.ID
			instance, err := getRdsInstanceByID(client, id)
			if err != nil {
				return err
			}
			if instance.Id != "" {
				return fmt.Errorf("%s (%s) still exists", rsType, id)
			}
		}
		return nil
	}
}

func testAccCheckRdsInstanceV3Exists(name string, instance *instances.RdsInstanceResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		id := rs.Primary.ID
		if id == "" {
			return fmt.Errorf("No ID is set")
		}

		conf := testAccProvider.Meta().(*Config)
		client, err := conf.RdsV3Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating flexibleengine rds client: %s", err)
		}

		found, err := getRdsInstanceByID(client, id)
		if err != nil {
			return fmt.Errorf("Error checking %s exist, err=%s", name, err)
		}
		if found.Id == "" {
			return fmt.Errorf("resource %s does not exist", name)
		}

		instance = found
		return nil
	}
}

func testAccRdsReplicaInstanceV3_basic(val string) string {
	return fmt.Sprintf(`
resource "flexibleengine_vpc_v1" "vpc_1" {
  name = "vpc-acc-%[1]s"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "subnet_1" {
  name       = "subnet-acc-%[1]s"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = flexibleengine_vpc_v1.vpc_1.id
}

resource "flexibleengine_networking_secgroup_v2" "secgroup" {
  name        = "sg-acc-%[1]s"
  description = "security group for rds instance"
}

resource "flexibleengine_rds_instance_v3" "instance" {
  name              = "rds_instance_%[1]s"
  flavor            = "rds.pg.s3.large.2"
  availability_zone = ["%[2]s"]
  security_group_id = flexibleengine_networking_secgroup_v2.secgroup.id
  vpc_id            = flexibleengine_vpc_v1.vpc_1.id
  subnet_id         = flexibleengine_vpc_subnet_v1.subnet_1.id

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
  name              = "replica_instance_%[1]s"
  flavor            = "rds.pg.s3.large.2.rr"
  replica_of_id     = flexibleengine_rds_instance_v3.instance.id
  availability_zone = "%[2]s"

  volume {
    type = "ULTRAHIGH"
  }
  tags = {
    key  = "value"
    func = "readonly"
  }
}
`, val, OS_AVAILABILITY_ZONE)
}
