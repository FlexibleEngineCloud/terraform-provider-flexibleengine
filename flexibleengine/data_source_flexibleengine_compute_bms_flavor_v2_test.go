package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccBMSV2FlavorDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccBmsFlavorPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBMSV2FlavorDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBMSV2FlavorDataSourceID("data.flexibleengine_compute_bms_flavors_v2.byName"),
					testAccCheckBMSV2FlavorDataSourceID("data.flexibleengine_compute_bms_flavors_v2.byCPU"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_compute_bms_flavors_v2.byName", "name", OS_BMS_FLAVOR_NAME),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_compute_bms_flavors_v2.byCPU", "vcpus", "32"),
				),
			},
		},
	})
}

func testAccCheckBMSV2FlavorDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Flavor data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Flavor data source ID not set")
		}

		return nil
	}
}

var testAccBMSV2FlavorDataSource_basic = fmt.Sprintf(`
data "flexibleengine_compute_bms_flavors_v2" "byName" {
  name = "%s"
}

data "flexibleengine_compute_bms_flavors_v2" "byCPU" {
  vcpus = 32
}
`, OS_BMS_FLAVOR_NAME)
