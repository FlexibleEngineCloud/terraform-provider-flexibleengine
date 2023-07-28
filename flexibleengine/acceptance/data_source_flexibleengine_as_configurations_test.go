package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceASConfiguration_basic(t *testing.T) {
	dataSourceName := "data.flexibleengine_as_configurations.configurations"
	name := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceASConfiguration_conf(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "configurations.0.scaling_configuration_name", name),
				),
			},
		},
	})
}

func testAccDataSourceASConfiguration_conf(name string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_as_configurations" "configurations" {
  name     = flexibleengine_as_configuration.acc_as_config.scaling_configuration_name
  image_id = flexibleengine_as_configuration.acc_as_config.instance_config.0.image

  depends_on = [flexibleengine_as_configuration.acc_as_config]
}
`, testASGroup_Base(name))
}
