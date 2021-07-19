package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/golangsdk/openstack/compute/v2/extensions/attachinterfaces"
)

func TestAccComputeV2InterfaceAttach_basic(t *testing.T) {
	var ai attachinterfaces.Interface

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InterfaceAttachDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InterfaceAttach_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InterfaceAttachExists("flexibleengine_compute_interface_attach_v2.ai_1", &ai),
				),
			},
		},
	})
}

func TestAccComputeV2InterfaceAttach_IP(t *testing.T) {
	var ai attachinterfaces.Interface

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InterfaceAttachDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InterfaceAttach_IP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InterfaceAttachExists("flexibleengine_compute_interface_attach_v2.ai_1", &ai),
					testAccCheckComputeV2InterfaceAttachIP(&ai, "192.168.1.100"),
				),
			},
		},
	})
}

func testAccCheckComputeV2InterfaceAttachDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	computeClient, err := config.computeV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_compute_interface_attach_v2" {
			continue
		}

		instanceId, portId, err := computeInterfaceAttachV2ParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = attachinterfaces.Get(computeClient, instanceId, portId).Extract()
		if err == nil {
			return fmt.Errorf("Volume attachment still exists")
		}
	}

	return nil
}

func testAccCheckComputeV2InterfaceAttachExists(n string, ai *attachinterfaces.Interface) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		computeClient, err := config.computeV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine compute client: %s", err)
		}

		instanceId, portId, err := computeInterfaceAttachV2ParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		found, err := attachinterfaces.Get(computeClient, instanceId, portId).Extract()
		if err != nil {
			return err
		}

		//if found.instanceID != instanceID || found.PortID != portId {
		if found.PortID != portId {
			return fmt.Errorf("InterfaceAttach not found")
		}

		*ai = *found

		return nil
	}
}

func testAccCheckComputeV2InterfaceAttachIP(
	ai *attachinterfaces.Interface, ip string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, i := range ai.FixedIPs {
			if i.IPAddress == ip {
				return nil
			}
		}
		return fmt.Errorf("Requested ip (%s) does not exist on port", ip)

	}
}

var testAccComputeV2InterfaceAttach_basic = fmt.Sprintf(`
resource "flexibleengine_networking_port_v2" "port_1" {
  name = "port_1"
  network_id = "%s"
  admin_state_up = "true"
}

resource "flexibleengine_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
}

resource "flexibleengine_compute_interface_attach_v2" "ai_1" {
  instance_id = "${flexibleengine_compute_instance_v2.instance_1.id}"
  port_id = "${flexibleengine_networking_port_v2.port_1.id}"
}
`, OS_NETWORK_ID, OS_NETWORK_ID)

var testAccComputeV2InterfaceAttach_IP = fmt.Sprintf(`
resource "flexibleengine_networking_network_v2" "network_1" {
  name = "network_1"
}

resource "flexibleengine_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  network_id = "${flexibleengine_networking_network_v2.network_1.id}"
  cidr = "192.168.1.0/24"
  ip_version = 4
  enable_dhcp = true
  no_gateway = true
}

resource "flexibleengine_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
}

resource "flexibleengine_compute_interface_attach_v2" "ai_1" {
  instance_id = "${flexibleengine_compute_instance_v2.instance_1.id}"
  network_id = "${flexibleengine_networking_network_v2.network_1.id}"
  fixed_ip = "192.168.1.100"
}
`, OS_NETWORK_ID)
