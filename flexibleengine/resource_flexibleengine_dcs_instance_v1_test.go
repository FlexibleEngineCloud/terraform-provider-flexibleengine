package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/golangsdk/openstack/dcs/v1/instances"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccDcsInstancesV1_basic(t *testing.T) {
	var instance instances.Instance
	var instanceName = fmt.Sprintf("dcs_instance_%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDcsV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1Instance_basic(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsV1InstanceExists("flexibleengine_dcs_instance_v1.instance_1", instance),
					resource.TestCheckResourceAttr(
						"flexibleengine_dcs_instance_v1.instance_1", "name", instanceName),
					resource.TestCheckResourceAttr(
						"flexibleengine_dcs_instance_v1.instance_1", "engine", "Redis"),
				),
			},
		},
	})
}

func testAccCheckDcsV1InstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	dcsClient, err := config.dcsV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating Flexibleengine instance client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_dcs_instance_v1" {
			continue
		}

		_, err := instances.Get(dcsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("The Dcs instance still exists.")
		}
	}
	return nil
}

func testAccCheckDcsV1InstanceExists(n string, instance instances.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		dcsClient, err := config.dcsV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating Flexibleengine instance client: %s", err)
		}

		v, err := instances.Get(dcsClient, rs.Primary.ID).Extract()
		if err != nil {
			return fmt.Errorf("Error getting Flexibleengine instance: %s, err: %s", rs.Primary.ID, err)
		}

		if v.InstanceID != rs.Primary.ID {
			return fmt.Errorf("The Dcs instance not found.")
		}
		instance = *v
		return nil
	}
}

func testAccDcsV1Instance_basic(instanceName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_networking_secgroup_v2" "secgroup_1" {
  name        = "secgroup_1"
  description = "secgroup_1"
}

resource "flexibleengine_vpc_v1" "vpc_1" {
  name = "terraform_vpc1"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "subnet_1" {
  name       = "terraform_subnet"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = flexibleengine_vpc_v1.vpc_1.id
}

resource "flexibleengine_dcs_instance_v1" "instance_1" {
  name              = "%s"
  engine            = "Redis"
  engine_version    = "3.0"
  password          = "Huawei_test"
  product_id        = "dcs.master_standby-h"
  capacity          = 2
  vpc_id            = flexibleengine_vpc_v1.vpc_1.id
  network_id        = flexibleengine_vpc_subnet_v1.subnet_1.id
  security_group_id = flexibleengine_networking_secgroup_v2.secgroup_1.id
  available_zones   = ["eu-west-0a"]

  save_days   = 1
  backup_type = "manual"
  begin_at    = "00:00-01:00"
  period_type = "weekly"
  backup_at   = [1]
}
	`, instanceName)
}
