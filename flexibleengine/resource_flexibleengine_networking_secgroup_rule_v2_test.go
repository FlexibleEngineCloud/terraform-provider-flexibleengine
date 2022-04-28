package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/security/groups"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/security/rules"
)

func TestAccNetworkingV2SecGroupRule_basic(t *testing.T) {
	var secgroup groups.SecGroup
	var secgroupRule rules.SecGroupRule

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "flexibleengine_networking_secgroup_rule_v2.secgroup_rule_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SecGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SecGroupRule_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SecGroupExists(
						"flexibleengine_networking_secgroup_v2.secgroup_1", &secgroup),
					testAccCheckNetworkingV2SecGroupRuleExists(resourceName, &secgroupRule),
					resource.TestCheckResourceAttr(resourceName, "direction", "ingress"),
					resource.TestCheckResourceAttr(resourceName, "port_range_min", "22"),
					resource.TestCheckResourceAttr(resourceName, "port_range_max", "22"),
					resource.TestCheckResourceAttr(resourceName, "ethertype", "IPv4"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "tcp"),
					resource.TestCheckResourceAttr(resourceName, "remote_ip_prefix", "0.0.0.0/0"),
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

func TestAccNetworkingV2SecGroupRule_remoteGroup(t *testing.T) {
	var secgroupRule rules.SecGroupRule

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "flexibleengine_networking_secgroup_rule_v2.secgroup_rule_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SecGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SecGroupRule_remoteGroup(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SecGroupRuleExists(resourceName, &secgroupRule),
					resource.TestCheckResourceAttr(resourceName, "direction", "ingress"),
					resource.TestCheckResourceAttr(resourceName, "port_range_min", "80"),
					resource.TestCheckResourceAttr(resourceName, "port_range_max", "80"),
					resource.TestCheckResourceAttr(resourceName, "ethertype", "IPv4"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "tcp"),
					resource.TestCheckResourceAttrSet(resourceName, "remote_group_id"),
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

func TestAccNetworkingV2SecGroupRule_ipv6(t *testing.T) {
	var secgroupRule rules.SecGroupRule

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "flexibleengine_networking_secgroup_rule_v2.secgroup_rule_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SecGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SecGroupRule_ipv6(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SecGroupRuleExists(resourceName, &secgroupRule),
					resource.TestCheckResourceAttr(resourceName, "remote_ip_prefix", "2001:558:fc00::/39"),
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

func TestAccNetworkingV2SecGroupRule_numericProtocol(t *testing.T) {
	var secgroupRule rules.SecGroupRule

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "flexibleengine_networking_secgroup_rule_v2.secgroup_rule_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SecGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SecGroupRule_numericProtocol(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SecGroupRuleExists(resourceName, &secgroupRule),
					resource.TestCheckResourceAttr(resourceName, "protocol", "115"),
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

func testAccCheckNetworkingV2SecGroupRuleDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.NetworkingV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_networking_secgroup_rule_v2" {
			continue
		}

		_, err := rules.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Security group rule still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingV2SecGroupRuleExists(n string, sgRule *rules.SecGroupRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.NetworkingV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine networking client: %s", err)
		}

		found, err := rules.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Security group rule not found")
		}

		*sgRule = *found

		return nil
	}
}

func testAccNetworkingV2SecGroupRule_base(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_networking_secgroup_v2" "secgroup_1" {
  name        = "secgroup-%s"
  description = "terraform security group rule acceptance test"
}
`, rName)
}

func testAccNetworkingV2SecGroupRule_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_networking_secgroup_rule_v2" "secgroup_rule_1" {
  direction         = "ingress"
  ethertype         = "IPv4"
  port_range_max    = 22
  port_range_min    = 22
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = flexibleengine_networking_secgroup_v2.secgroup_1.id
}
`, testAccNetworkingV2SecGroupRule_base(rName))
}

func testAccNetworkingV2SecGroupRule_remoteGroup(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_networking_secgroup_rule_v2" "secgroup_rule_1" {
  direction         = "ingress"
  ethertype         = "IPv4"
  port_range_max    = 80
  port_range_min    = 80
  protocol          = "tcp"
  remote_group_id   = flexibleengine_networking_secgroup_v2.secgroup_1.id
  security_group_id = flexibleengine_networking_secgroup_v2.secgroup_1.id
}
`, testAccNetworkingV2SecGroupRule_base(rName))
}

func testAccNetworkingV2SecGroupRule_ipv6(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_networking_secgroup_rule_v2" "secgroup_rule_1" {
  direction         = "ingress"
  ethertype         = "IPv6"
  port_range_max    = 22
  port_range_min    = 22
  protocol          = "tcp"
  remote_ip_prefix  = "2001:558:FC00::/39"
  security_group_id = flexibleengine_networking_secgroup_v2.secgroup_1.id
}
`, testAccNetworkingV2SecGroupRule_base(rName))
}

func testAccNetworkingV2SecGroupRule_numericProtocol(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_networking_secgroup_rule_v2" "secgroup_rule_1" {
  direction         = "ingress"
  ethertype         = "IPv4"
  port_range_max    = 22
  port_range_min    = 22
  protocol          = "115"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = flexibleengine_networking_secgroup_v2.secgroup_1.id
}
`, testAccNetworkingV2SecGroupRule_base(rName))
}
