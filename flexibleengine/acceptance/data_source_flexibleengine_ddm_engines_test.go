package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceDdmEngines_basic(t *testing.T) {
	rName := "data.flexibleengine_ddm_engines.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDdmEngines_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "engines.#", "1"),
					resource.TestCheckResourceAttr(rName, "engines.0.version", "3.0.8.5"),
				),
			},
		},
	})
}

func testAccDatasourceDdmEngines_basic() string {
	return `
data "flexibleengine_ddm_engines" "test" {
  version    = "3.0.8.5"
}
`
}
