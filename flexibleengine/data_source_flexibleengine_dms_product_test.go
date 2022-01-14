package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDmsProductDataSource_basic(t *testing.T) {
	dataSourceName := "data.flexibleengine_dms_product.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsProductDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsProductDataSourceID(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "engine", "kafka"),
					resource.TestCheckResourceAttr(dataSourceName, "partition_num", "300"),
					resource.TestCheckResourceAttr(dataSourceName, "storage_spec_codes.#", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "availability_zones.#", "3"),
					resource.TestCheckResourceAttrSet(dataSourceName, "cpu_arch"),
					resource.TestCheckResourceAttrSet(dataSourceName, "spec_code"),
				),
			},
		},
	})
}

func testAccCheckDmsProductDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find DMS product data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("DMS product data source ID not set")
		}

		return nil
	}
}

var testAccDmsProductDataSource_basic = fmt.Sprintf(`
data "flexibleengine_dms_product" "test" {
  bandwidth = "100MB"
}
`)
