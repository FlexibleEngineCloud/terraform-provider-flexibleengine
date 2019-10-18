package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/networking/v1/subnets"
)

func TestAccFlexibleEngineVpcSubnetV1_basic(t *testing.T) {
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFlexibleEngineVpcSubnetV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFlexibleEngineVpcSubnetV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlexibleEngineVpcSubnetV1Exists("flexibleengine_vpc_subnet_v1.subnet_1", &subnet),
					resource.TestCheckResourceAttr(
						"flexibleengine_vpc_subnet_v1.subnet_1", "name", "flexibleengine_subnet"),
					resource.TestCheckResourceAttr(
						"flexibleengine_vpc_subnet_v1.subnet_1", "cidr", "192.168.0.0/16"),
					resource.TestCheckResourceAttr(
						"flexibleengine_vpc_subnet_v1.subnet_1", "gateway_ip", "192.168.0.1"),
				),
			},
			{
				Config: testAccFlexibleEngineVpcSubnetV1_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"flexibleengine_vpc_subnet_v1.subnet_1", "name", "flexibleengine_subnet_1"),
				),
			},
		},
	})
}

func TestAccFlexibleEngineVpcSubnetV1_timeout(t *testing.T) {
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFlexibleEngineVpcSubnetV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFlexibleEngineVpcSubnetV1_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlexibleEngineVpcSubnetV1Exists("flexibleengine_vpc_subnet_v1.subnet_1", &subnet),
				),
			},
		},
	})
}

func testAccCheckFlexibleEngineVpcSubnetV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	subnetClient, err := config.networkingV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine vpc client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_vpc_subnet_v1" {
			continue
		}

		_, err := subnets.Get(subnetClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Subnet still exists")
		}
	}

	return nil
}
func testAccCheckFlexibleEngineVpcSubnetV1Exists(n string, subnet *subnets.Subnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		subnetClient, err := config.networkingV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine Vpc client: %s", err)
		}

		found, err := subnets.Get(subnetClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Subnet not found")
		}

		*subnet = *found

		return nil
	}
}

const testAccFlexibleEngineVpcSubnetV1_basic = `
resource "flexibleengine_vpc_v1" "vpc_1" {
  name = "vpc_test"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "subnet_1" {
  name = "flexibleengine_subnet"
  cidr = "192.168.0.0/16"
  gateway_ip = "192.168.0.1"
  vpc_id = "${flexibleengine_vpc_v1.vpc_1.id}"


}
`
const testAccFlexibleEngineVpcSubnetV1_update = `
resource "flexibleengine_vpc_v1" "vpc_1" {
  name = "vpc_test"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "subnet_1" {
  name = "flexibleengine_subnet_1"
  cidr = "192.168.0.0/16"
  gateway_ip = "192.168.0.1"
  vpc_id = "${flexibleengine_vpc_v1.vpc_1.id}"

 }
`

const testAccFlexibleEngineVpcSubnetV1_timeout = `
resource "flexibleengine_vpc_v1" "vpc_1" {
  name = "vpc_test"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "subnet_1" {
  name = "flexibleengine_subnet"
  cidr = "192.168.0.0/16"
  gateway_ip = "192.168.0.1"
  vpc_id = "${flexibleengine_vpc_v1.vpc_1.id}"

 timeouts {
    create = "5m"
    delete = "5m"
  }

}
`
