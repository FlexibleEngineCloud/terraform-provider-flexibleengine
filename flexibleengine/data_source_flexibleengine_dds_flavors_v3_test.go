package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDDSFlavorsV3DataSource_basic(t *testing.T) {
	resourceName := "data.flexibleengine_dds_flavors_v3.flavor"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDDSFlavorsV3DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDDSFlavorsV3DataSourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "engine_name", "DDS-Community"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.#"),
				),
			},
		},
	})
}

func testAccCheckDDSFlavorsV3DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find DDS Flavor data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("DDS Flavor data source ID not set")
		}

		return nil
	}
}

var testAccDDSFlavorsV3DataSource_basic = `
data "flexibleengine_dds_flavors_v3" "flavor" {
  vcpus  = 8
  memory = 32
}
`
