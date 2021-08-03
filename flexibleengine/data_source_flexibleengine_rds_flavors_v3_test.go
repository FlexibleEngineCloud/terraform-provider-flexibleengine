package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccRdsFlavorV3DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsFlavorV3DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsFlavorV3DataSourceID("data.flexibleengine_rds_flavors_v3.flavor"),
				),
			},
		},
	})
}

func testAccCheckRdsFlavorV3DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find rds data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Rds data source ID not set ")
		}

		return nil
	}
}

var testAccRdsFlavorV3DataSource_basic = `

data "flexibleengine_rds_flavors_v3" "flavor" {
  db_type = "PostgreSQL"
  db_version = "9.5"
  instance_mode = "ha"
}
`
