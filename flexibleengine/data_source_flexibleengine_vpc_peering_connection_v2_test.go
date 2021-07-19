package flexibleengine

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

func TestAccFlexibleEngineVpcPeeringConnectionV2DataSource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceFlexibleEngineVpcPeeringConnectionV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlexibleEngineVpcPeeringConnectionV2DataSourceID("data.flexibleengine_vpc_peering_connection_v2.by_id"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_vpc_peering_connection_v2.by_id", "name", "flexibleengine_peering"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_vpc_peering_connection_v2.by_id", "status", "ACTIVE"),
					testAccCheckFlexibleEngineVpcPeeringConnectionV2DataSourceID("data.flexibleengine_vpc_peering_connection_v2.by_vpc_id"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_vpc_peering_connection_v2.by_vpc_id", "name", "flexibleengine_peering"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_vpc_peering_connection_v2.by_vpc_id", "status", "ACTIVE"),
					testAccCheckFlexibleEngineVpcPeeringConnectionV2DataSourceID("data.flexibleengine_vpc_peering_connection_v2.by_peer_vpc_id"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_vpc_peering_connection_v2.by_peer_vpc_id", "name", "flexibleengine_peering"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_vpc_peering_connection_v2.by_peer_vpc_id", "status", "ACTIVE"),
					testAccCheckFlexibleEngineVpcPeeringConnectionV2DataSourceID("data.flexibleengine_vpc_peering_connection_v2.by_vpc_ids"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_vpc_peering_connection_v2.by_vpc_ids", "name", "flexibleengine_peering"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_vpc_peering_connection_v2.by_vpc_ids", "status", "ACTIVE"),
				),
			},
		},
	})
}

func testAccCheckFlexibleEngineVpcPeeringConnectionV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find vpc peering connection data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("vpc peering connection data source ID not set")
		}

		return nil
	}
}

const testAccDataSourceFlexibleEngineVpcPeeringConnectionV2Config = `
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

data "flexibleengine_vpc_peering_connection_v2" "by_id" {
		id = "${flexibleengine_vpc_peering_connection_v2.peering_1.id}"
}

data "flexibleengine_vpc_peering_connection_v2" "by_vpc_id" {
		vpc_id = "${flexibleengine_vpc_peering_connection_v2.peering_1.vpc_id}"
}

data "flexibleengine_vpc_peering_connection_v2" "by_peer_vpc_id" {
		peer_vpc_id = "${flexibleengine_vpc_peering_connection_v2.peering_1.peer_vpc_id}"
}

data "flexibleengine_vpc_peering_connection_v2" "by_vpc_ids" {
		vpc_id = "${flexibleengine_vpc_peering_connection_v2.peering_1.vpc_id}"
		peer_vpc_id = "${flexibleengine_vpc_peering_connection_v2.peering_1.peer_vpc_id}"
}
`
