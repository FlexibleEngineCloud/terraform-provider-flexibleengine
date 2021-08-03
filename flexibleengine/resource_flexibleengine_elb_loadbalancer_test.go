package flexibleengine

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/elbaas/loadbalancer_elbs"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/security/groups"
)

func TestAccELBLoadBalancer_basic(t *testing.T) {

	var lb loadbalancer_elbs.LoadBalancer

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckELBLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccELBLoadBalancerConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckELBLoadBalancerExists("flexibleengine_elb_loadbalancer.loadbalancer_1", &lb),
				),
			},
			{
				Config: testAccELBLoadBalancerConfig_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"flexibleengine_elb_loadbalancer.loadbalancer_1", "name", "loadbalancer_1_updated"),
				),
			},
		},
	})
}

func TestAccELBLoadBalancer_secGroup(t *testing.T) {
	var lb loadbalancer_elbs.LoadBalancer
	var sg_1, sg_2 groups.SecGroup

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckELBLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccELBLoadBalancer_secGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckELBLoadBalancerExists(
						"flexibleengine_elb_loadbalancer.loadbalancer_1", &lb),
					testAccCheckNetworkingV2SecGroupExists(
						"flexibleengine_networking_secgroup_v2.secgroup_1", &sg_1),
					testAccCheckNetworkingV2SecGroupExists(
						"flexibleengine_networking_secgroup_v2.secgroup_2", &sg_2),
					testAccCheckELBLoadBalancerHasSecGroup(&lb, &sg_1),
				),
			},
			{
				Config: testAccELBLoadBalancer_secGroup_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckELBLoadBalancerExists(
						"flexibleengine_elb_loadbalancer.loadbalancer_1", &lb),
					testAccCheckELBLoadBalancerHasSecGroup(&lb, &sg_2),
				),
			},
		},
	})
}

func testAccCheckELBLoadBalancerDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.otcV1Client(OS_REGION_NAME)
	if err != nil {
		fmt.Printf("@@@@@@@@@@@@@@@@ testAccCheckELBLoadBalancerDestroy FlexibleEngine networking client: %s", err)

		return fmt.Errorf("Error creating FlexibleEngine networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_elb_loadbalancer" {
			continue
		}

		_, err := loadbalancer_elbs.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			fmt.Printf("@@@@@@@@@@@@@@@@ testAccCheckELBLoadBalancerDestroy LoadBalancer still exists: %s", rs.Primary.ID)

			return fmt.Errorf("LoadBalancer still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckELBLoadBalancerExists(
	n string, lb *loadbalancer_elbs.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			fmt.Printf("@@@@@@@@@@@@@@@@ testAccCheckELBLoadBalancerExists Not found: %s \n", n)

			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			fmt.Printf("@@@@@@@@@@@@@@@@ testAccCheckELBLoadBalancerExists No ID is set \n")
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.otcV1Client(OS_REGION_NAME)
		if err != nil {
			fmt.Printf("@@@@@@@@@@@@@@@@ testAccCheckELBLoadBalancerExists Error creating FlexibleEngine networking client: %s", err)
			return fmt.Errorf("Error creating FlexibleEngine networking client: %s", err)
		}
		fmt.Printf("@@@@@@@@@@@@@@@@ testAccCheckELBLoadBalancerExists  middle \n ")
		found, err := loadbalancer_elbs.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			log.Printf("[#####ERR#####] : %v", err)

			fmt.Printf("@@@@@@@@@@@@@@@@ testAccCheckELBLoadBalancerExists err1 =%v\n ", err)
			return err
		}

		if found.ID != rs.Primary.ID {
			fmt.Printf("@@@@@@@@@@@@@@@@ testAccCheckELBLoadBalancerExists err2 Member not found \n ")

			return fmt.Errorf("Member not found")
		}

		*lb = *found

		return nil
	}
}

func testAccCheckELBLoadBalancerHasSecGroup(
	lb *loadbalancer_elbs.LoadBalancer, sg *groups.SecGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		_, err := config.otcV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine networking client: %s", err)
		}

		return nil
	}
}

var testAccELBLoadBalancerConfig_basic = fmt.Sprintf(`
resource "flexibleengine_elb_loadbalancer" "loadbalancer_1" {
  name = "loadbalancer_1"
  vpc_id = "%s"
  type = "External"
  bandwidth = 5

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, OS_VPC_ID)

var testAccELBLoadBalancerConfig_update = fmt.Sprintf(`
resource "flexibleengine_elb_loadbalancer" "loadbalancer_1" {
  name = "loadbalancer_1_updated"
  admin_state_up = "true"
  vpc_id = "%s"
  type = "External"
  bandwidth = 3

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, OS_VPC_ID)

var testAccELBLoadBalancer_secGroup = fmt.Sprintf(`
resource "flexibleengine_networking_secgroup_v2" "secgroup_1" {
  name = "secgroup_1"
  description = "secgroup_1"
}

resource "flexibleengine_networking_secgroup_v2" "secgroup_2" {
  name = "secgroup_2"
  description = "secgroup_2"
}

resource "flexibleengine_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "flexibleengine_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  network_id = "${flexibleengine_networking_network_v2.network_1.id}"
  cidr = "192.168.199.0/24"
}

resource "flexibleengine_elb_loadbalancer" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = "${flexibleengine_networking_subnet_v2.subnet_1.id}"
  vpc_id = "%s"
  type = "External"
  bandwidth = 3
  security_group_id = "${flexibleengine_networking_secgroup_v2.secgroup_1.id}"
}
`, OS_VPC_ID)

var testAccELBLoadBalancer_secGroup_update = fmt.Sprintf(`
resource "flexibleengine_networking_secgroup_v2" "secgroup_1" {
  name = "secgroup_1"
  description = "secgroup_1"
}

resource "flexibleengine_networking_secgroup_v2" "secgroup_2" {
  name = "secgroup_2"
  description = "secgroup_2"
}

resource "flexibleengine_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "flexibleengine_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  network_id = "${flexibleengine_networking_network_v2.network_1.id}"
  cidr = "192.168.199.0/24"
}

resource "flexibleengine_elb_loadbalancer" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = "${flexibleengine_networking_subnet_v2.subnet_1.id}"
  vpc_id = "%s"
  type = "External"
  bandwidth = 3
  security_group_id = "${flexibleengine_networking_secgroup_v2.secgroup_2.id}"
}
`, OS_VPC_ID)
