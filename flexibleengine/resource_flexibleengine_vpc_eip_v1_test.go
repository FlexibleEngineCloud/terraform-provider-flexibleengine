package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/golangsdk/openstack/networking/v1/eips"
)

func TestAccVpcV1EIP_basic(t *testing.T) {
	var eip eips.PublicIp
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "flexibleengine_vpc_eip_v1.eip_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcV1EIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcV1EIP_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcV1EIPExists(resourceName, &eip),
					resource.TestCheckResourceAttr(resourceName, "status", "UNBOUND"),
					resource.TestCheckResourceAttr(resourceName, "publicip.0.type", "5_bgp"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth.0.name", rName),
					resource.TestCheckResourceAttr(resourceName, "bandwidth.0.share_type", "PER"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth.0.charge_mode", "traffic"),
				),
			},
			{
				Config: testAccVpcV1EIP_tags(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcV1EIPExists(resourceName, &eip),
					resource.TestCheckResourceAttr(resourceName, "status", "UNBOUND"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
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

func testAccCheckVpcV1EIPDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.networkingV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating EIP: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_vpc_eip_v1" {
			continue
		}

		_, err := eips.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("EIP still exists")
		}
	}

	return nil
}

func testAccCheckVpcV1EIPExists(n string, kp *eips.PublicIp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.networkingV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating networking client: %s", err)
		}

		found, err := eips.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("EIP not found")
		}

		kp = &found

		return nil
	}
}

func testAccVpcV1EIP_basic(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_vpc_eip_v1" "eip_1" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    share_type = "PER"
    name       = "%s"
    size       = 5
  }
}
`, rName)
}

func testAccVpcV1EIP_tags(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_vpc_eip_v1" "eip_1" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    share_type = "PER"
    name       = "%s"
    size       = 5
  }
  tags = {
    foo = "bar"
    key = "value"
  }
}
`, rName)
}
