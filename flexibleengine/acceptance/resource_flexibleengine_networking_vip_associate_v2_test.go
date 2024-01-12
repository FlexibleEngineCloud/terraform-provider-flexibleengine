package acceptance

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v2/ports"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccNetworkingV2VIPAssociate_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2VIPAssociateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2VIPAssociateConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("flexibleengine_networking_vip_associate_v2.vip_associate_1",
						"port_ids.0", "flexibleengine_compute_instance_v2.test", "network.0.port"),
					resource.TestCheckResourceAttrPair("flexibleengine_networking_vip_associate_v2.vip_associate_1",
						"vip_id", "flexibleengine_networking_vip_v2.vip_1", "id"),
				),
			},
			{
				ResourceName:      "flexibleengine_networking_vip_associate_v2.vip_associate_1",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccNetworkingV2VIPAssociateImportStateIdFunc(),
			},
		},
	})
}

func testAccCheckNetworkingV2VIPAssociateDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	networkingClient, err := config.NetworkingV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_networking_vip_associate_v2" {
			continue
		}

		vipID := rs.Primary.Attributes["vip_id"]
		_, err = ports.Get(networkingClient, vipID).Extract()
		if err != nil {
			// If the error is a 404, then the vip port does not exist,
			// and therefore the floating IP cannot be associated to it.
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil
			}
			return err
		}
	}

	log.Printf("[DEBUG] Destroy NetworkingVIPAssociated success!")
	return nil
}

func testAccNetworkingV2VIPAssociateImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		vip, ok := s.RootModule().Resources["flexibleengine_networking_vip_v2.vip_1"]
		if !ok {
			return "", fmt.Errorf("vip not found: %s", vip)
		}
		instance, ok := s.RootModule().Resources["flexibleengine_compute_instance_v2.test"]
		if !ok {
			return "", fmt.Errorf("port not found: %s", instance)
		}
		if vip.Primary.ID == "" || instance.Primary.Attributes["network.0.port"] == "" {
			return "", fmt.Errorf("resource not found: %s/%s", vip.Primary.ID,
				instance.Primary.Attributes["network.0.port"])
		}
		return fmt.Sprintf("%s/%s", vip.Primary.ID, instance.Primary.Attributes["network.0.port"]), nil
	}
}

const testAccCompute_data = `
data "flexibleengine_availability_zones" "test" {}

data "flexibleengine_compute_flavors_v2" "test" {
  availability_zone = data.flexibleengine_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core          = 2
  memory_size       = 4
}

data "flexibleengine_images_image" "test" {
  name = "OBS Ubuntu 20.04"
}

data "flexibleengine_networking_secgroup_v2" "test" {
  name = "default"
}

resource "flexibleengine_vpc_v1" "vpc_1" {
  name = "tf-test7bd"
  cidr = "192.168.0.0/24"
}

resource "flexibleengine_vpc_subnet_v1" "test" {
  name       = "tf-test7bd"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = flexibleengine_vpc_v1.vpc_1.id
}
`

func testAccComputeInstance_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_compute_instance_v2" "test" {
  name                = "%s"
  image_id            = data.flexibleengine_images_image.test.id
  flavor_id           = data.flexibleengine_compute_flavors_v2.test.flavors[0]
  security_groups     = [data.flexibleengine_networking_secgroup_v2.test.name]
  stop_before_destroy = true

  network {
    uuid = flexibleengine_vpc_subnet_v1.test.id
  }
}
`, testAccCompute_data, rName)
}

func testAccNetworkingV2VIPAssociateConfig_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_networking_port" "port" {
  port_id = flexibleengine_compute_instance_v2.test.network[0].port
}

resource "flexibleengine_networking_vip_v2" "vip_1" {
  network_id = flexibleengine_vpc_subnet_v1.test.id
}

resource "flexibleengine_networking_vip_associate_v2" "vip_associate_1" {
  vip_id   = flexibleengine_networking_vip_v2.vip_1.id
  port_ids = [flexibleengine_compute_instance_v2.test.network[0].port]
}
`, testAccComputeInstance_basic(rName))
}
