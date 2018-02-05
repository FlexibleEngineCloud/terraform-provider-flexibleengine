package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// PASS
func TestAccOrangeCloudNetworkingSecGroupV2DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOrangeCloudNetworkingSecGroupV2DataSource_group,
			},
			resource.TestStep{
				Config: testAccOrangeCloudNetworkingSecGroupV2DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupV2DataSourceID("data.flexibleengine_networking_secgroup_v2.secgroup_1"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_networking_secgroup_v2.secgroup_1", "name", "orange_acctest_secgroup"),
				),
			},
		},
	})
}

// PASS
func TestAccOrangeCloudNetworkingSecGroupV2DataSource_secGroupID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOrangeCloudNetworkingSecGroupV2DataSource_group,
			},
			resource.TestStep{
				Config: testAccOrangeCloudNetworkingSecGroupV2DataSource_secGroupID,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupV2DataSourceID("data.flexibleengine_networking_secgroup_v2.secgroup_1"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_networking_secgroup_v2.secgroup_1", "name", "orange_acctest_secgroup"),
				),
			},
		},
	})
}

func testAccCheckNetworkingSecGroupV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find security group data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Security group data source ID not set")
		}

		return nil
	}
}

const testAccOrangeCloudNetworkingSecGroupV2DataSource_group = `
resource "flexibleengine_networking_secgroup_v2" "secgroup_1" {
        name        = "orange_acctest_secgroup"
	description = "My neutron security group for orange acctest"
}
`

var testAccOrangeCloudNetworkingSecGroupV2DataSource_basic = fmt.Sprintf(`
%s

data "flexibleengine_networking_secgroup_v2" "secgroup_1" {
	name = "${flexibleengine_networking_secgroup_v2.secgroup_1.name}"
}
`, testAccOrangeCloudNetworkingSecGroupV2DataSource_group)

var testAccOrangeCloudNetworkingSecGroupV2DataSource_secGroupID = fmt.Sprintf(`
%s

data "flexibleengine_networking_secgroup_v2" "secgroup_1" {
	secgroup_id = "${flexibleengine_networking_secgroup_v2.secgroup_1.id}"
}
`, testAccOrangeCloudNetworkingSecGroupV2DataSource_group)
