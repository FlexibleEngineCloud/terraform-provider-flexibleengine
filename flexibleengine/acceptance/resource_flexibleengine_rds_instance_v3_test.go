package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/rds/v3/instances"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/rds"
)

func TestAccRdsInstanceV3_basic(t *testing.T) {
	var instance instances.RdsInstanceResponse
	name := acctest.RandString(4)
	resourceType := "flexibleengine_rds_instance_v3"
	resourceName := "flexibleengine_rds_instance_v3.instance"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceV3Destroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceV3_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("rds_acc_instance-%s", name)),
					resource.TestCheckResourceAttr(resourceName, "flavor", "rds.pg.s3.medium.4"),
					resource.TestCheckResourceAttr(resourceName, "time_zone", "UTC+01:00"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "5432"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "1"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "60"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.name", "div_precision_increment"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.value", "12"),
					resource.TestCheckResourceAttrSet(resourceName, "fixed_ip"),
				),
			},
			{
				Config: testAccRdsInstanceV3_update(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tf_acc_instance-%s", name)),
					resource.TestCheckResourceAttr(resourceName, "flavor", "rds.pg.s3.large.4"),
					resource.TestCheckResourceAttr(resourceName, "time_zone", "UTC+01:00"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "5436"),
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
		conf := testAccProvider.Meta().(*config.Config)
		client, err := conf.RdsV3Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating flexibleengine rds client: %s", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != rsType {
				continue
			}

			id := rs.Primary.ID
			instance, err := rds.GetRdsInstanceByID(client, id)
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

		conf := testAccProvider.Meta().(*config.Config)
		client, err := conf.RdsV3Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating flexibleengine rds client: %s", err)
		}

		found, err := rds.GetRdsInstanceByID(client, id)
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

func testAccRdsInstanceV3_network(rName string) string {
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

`, rName)
}

func testAccRdsInstanceV3_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_rds_instance_v3" "instance" {
  name              = "rds_acc_instance-%s"
  flavor            = "rds.pg.s3.medium.4"
  availability_zone = ["%s"]
  security_group_id = flexibleengine_networking_secgroup_v2.secgroup.id
  vpc_id            = flexibleengine_vpc_v1.vpc_1.id
  subnet_id         = flexibleengine_vpc_subnet_v1.subnet_1.id
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

  parameters {
    name  = "div_precision_increment"
    value = "12"
  }

  tags = {
    key = "value"
    foo = "bar"
  }
}
`, testAccRdsInstanceV3_network(rName), rName, OS_AVAILABILITY_ZONE)
}

// name, flavor, port, volume, backup_strategy and tags will be updated
func testAccRdsInstanceV3_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_rds_instance_v3" "instance" {
  name              = "tf_acc_instance-%s"
  flavor            = "rds.pg.s3.large.4"
  availability_zone = ["%s"]
  security_group_id = flexibleengine_networking_secgroup_v2.secgroup.id
  vpc_id            = flexibleengine_vpc_v1.vpc_1.id
  subnet_id         = flexibleengine_vpc_subnet_v1.subnet_1.id
  time_zone         = "UTC+01:00"

  db {
    password = "Huangwei!120521"
    type     = "PostgreSQL"
    version  = "11"
    port     = 5436
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
`, testAccRdsInstanceV3_network(rName), rName, OS_AVAILABILITY_ZONE)
}
