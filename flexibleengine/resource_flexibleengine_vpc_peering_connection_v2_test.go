package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/networking/v2/peerings"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccFlexibleEngineVpcPeeringConnectionV2_basic(t *testing.T) {
	var peering peerings.Peering
	resourceName := "flexibleengine_vpc_peering_connection_v2.peering_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFlexibleEngineVpcPeeringConnectionV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFlexibleEngineVpcPeeringConnectionV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlexibleEngineVpcPeeringConnectionV2Exists(resourceName, &peering),
					resource.TestCheckResourceAttr(resourceName, "name", "flexibleengine_peering"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
				),
			},
			{
				Config: testAccFlexibleEngineVpcPeeringConnectionV2_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "flexibleengine_peering_1"),
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

func testAccCheckFlexibleEngineVpcPeeringConnectionV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	peeringClient, err := config.NetworkingV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine Peering client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_vpc_peering_connection_v2" {
			continue
		}

		_, err := peerings.Get(peeringClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Vpc Peering Connection still exists")
		}
	}

	return nil
}

func testAccCheckFlexibleEngineVpcPeeringConnectionV2Exists(n string, peering *peerings.Peering) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		peeringClient, err := config.NetworkingV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine Peering client: %s", err)
		}

		found, err := peerings.Get(peeringClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Vpc peering Connection not found")
		}

		*peering = *found

		return nil
	}
}

const testAccFlexibleEngineVpcPeeringConnectionV2_basic = `
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
`
const testAccFlexibleEngineVpcPeeringConnectionV2_update = `
resource "flexibleengine_vpc_v1" "vpc_1" {
  name = "vpc_test"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_v1" "vpc_2" {
  name = "vpc_test1"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_peering_connection_v2" "peering_1" {
  name        = "flexibleengine_peering_1"
  vpc_id      = flexibleengine_vpc_v1.vpc_1.id
  peer_vpc_id = flexibleengine_vpc_v1.vpc_2.id
}
`
