package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v1/ports"
)

// TestAccNetworkingV2VIPAssociate_basic is basic acc test.
func TestAccNetworkingV2VIPAssociate_basic(t *testing.T) {
	rName := fmt.Sprintf("tf_test_%s", acctest.RandString(5))
	resourceName := "flexibleengine_networking_vip_associate_v2.vip_associate_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFloatingIP(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2VIPAssociateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2VIPAssociate_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2VIPAssociateAssociated(),
					resource.TestCheckResourceAttrPair(resourceName, "vip_id",
						"flexibleengine_networking_vip_v2.vip_1", "id"),
				),
			},
		},
	})
}

// testAccCheckNetworkingV2VIPAssociateDestroy checks destory.
func testAccCheckNetworkingV2VIPAssociateDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	vpcClient, err := config.NetworkingV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine VPC client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_networking_vip_associate_v2" {
			continue
		}

		vipId, portIds, err := parseNetworkingVIPAssociateID(rs.Primary.ID)
		if err != nil {
			return err
		}

		vipPort, err := ports.Get(vpcClient, vipId)
		if err != nil {
			// If the error is a 404, then the vip port does not exist,
			// and therefore the floating IP cannot be associated to it.
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil
			}
			return err
		}

		// port by port
		for _, portId := range portIds {
			p, err := ports.Get(vpcClient, portId)
			if err != nil {
				// If the error is a 404, then the port does not exist,
				// and therefore the floating IP cannot be associated to it.
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					return nil
				}
				return err
			}

			// But if the port and vip still exists
			for _, ip := range p.FixedIps {
				for _, addresspair := range vipPort.AllowedAddressPairs {
					if ip.IpAddress == addresspair.IpAddress {
						return fmt.Errorf("VIP %s is still associated to port %s", vipId, portId)
					}
				}
			}
		}
	}

	return nil
}

func testAccCheckNetworkingV2VIPAssociateAssociated() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		vpcClient, err := config.NetworkingV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine VPC client: %s", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "flexibleengine_networking_vip_associate_v2" {
				continue
			}

			vipId, portIds, err := parseNetworkingVIPAssociateID(rs.Primary.ID)
			if err != nil {
				return err
			}

			vipPort, err := ports.Get(vpcClient, vipId)
			if err != nil {
				return err
			}

			// port by port
			for _, portId := range portIds {
				p, err := ports.Get(vpcClient, portId)
				if err != nil {
					return err
				}

				isAllowed := false
				for _, ip := range p.FixedIps {
					for _, addresspair := range vipPort.AllowedAddressPairs {
						if ip.IpAddress == addresspair.IpAddress {
							// port it associated
							isAllowed = true
							break
						}
					}
				}
				if !isAllowed {
					return fmt.Errorf("VIP %s was not attached to port %s", vipPort.ID, portId)
				}
			}
		}

		return nil
	}
}

func testAccNetworkingV2VIPAssociate_basic(rName string) string {
	return fmt.Sprintf(`
data "flexibleengine_images_image_v2" "ubuntu" {
  name = "OBS Ubuntu 20.04"
}

resource "flexibleengine_vpc_v1" "vpc_1" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "subnet_1" {
  name       = "%[1]s"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = flexibleengine_vpc_v1.vpc_1.id
}

resource "flexibleengine_compute_instance_v2" "servers" {
  count = 2

  name            = "instance_${count.index}"
  flavor_id       = "s6.small.1"
  image_id        = data.flexibleengine_images_image_v2.ubuntu.id
  security_groups = ["default"]
  network {
    uuid = flexibleengine_vpc_subnet_v1.subnet_1.id
  }
}

resource "flexibleengine_networking_vip_v2" "vip_1" {
  name       = "%[1]s"
  network_id = flexibleengine_vpc_subnet_v1.subnet_1.id
}

resource "flexibleengine_networking_vip_associate_v2" "vip_associate_1" {
  vip_id   = flexibleengine_networking_vip_v2.vip_1.id
  port_ids = [
    flexibleengine_compute_instance_v2.servers.0.network.0.port,
    flexibleengine_compute_instance_v2.servers.1.network.0.port,
  ]
}
`, rName)
}
