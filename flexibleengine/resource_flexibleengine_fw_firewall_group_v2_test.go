package flexibleengine

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/fwaas_v2/firewall_groups"
)

func TestAccFWFirewallGroupV2_basic(t *testing.T) {
	var epolicyID *string
	var ipolicyID *string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWFirewallGroupV2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccFWFirewallGroupV2_basic_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWFirewallGroupV2("flexibleengine_fw_firewall_group_v2.fw_1", "", "", ipolicyID, epolicyID),
				),
			},
			resource.TestStep{
				Config: testAccFWFirewallGroupV2_basic_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWFirewallGroupV2(
						"flexibleengine_fw_firewall_group_v2.fw_1", "fw_1", "terraform acceptance test", ipolicyID, epolicyID),
				),
			},
		},
	})
}

func TestAccFWFirewallGroupV2_port0(t *testing.T) {
	var firewall_group FirewallGroup

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWFirewallGroupV2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccFWFirewallV2_port,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWFirewallGroupV2Exists("flexibleengine_fw_firewall_group_v2.fw_1", &firewall_group),
					testAccCheckFWFirewallPortCount(&firewall_group, 1),
				),
			},
		},
	})
}

func TestAccFWFirewallGroupV2_no_ports(t *testing.T) {
	var firewall_group FirewallGroup

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWFirewallGroupV2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccFWFirewallV2_no_ports,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWFirewallGroupV2Exists("flexibleengine_fw_firewall_group_v2.fw_1", &firewall_group),
					resource.TestCheckResourceAttr("flexibleengine_fw_firewall_group_v2.fw_1", "description", "firewall router test"),
					testAccCheckFWFirewallPortCount(&firewall_group, 0),
				),
			},
		},
	})
}

func TestAccFWFirewallGroupV2_port_update(t *testing.T) {
	var firewall_group FirewallGroup

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWFirewallGroupV2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccFWFirewallV2_port,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWFirewallGroupV2Exists("flexibleengine_fw_firewall_group_v2.fw_1", &firewall_group),
					testAccCheckFWFirewallPortCount(&firewall_group, 1),
				),
			},
			resource.TestStep{
				Config: testAccFWFirewallV2_port_add,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWFirewallGroupV2Exists("flexibleengine_fw_firewall_group_v2.fw_1", &firewall_group),
					testAccCheckFWFirewallPortCount(&firewall_group, 2),
				),
			},
		},
	})
}

func TestAccFWFirewallGroupV2_port_remove(t *testing.T) {
	var firewall_group FirewallGroup

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWFirewallGroupV2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccFWFirewallV2_port,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWFirewallGroupV2Exists("flexibleengine_fw_firewall_group_v2.fw_1", &firewall_group),
					testAccCheckFWFirewallPortCount(&firewall_group, 1),
				),
			},
			resource.TestStep{
				Config: testAccFWFirewallV2_port_remove,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWFirewallGroupV2Exists("flexibleengine_fw_firewall_group_v2.fw_1", &firewall_group),
					testAccCheckFWFirewallPortCount(&firewall_group, 0),
				),
			},
		},
	})
}

func testAccCheckFWFirewallGroupV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.hwNetworkV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine networking client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_firewall_group" {
			continue
		}

		_, err = firewall_groups.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Firewall group (%s) still exists.", rs.Primary.ID)
		}
		if _, ok := err.(golangsdk.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckFWFirewallGroupV2Exists(n string, firewall_group *FirewallGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.hwNetworkV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Exists) Error creating FlexibleEngine networking client: %s", err)
		}

		var found FirewallGroup
		err = firewall_groups.Get(networkingClient, rs.Primary.ID).ExtractInto(&found)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Firewall group not found")
		}

		*firewall_group = found

		return nil
	}
}

func testAccCheckFWFirewallPortCount(firewall_group *FirewallGroup, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(firewall_group.PortIDs) != expected {
			return fmt.Errorf("Expected %d Ports, got %d", expected, len(firewall_group.PortIDs))
		}

		return nil
	}
}

func testAccCheckFWFirewallGroupV2(n, expectedName, expectedDescription string, ipolicyID *string, epolicyID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.hwNetworkV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Exists) Error creating FlexibleEngine networking client: %s", err)
		}

		var found *firewall_groups.FirewallGroup
		for i := 0; i < 5; i++ {
			// Firewall creation is asynchronous. Retry some times
			// if we get a 404 error. Fail on any other error.
			found, err = firewall_groups.Get(networkingClient, rs.Primary.ID).Extract()
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					time.Sleep(time.Second)
					continue
				}
				return err
			}
			break
		}

		switch {
		case found.Name != expectedName:
			err = fmt.Errorf("Expected Name to be <%s> but found <%s>", expectedName, found.Name)
		case found.Description != expectedDescription:
			err = fmt.Errorf("Expected Description to be <%s> but found <%s>",
				expectedDescription, found.Description)
		case found.IngressPolicyID == "":
			err = fmt.Errorf("Ingress Policy should not be empty")
		case found.EgressPolicyID == "":
			err = fmt.Errorf("Egress Policy should not be empty")
		case ipolicyID != nil && found.IngressPolicyID == *ipolicyID:
			err = fmt.Errorf("Ingress Policy had not been correctly updated. Went from <%s> to <%s>",
				expectedName, found.Name)
		case epolicyID != nil && found.EgressPolicyID == *epolicyID:
			err = fmt.Errorf("Egress Policy had not been correctly updated. Went from <%s> to <%s>",
				expectedName, found.Name)
		}

		if err != nil {
			return err
		}

		ipolicyID = &found.IngressPolicyID
		epolicyID = &found.EgressPolicyID

		return nil
	}
}

const testAccFWFirewallGroupV2_basic_1 = `
resource "flexibleengine_fw_firewall_group_v2" "fw_1" {
  ingress_policy_id = "${flexibleengine_fw_policy_v2.policy_1.id}"
  egress_policy_id = "${flexibleengine_fw_policy_v2.policy_1.id}"

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}

resource "flexibleengine_fw_policy_v2" "policy_1" {
  name = "policy_1"
}
`

const testAccFWFirewallGroupV2_basic_2 = `
resource "flexibleengine_fw_firewall_group_v2" "fw_1" {
  name = "fw_1"
  description = "terraform acceptance test"
  ingress_policy_id = "${flexibleengine_fw_policy_v2.policy_2.id}"
  egress_policy_id = "${flexibleengine_fw_policy_v2.policy_2.id}"
  admin_state_up = true

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}

resource "flexibleengine_fw_policy_v2" "policy_2" {
  name = "policy_2"
}
`

var testAccFWFirewallV2_port = fmt.Sprintf(`
resource "flexibleengine_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "flexibleengine_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  enable_dhcp = true
  network_id = "${flexibleengine_networking_network_v2.network_1.id}"
}

resource "flexibleengine_networking_router_v2" "router_1" {
  name = "router_1"
  admin_state_up = "true"
  external_gateway = "%s"
}

resource "flexibleengine_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${flexibleengine_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${flexibleengine_networking_subnet_v2.subnet_1.id}"
    #ip_address = "192.168.199.23"
  }
}

resource "flexibleengine_networking_router_interface_v2" "router_interface_1" {
  router_id = "${flexibleengine_networking_router_v2.router_1.id}"
  port_id = "${flexibleengine_networking_port_v2.port_1.id}"
}

resource "flexibleengine_fw_policy_v2" "policy_1" {
  name = "policy_1"
}

resource "flexibleengine_fw_firewall_group_v2" "fw_1" {
  name = "firewall_1"
  description = "firewall router test"
  ingress_policy_id = "${flexibleengine_fw_policy_v2.policy_1.id}"
  #egress_policy_id = "${flexibleengine_fw_policy_v2.policy_1.id}"
  ports = [
	"${flexibleengine_networking_port_v2.port_1.id}"
  ]
  depends_on = ["flexibleengine_networking_router_interface_v2.router_interface_1"]
}
`, OS_EXTGW_ID)

var testAccFWFirewallV2_port_add = fmt.Sprintf(`
resource "flexibleengine_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "flexibleengine_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${flexibleengine_networking_network_v2.network_1.id}"
}

resource "flexibleengine_networking_router_v2" "router_1" {
  name = "router_1"
  admin_state_up = "true"
  external_gateway = "%s"
}

resource "flexibleengine_networking_router_v2" "router_2" {
  name = "router_2"
  admin_state_up = "true"
  external_gateway = "%s"
}

resource "flexibleengine_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${flexibleengine_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${flexibleengine_networking_subnet_v2.subnet_1.id}"
    #ip_address = "192.168.199.23"
  }
}

resource "flexibleengine_networking_port_v2" "port_2" {
  name = "port_2"
  admin_state_up = "true"
  network_id = "${flexibleengine_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${flexibleengine_networking_subnet_v2.subnet_1.id}"
    #ip_address = "192.168.199.24"
  }
}

resource "flexibleengine_networking_router_interface_v2" "router_interface_1" {
  router_id = "${flexibleengine_networking_router_v2.router_1.id}"
  port_id = "${flexibleengine_networking_port_v2.port_1.id}"
}

resource "flexibleengine_networking_router_interface_v2" "router_interface_2" {
  router_id = "${flexibleengine_networking_router_v2.router_2.id}"
  port_id = "${flexibleengine_networking_port_v2.port_2.id}"
}

resource "flexibleengine_fw_policy_v2" "policy_1" {
  name = "policy_1"
}

resource "flexibleengine_fw_firewall_group_v2" "fw_1" {
  name = "firewall_1"
  description = "firewall router test"
  ingress_policy_id = "${flexibleengine_fw_policy_v2.policy_1.id}"
  egress_policy_id = "${flexibleengine_fw_policy_v2.policy_1.id}"
  ports = [
	"${flexibleengine_networking_port_v2.port_1.id}",
	"${flexibleengine_networking_port_v2.port_2.id}"
  ]
  depends_on = ["flexibleengine_networking_router_interface_v2.router_interface_1", "flexibleengine_networking_router_interface_v2.router_interface_2"]
}
`, OS_EXTGW_ID, OS_EXTGW_ID)

const testAccFWFirewallV2_port_remove = `
resource "flexibleengine_fw_policy_v2" "policy_1" {
  name = "policy_1"
}

resource "flexibleengine_fw_firewall_group_v2" "fw_1" {
  name = "firewall_1"
  description = "firewall router test"
  ingress_policy_id = "${flexibleengine_fw_policy_v2.policy_1.id}"
  egress_policy_id = "${flexibleengine_fw_policy_v2.policy_1.id}"
}
`

const testAccFWFirewallV2_no_ports = `
resource "flexibleengine_fw_policy_v2" "policy_1" {
  name = "policy_1"
}

resource "flexibleengine_fw_firewall_group_v2" "fw_1" {
  name = "firewall_1"
  description = "firewall router test"
  ingress_policy_id = "${flexibleengine_fw_policy_v2.policy_1.id}"
  egress_policy_id = "${flexibleengine_fw_policy_v2.policy_1.id}"
}
`
