package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccFlexibleEngineVpcRouteIdsV2DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFlexibleEngineRouteIdV2DataSource_vpcroute,
			},
			{
				Config: testAccFlexibleEngineRouteIdV2DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccFlexibleEngineRouteIdV2DataSourceID("data.flexibleengine_vpc_route_ids_v2.route_ids"),
					resource.TestCheckResourceAttr("data.flexibleengine_vpc_route_ids_v2.route_ids", "ids.#", "1"),
				),
			},
		},
	})
}
func testAccFlexibleEngineRouteIdV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find vpc route data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Vpc Route data source ID not set")
		}

		return nil
	}
}

const testAccFlexibleEngineRouteIdV2DataSource_vpcroute = `
resource "flexibleengine_vpc_v1" "vpc_1" {
name = "vpc_test"
cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_v1" "vpc_2" {
		name = "vpc_test1"
        cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_peering_connection_v2" "peering_1" {
		name = "flexibleengine_peering"
		vpc_id = "${flexibleengine_vpc_v1.vpc_1.id}"
		peer_vpc_id = "${flexibleengine_vpc_v1.vpc_2.id}"
}

resource "flexibleengine_vpc_route_v2" "route_1" {
   type = "peering"
  nexthop = "${flexibleengine_vpc_peering_connection_v2.peering_1.id}"
  destination = "192.168.0.0/16"
  vpc_id ="${flexibleengine_vpc_v1.vpc_1.id}"
}
`

var testAccFlexibleEngineRouteIdV2DataSource_basic = fmt.Sprintf(`
%s
data "flexibleengine_vpc_route_ids_v2" "route_ids" {
  vpc_id = "${flexibleengine_vpc_route_v2.route_1.vpc_id}"
}
`, testAccFlexibleEngineRouteIdV2DataSource_vpcroute)
