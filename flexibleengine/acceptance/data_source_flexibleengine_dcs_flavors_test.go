package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDcsFlavors_basic(t *testing.T) {
	dataSourceName := "data.flexibleengine_dcs_flavors.flavors"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsFlavors_conf(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.#"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.engine", "redis"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.capacity", "0.125"),
				),
			},
		},
	})
}

func testAccDcsFlavors_conf() string {
	return `
data "flexibleengine_dcs_flavors" "flavors" {
  engine   = "Redis"
  capacity = 0.125
}
`
}
