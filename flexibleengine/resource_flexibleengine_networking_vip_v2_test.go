package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/networking/v1/ports"
)

// TestAccNetworkingV2VIP_basic is basic acc test.
func TestAccNetworkingV2VIP_basic(t *testing.T) {
	var vip ports.Port
	resourceName := "flexibleengine_networking_vip_v2.vip_1"
	rName := fmt.Sprintf("tf_test_%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFloatingIP(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2VIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingVIP_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2VIPExists(resourceName, &vip),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "ip_version", "4"),
					resource.TestCheckResourceAttrSet(resourceName, "mac_address"),
				),
			},
			{
				Config: testAccNetworkingVIP_basic(rName + "_update"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName+"_update"),
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

// testAccCheckNetworkingV2VIPDestroy checks destory.
func testAccCheckNetworkingV2VIPDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	vpcClient, err := config.NetworkingV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine VPC v1 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_networking_vip_v2" {
			continue
		}

		_, err := ports.Get(vpcClient, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("VIP still exists")
		}
	}

	return nil
}

// testAccCheckNetworkingV2VIPExists checks exist.
func testAccCheckNetworkingV2VIPExists(n string, vip *ports.Port) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		vpcClient, err := config.NetworkingV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine VPC v1 client: %s", err)
		}

		found, err := ports.Get(vpcClient, rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("VIP not found")
		}

		*vip = *found
		return nil
	}
}

func testAccNetworkingVIP_basic(rName string) string {
	return fmt.Sprintf(`
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

resource "flexibleengine_networking_vip_v2" "vip_1" {
  name       = "%[1]s"
  network_id = flexibleengine_vpc_subnet_v1.subnet_1.id
}
`, rName)
}
