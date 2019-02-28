package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/layer3/routers"
)

func TestAccNetworkingV2Router_basic(t *testing.T) {
	var router routers.Router
	var routerName = fmt.Sprintf("terra-test-%s", acctest.RandString(5))
	var routerNameUpdate = fmt.Sprintf("%s-update", routerName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2RouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Router_basic(routerName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2RouterExists("flexibleengine_networking_router_v2.router_1", &router),
				),
			},
			{
				Config: testAccNetworkingV2Router_update(routerNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"flexibleengine_networking_router_v2.router_1", "name", routerNameUpdate),
				),
			},
		},
	})
}

func TestAccNetworkingV2Router_update_external_gw(t *testing.T) {
	var router routers.Router
	var routerName = fmt.Sprintf("terra-test-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2RouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Router_update_external_gw_1(routerName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2RouterExists("flexibleengine_networking_router_v2.router_1", &router),
				),
			},
			{
				Config: testAccNetworkingV2Router_update_external_gw_2(routerName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"flexibleengine_networking_router_v2.router_1", "external_gateway", OS_EXTGW_ID),
				),
			},
		},
	})
}

func TestAccNetworkingV2Router_timeout(t *testing.T) {
	var router routers.Router
	var routerName = fmt.Sprintf("terra-test-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2RouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Router_timeout(routerName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2RouterExists("flexibleengine_networking_router_v2.router_1", &router),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2RouterDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_networking_router_v2" {
			continue
		}

		_, err := routers.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Router still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingV2RouterExists(n string, router *routers.Router) resource.TestCheckFunc {
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

		found, err := routers.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Router not found")
		}

		*router = *found

		return nil
	}
}

func testAccNetworkingV2Router_basic(routerName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_networking_router_v2" "router_1" {
	name = "%s"
	admin_state_up = "true"
	distributed = "false"
}
`, routerName)
}

func testAccNetworkingV2Router_update(routerName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_networking_router_v2" "router_1" {
	name = "%s"
	admin_state_up = "true"
	distributed = "false"
}
`, routerName)
}

func testAccNetworkingV2Router_update_external_gw_1(routerName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_networking_router_v2" "router_1" {
	name = "%s"
	admin_state_up = "true"
	distributed = "false"
}
`, routerName)
}

func testAccNetworkingV2Router_update_external_gw_2(routerName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_networking_router_v2" "router_1" {
	name = "%s"
	admin_state_up = "true"
	distributed = "false"
	external_gateway = "%s"
}
`, OS_EXTGW_ID, routerName)
}

func testAccNetworkingV2Router_timeout(routerName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_networking_router_v2" "router_1" {
	name = "%s"
	admin_state_up = "true"
	distributed = "false"

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`, routerName)
}
