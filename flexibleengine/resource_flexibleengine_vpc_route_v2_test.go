package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/networking/v2/routes"
)

func TestAccFlexibleEngineVpcRouteV2_basic(t *testing.T) {
	var route routes.Route
	resourceName := "flexibleengine_vpc_route_v2.route_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFlexibleEngineRouteV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlexibleEngineRouteV2Exists(resourceName, &route),
					resource.TestCheckResourceAttr(resourceName, "destination", "192.168.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "type", "peering"),
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

func testAccCheckFlexibleEngineRouteV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	routeClient, err := config.NetworkingV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine route client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_vpc_route_v2" {
			continue
		}

		_, err := routes.Get(routeClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Route still exists")
		}
	}

	return nil
}

func testAccCheckFlexibleEngineRouteV2Exists(n string, route *routes.Route) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		routeClient, err := config.NetworkingV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine route client: %s", err)
		}

		found, err := routes.Get(routeClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.RouteID != rs.Primary.ID {
			return fmt.Errorf("route not found")
		}

		*route = *found

		return nil
	}
}

const testAccRouteV2_basic = `
resource "flexibleengine_vpc_v1" "vpc_1" {
  name = "vpc_test"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_v1" "vpc_2" {
  name = "vpc_test1"
  cidr = "192.168.0.0/16"
}
resource "flexibleengine_vpc_peering_connection_v2" "peering_1" {
  name        = "flexibleengine_peering"
  vpc_id      = flexibleengine_vpc_v1.vpc_1.id
  peer_vpc_id = flexibleengine_vpc_v1.vpc_2.id
}
resource "flexibleengine_vpc_route_v2" "route_1" {
  type        = "peering"
  nexthop     = flexibleengine_vpc_peering_connection_v2.peering_1.id
  destination = "192.168.0.0/16"
  vpc_id      = flexibleengine_vpc_v1.vpc_1.id

}
`
