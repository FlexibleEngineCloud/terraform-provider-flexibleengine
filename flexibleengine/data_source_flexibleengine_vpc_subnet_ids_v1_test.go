package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccFlexibleEngineVpcSubnetIdsV2DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFlexibleEngineSubnetIdV2DataSource_vpcsubnet,
			},
			{
				Config: testAccFlexibleEngineSubnetIdV2DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccFlexibleEngineSubnetIdV2DataSourceID("data.flexibleengine_vpc_subnet_ids_v1.subnet_ids"),
					resource.TestCheckResourceAttr("data.flexibleengine_vpc_subnet_ids_v1.subnet_ids", "ids.#", "1"),
				),
			},
		},
	})
}
func testAccFlexibleEngineSubnetIdV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find vpc subnet data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Vpc Subnet data source ID not set")
		}

		return nil
	}
}

const testAccFlexibleEngineSubnetIdV2DataSource_vpcsubnet = `
resource "flexibleengine_vpc_v1" "vpc_1" {
	name = "test_vpc"
	cidr= "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "subnet_1" {
  name = "flexibleengine_subnet"
  cidr = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id = "${flexibleengine_vpc_v1.vpc_1.id}"
}
`

var testAccFlexibleEngineSubnetIdV2DataSource_basic = fmt.Sprintf(`
%s
data "flexibleengine_vpc_subnet_ids_v1" "subnet_ids" {
  vpc_id = "${flexibleengine_vpc_v1.vpc_1.id}"
}
`, testAccFlexibleEngineSubnetIdV2DataSource_vpcsubnet)
