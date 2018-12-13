package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDcsMaintainWindowV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDcsMaintainWindowV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsMaintainWindowV1DataSourceID("data.flexibleengine_dcs_maintainwindow_v1.maintainwindow1"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_dcs_maintainwindow_v1.maintainwindow1", "seq", "1"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_dcs_maintainwindow_v1.maintainwindow1", "begin", "22"),
				),
			},
		},
	})
}

func testAccCheckDcsMaintainWindowV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Dcs maintainwindow data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Dcs maintainwindow data source ID not set")
		}

		return nil
	}
}

var testAccDcsMaintainWindowV1DataSource_basic = fmt.Sprintf(`
data "flexibleengine_dcs_maintainwindow_v1" "maintainwindow1" {
seq = 1
}
`)
