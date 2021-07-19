package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/golangsdk/openstack/rds/v3/instances"
)

func TestAccRdsInstanceV3_basic(t *testing.T) {
	var instance instances.RdsInstanceResponse
	name := acctest.RandString(4)
	resourceType := "flexibleengine_rds_instance_v3"
	resourceName := "flexibleengine_rds_instance_v3.instance"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRdsInstanceV3Destroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceV3_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("rds_acc_instance-%s", name)),
					resource.TestCheckResourceAttr(resourceName, "flavor", "rds.pg.s1.medium"),
					resource.TestCheckResourceAttr(resourceName, "time_zone", "UTC+01:00"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "5432"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "1"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "60"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttrSet(resourceName, "fixed_ip"),
				),
			},
			{
				Config: testAccRdsInstanceV3_update(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "flavor", "rds.pg.s1.large"),
					resource.TestCheckResourceAttr(resourceName, "time_zone", "UTC+01:00"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "5432"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "2"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "100"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"db",
					"status",
				},
			},
		},
	})
}

func testAccCheckRdsInstanceV3Destroy(rsType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		client, err := config.RdsV3Client(OS_REGION_NAME)
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

		config := testAccProvider.Meta().(*Config)
		client, err := config.RdsV3Client(OS_REGION_NAME)
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

func testAccRdsInstanceV3_basic(val string) string {
	return fmt.Sprintf(`
resource "flexibleengine_networking_secgroup_v2" "secgroup" {
  name        = "sg-acc-%s"
  description = "security group for rds instance"
}

resource "flexibleengine_rds_instance_v3" "instance" {
  name              = "rds_acc_instance-%s"
  flavor            = "rds.pg.s1.medium"
  availability_zone = ["%s"]
  security_group_id = flexibleengine_networking_secgroup_v2.secgroup.id
  vpc_id            = "%s"
  subnet_id         = "%s"
  time_zone         = "UTC+01:00"

  db {
    password = "Huangwei!120521"
    type     = "PostgreSQL"
    version  = "11"
  }
  volume {
    type = "COMMON"
    size = 60
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }

  tags = {
    key = "value"
    foo = "bar"
  }
}
	`, val, val, OS_AVAILABILITY_ZONE, OS_VPC_ID, OS_NETWORK_ID)
}

// volume, backup_strategy and tags will be updated
func testAccRdsInstanceV3_update(val string) string {
	return fmt.Sprintf(`
resource "flexibleengine_networking_secgroup_v2" "secgroup" {
  name        = "sg-acc-%s"
  description = "security group for rds instance"
}

resource "flexibleengine_rds_instance_v3" "instance" {
  name              = "rds_acc_instance-%s"
  flavor            = "rds.pg.s1.large"
  availability_zone = ["%s"]
  security_group_id = flexibleengine_networking_secgroup_v2.secgroup.id
  vpc_id            = "%s"
  subnet_id         = "%s"
  time_zone         = "UTC+01:00"

  db {
    password = "Huangwei!120521"
    type     = "PostgreSQL"
    version  = "11"
  }
  volume {
    type = "COMMON"
    size = 100
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 2
  }

  tags = {
    key   = "value1"
    owner = "terraform"
  }
}
	`, val, val, OS_AVAILABILITY_ZONE, OS_VPC_ID, OS_NETWORK_ID)
}
