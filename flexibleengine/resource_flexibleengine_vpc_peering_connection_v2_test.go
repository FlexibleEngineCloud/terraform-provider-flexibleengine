package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/peerings"
)

func TestAccFlexibleEngineVpcPeeringConnectionV2_basic(t *testing.T) {
	var peering peerings.Peering

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFlexibleEngineVpcPeeringConnectionV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFlexibleEngineVpcPeeringConnectionV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlexibleEngineVpcPeeringConnectionV2Exists("flexibleengine_vpc_peering_connection_v2.peering_1", &peering),
					resource.TestCheckResourceAttr(
						"flexibleengine_vpc_peering_connection_v2.peering_1", "name", "flexibleengine_peering"),
					resource.TestCheckResourceAttr(
						"flexibleengine_vpc_peering_connection_v2.peering_1", "status", "ACTIVE"),
				),
			},
			{
				Config: testAccFlexibleEngineVpcPeeringConnectionV2_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"flexibleengine_vpc_peering_connection_v2.peering_1", "name", "flexibleengine_peering_1"),
				),
			},
		},
	})
}

func TestAccFlexibleEngineVpcPeeringConnectionV2_timeout(t *testing.T) {
	var peering peerings.Peering

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFlexibleEngineVpcPeeringConnectionV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFlexibleEngineVpcPeeringConnectionV2_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlexibleEngineVpcPeeringConnectionV2Exists("flexibleengine_vpc_peering_connection_v2.peering_1", &peering),
				),
			},
		},
	})
}

func testAccCheckFlexibleEngineVpcPeeringConnectionV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	peeringClient, err := config.networkingV2Client(OS_REGION_NAME)
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
		peeringClient, err := config.networkingV2Client(OS_REGION_NAME)
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
  name = "flexibleengine_peering"
  vpc_id = "${flexibleengine_vpc_v1.vpc_1.id}"
  peer_vpc_id = "${flexibleengine_vpc_v1.vpc_2.id}"
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
  name = "flexibleengine_peering_1"
  vpc_id = "${flexibleengine_vpc_v1.vpc_1.id}"
  peer_vpc_id = "${flexibleengine_vpc_v1.vpc_2.id}"
}
`
const testAccFlexibleEngineVpcPeeringConnectionV2_timeout = `
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

 timeouts {
    create = "5m"
    delete = "5m"
  }
}
`
