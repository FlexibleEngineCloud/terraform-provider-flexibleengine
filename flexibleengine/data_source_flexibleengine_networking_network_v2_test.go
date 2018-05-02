package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// PASS
func TestAccFlexibleEngineNetworkingNetworkV2DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccFlexibleEngineNetworkingNetworkV2DataSource_network,
			},
			resource.TestStep{
				Config: testAccFlexibleEngineNetworkingNetworkV2DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkV2DataSourceID("data.flexibleengine_networking_network_v2.net"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_networking_network_v2.net", "name", "flexibleengine_test_network"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_networking_network_v2.net", "admin_state_up", "true"),
				),
			},
		},
	})
}

// PASS
func TestAccFlexibleEngineNetworkingNetworkV2DataSource_subnet(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccFlexibleEngineNetworkingNetworkV2DataSource_network,
			},
			resource.TestStep{
				Config: testAccFlexibleEngineNetworkingNetworkV2DataSource_subnet,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkV2DataSourceID("data.flexibleengine_networking_network_v2.net"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_networking_network_v2.net", "name", "flexibleengine_test_network"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_networking_network_v2.net", "admin_state_up", "true"),
				),
			},
		},
	})
}

// PASS
func TestAccFlexibleEngineNetworkingNetworkV2DataSource_networkID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccFlexibleEngineNetworkingNetworkV2DataSource_network,
			},
			resource.TestStep{
				Config: testAccFlexibleEngineNetworkingNetworkV2DataSource_networkID,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkV2DataSourceID("data.flexibleengine_networking_network_v2.net"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_networking_network_v2.net", "name", "flexibleengine_test_network"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_networking_network_v2.net", "admin_state_up", "true"),
				),
			},
		},
	})
}

func testAccCheckNetworkingNetworkV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find network data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Network data source ID not set")
		}

		return nil
	}
}

const testAccFlexibleEngineNetworkingNetworkV2DataSource_network = `
resource "flexibleengine_networking_network_v2" "net" {
        name = "flexibleengine_test_network"
        admin_state_up = "true"
}

resource "flexibleengine_networking_subnet_v2" "subnet" {
  name = "flexibleengine_test_subnet"
  cidr = "192.168.198.0/24"
  no_gateway = true
  network_id = "${flexibleengine_networking_network_v2.net.id}"
}
`

var testAccFlexibleEngineNetworkingNetworkV2DataSource_basic = fmt.Sprintf(`
%s

data "flexibleengine_networking_network_v2" "net" {
	name = "${flexibleengine_networking_network_v2.net.name}"
}
`, testAccFlexibleEngineNetworkingNetworkV2DataSource_network)

var testAccFlexibleEngineNetworkingNetworkV2DataSource_subnet = fmt.Sprintf(`
%s

data "flexibleengine_networking_network_v2" "net" {
	matching_subnet_cidr = "${flexibleengine_networking_subnet_v2.subnet.cidr}"
}
`, testAccFlexibleEngineNetworkingNetworkV2DataSource_network)

var testAccFlexibleEngineNetworkingNetworkV2DataSource_networkID = fmt.Sprintf(`
%s

data "flexibleengine_networking_network_v2" "net" {
	network_id = "${flexibleengine_networking_network_v2.net.id}"
}
`, testAccFlexibleEngineNetworkingNetworkV2DataSource_network)
