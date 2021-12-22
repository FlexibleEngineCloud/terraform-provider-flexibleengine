package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDmsMaintainWindowDataSource_basic(t *testing.T) {
	dataSourceName := "data.flexibleengine_dms_maintainwindow.maintainwindow1"
	dc := initDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsMaintainWindowDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "seq", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "begin", "22:00:00"),
				),
			},
		},
	})
}

var testAccDmsMaintainWindowDataSource_basic = fmt.Sprintf(`
data "flexibleengine_dms_maintainwindow" "maintainwindow1" {
  seq = 1
}
`)
