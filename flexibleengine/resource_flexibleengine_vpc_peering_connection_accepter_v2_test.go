package flexibleengine

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"regexp"
)

func TestAccFlexibleEngineVpcPeeringConnectionAccepterV2_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFlexibleEngineVpcPeeringConnectionAccepterDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccFlexibleEngineVpcPeeringConnectionAccepterV2_basic, //TODO: Research why normal scenario with peer tenant id is not working in acceptance tests
				ExpectError: regexp.MustCompile(`VPC peering action not permitted: Can not accept/reject peering request not in PENDING_ACCEPTANCE state.`),
			},
		},
	})
}

func testAccCheckFlexibleEngineVpcPeeringConnectionAccepterDestroy(s *terraform.State) error {
	// We don't destroy the underlying VPC Peering Connection.
	return nil
}

const testAccFlexibleEngineVpcPeeringConnectionAccepterV2_basic = `
resource "flexibleengine_vpc_v1" "vpc_1" {
  name = "otc_vpc_1"
  cidr = "192.168.0.0/16"
}
resource "flexibleengine_vpc_v1" "vpc_2" {
  name = "otc_vpc_2"
  cidr = "192.168.0.0/16"
}
resource "flexibleengine_vpc_peering_connection_v2" "peering_1" {
    name = "flexibleengine"
    vpc_id = "${flexibleengine_vpc_v1.vpc_1.id}"
    peer_vpc_id = "${flexibleengine_vpc_v1.vpc_2.id}"
  }
resource "flexibleengine_vpc_peering_connection_accepter_v2" "peer" {
  vpc_peering_connection_id = "${flexibleengine_vpc_peering_connection_v2.peering_1.id}"
  accept = true

}
`
