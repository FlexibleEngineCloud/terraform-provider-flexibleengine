package acceptance

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDcsMaintainWindowDataSource_basic(t *testing.T) {
	sourceName := "data.flexibleengine_dcs_maintainwindow_v1.maintainwindow1"
	dc := acceptance.InitDataSourceCheck(sourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsMaintainWindowDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(sourceName, "seq", "1"),
					resource.TestMatchResourceAttr(sourceName, "begin", regexp.MustCompile(`^\d{2}$`)),
				),
			},
		},
	})
}

var testAccDcsMaintainWindowDataSource_basic = `
data "flexibleengine_dcs_maintainwindow_v1" "maintainwindow1" {
  seq = 1
}
`
