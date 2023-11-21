package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDmsMaintainWindowDataSource_basic(t *testing.T) {
	dataSourceName := "data.flexibleengine_dms_maintainwindow.maintainwindow1"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

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

var testAccDmsMaintainWindowDataSource_basic = `
data "flexibleengine_dms_maintainwindow" "maintainwindow1" {
  seq = 1
}
`
