package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/golangsdk/openstack/networking/v2/networks"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/ports"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/subnets"
)

func TestAccNetworkingV2Port_basic(t *testing.T) {
	var network networks.Network
	var port ports.Port
	var subnet subnets.Subnet
	resourceName := "flexibleengine_networking_port_v2.port_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Port_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("flexibleengine_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("flexibleengine_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists(resourceName, &port),
					resource.TestCheckResourceAttr(resourceName, "name", "port_1"),
					resource.TestCheckResourceAttr(resourceName, "admin_state_up", "true"),
					resource.TestCheckResourceAttr(resourceName, "fixed_ip.#", "1"),
					resource.TestCheckResourceAttrPtr(resourceName, "network_id", &network.ID),
				),
			},
		},
	})
}

func TestAccNetworkingV2Port_fixedip(t *testing.T) {
	var network networks.Network
	var port ports.Port
	var subnet subnets.Subnet
	resourceName := "flexibleengine_networking_port_v2.port_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Port_fixedip,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("flexibleengine_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("flexibleengine_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists(resourceName, &port),
					resource.TestCheckResourceAttr(resourceName, "name", "port_1"),
					resource.TestCheckResourceAttr(resourceName, "admin_state_up", "true"),
					resource.TestCheckResourceAttr(resourceName, "fixed_ip.#", "1"),
					resource.TestCheckResourceAttrPtr(resourceName, "network_id", &network.ID),
				),
			},
			{
				Config: testAccNetworkingV2Port_fixedip_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "fixed_ip.#", "3"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Port_allowedAddressPairs(t *testing.T) {
	var network networks.Network
	var subnet subnets.Subnet
	var vrrp_port_1, vrrp_port_2, instance_port ports.Port

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Port_allowedAddressPairs,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("flexibleengine_networking_subnet_v2.vrrp_subnet", &subnet),
					testAccCheckNetworkingV2NetworkExists("flexibleengine_networking_network_v2.vrrp_network", &network),
					testAccCheckNetworkingV2PortExists("flexibleengine_networking_port_v2.vrrp_port_1", &vrrp_port_1),
					testAccCheckNetworkingV2PortExists("flexibleengine_networking_port_v2.vrrp_port_2", &vrrp_port_2),
					testAccCheckNetworkingV2PortExists("flexibleengine_networking_port_v2.instance_port", &instance_port),
					resource.TestCheckResourceAttr("flexibleengine_networking_port_v2.instance_port", "allowed_address_pairs.#", "2"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Port_securityGroup(t *testing.T) {
	var network networks.Network
	var port ports.Port
	var subnet subnets.Subnet
	resourceName := "flexibleengine_networking_port_v2.port_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Port_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("flexibleengine_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("flexibleengine_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists(resourceName, &port),
					resource.TestCheckResourceAttr(resourceName, "name", "port_1"),
					resource.TestCheckResourceAttr(resourceName, "admin_state_up", "true"),
					resource.TestCheckResourceAttr(resourceName, "fixed_ip.#", "1"),
					// default security group
					resource.TestCheckResourceAttr(resourceName, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttrPtr(resourceName, "network_id", &network.ID),
				),
			},
			{
				Config: testAccNetworkingV2Port_securityGroups,
				Check: resource.ComposeTestCheckFunc(
					// user defined security group
					resource.TestCheckResourceAttr(resourceName, "security_group_ids.#", "1"),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2PortDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_networking_port_v2" {
			continue
		}

		_, err := ports.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Port still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingV2PortExists(n string, port *ports.Port) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine networking client: %s", err)
		}

		found, err := ports.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Port not found")
		}

		*port = *found

		return nil
	}
}

const testAccNetworkingV2Port_preCondition string = `
resource "flexibleengine_networking_network_v2" "network_1" {
  name           = "network_1"
  admin_state_up = true
}

resource "flexibleengine_networking_subnet_v2" "subnet_1" {
  name       = "subnet_1"
  cidr       = "192.168.199.0/24"
  ip_version = 4
  network_id = flexibleengine_networking_network_v2.network_1.id
}
`

var testAccNetworkingV2Port_basic string = fmt.Sprintf(`
%s

resource "flexibleengine_networking_port_v2" "port_1" {
  name           = "port_1"
  admin_state_up = true
  network_id     = flexibleengine_networking_network_v2.network_1.id

  fixed_ip {
    subnet_id  = flexibleengine_networking_subnet_v2.subnet_1.id
    ip_address = "192.168.199.23"
  }
}
`, testAccNetworkingV2Port_preCondition)

var testAccNetworkingV2Port_fixedip string = fmt.Sprintf(`
%s

resource "flexibleengine_networking_port_v2" "port_1" {
  name           = "port_1"
  admin_state_up = true
  network_id     = flexibleengine_networking_network_v2.network_1.id

  fixed_ip {
    subnet_id = flexibleengine_networking_subnet_v2.subnet_1.id
  }
}
`, testAccNetworkingV2Port_preCondition)

var testAccNetworkingV2Port_fixedip_update string = fmt.Sprintf(`
%s

resource "flexibleengine_networking_port_v2" "port_1" {
  name           = "port_1"
  admin_state_up = true
  network_id     = flexibleengine_networking_network_v2.network_1.id

  fixed_ip {
    subnet_id  = flexibleengine_networking_subnet_v2.subnet_1.id
    ip_address = "192.168.199.20"
  }

  fixed_ip {
    subnet_id  = flexibleengine_networking_subnet_v2.subnet_1.id
    ip_address = "192.168.199.23"
  }

  fixed_ip {
    subnet_id  = flexibleengine_networking_subnet_v2.subnet_1.id
    ip_address = "192.168.199.40"
  }
}
`, testAccNetworkingV2Port_preCondition)

const testAccNetworkingV2Port_allowedAddressPairs = `
resource "flexibleengine_networking_network_v2" "vrrp_network" {
  name           = "vrrp_network"
  admin_state_up = true
}

resource "flexibleengine_networking_subnet_v2" "vrrp_subnet" {
  name       = "vrrp_subnet"
  cidr       = "10.0.0.0/24"
  ip_version = 4
  network_id = flexibleengine_networking_network_v2.vrrp_network.id

  allocation_pools {
    start = "10.0.0.2"
    end   = "10.0.0.200"
  }
}

resource "flexibleengine_networking_router_v2" "vrrp_router" {
  name = "vrrp_router"
}

resource "flexibleengine_networking_router_interface_v2" "vrrp_interface" {
  router_id = flexibleengine_networking_router_v2.vrrp_router.id
  subnet_id = flexibleengine_networking_subnet_v2.vrrp_subnet.id
}

resource "flexibleengine_networking_port_v2" "vrrp_port_1" {
  name           = "vrrp_port_1"
  admin_state_up = true
  network_id     = flexibleengine_networking_network_v2.vrrp_network.id

  fixed_ip {
    subnet_id  = flexibleengine_networking_subnet_v2.vrrp_subnet.id
    ip_address = "10.0.0.202"
  }
}

resource "flexibleengine_networking_port_v2" "vrrp_port_2" {
  name           = "vrrp_port_2"
  admin_state_up = true
  network_id     = flexibleengine_networking_network_v2.vrrp_network.id

  fixed_ip {
    subnet_id  = flexibleengine_networking_subnet_v2.vrrp_subnet.id
    ip_address = "10.0.0.201"
  }
}

resource "flexibleengine_networking_port_v2" "instance_port" {
  name           = "instance_port"
  admin_state_up = true
  network_id     = flexibleengine_networking_network_v2.vrrp_network.id

  allowed_address_pairs {
    ip_address  = tolist(flexibleengine_networking_port_v2.vrrp_port_1.fixed_ip).0.ip_address
    mac_address = flexibleengine_networking_port_v2.vrrp_port_1.mac_address
  }

  allowed_address_pairs {
    ip_address  = tolist(flexibleengine_networking_port_v2.vrrp_port_2.fixed_ip).0.ip_address
    mac_address = flexibleengine_networking_port_v2.vrrp_port_2.mac_address
  }
}
`

var testAccNetworkingV2Port_securityGroups = fmt.Sprintf(`
%s

resource "flexibleengine_networking_secgroup_v2" "secgroup_1" {
  name        = "security_group"
  description = "terraform security group acceptance test"
}

resource "flexibleengine_networking_port_v2" "port_1" {
  name               = "port_1"
  admin_state_up     = true
  network_id         = flexibleengine_networking_network_v2.network_1.id
  security_group_ids = [flexibleengine_networking_secgroup_v2.secgroup_1.id]

  fixed_ip {
    subnet_id  = flexibleengine_networking_subnet_v2.subnet_1.id
    ip_address = "192.168.199.23"
  }
}
`, testAccNetworkingV2Port_preCondition)
