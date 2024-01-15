package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataPublicGateway_basic(t *testing.T) {
	var (
		name            = acceptance.RandomAccResourceName()
		nameFilter      = acceptance.InitDataSourceCheck("data.flexibleengine_nat_gateway_v2.name_filter")
		idFilter        = acceptance.InitDataSourceCheck("data.flexibleengine_nat_gateway_v2.id_filter")
		allParamsFilter = acceptance.InitDataSourceCheck("data.flexibleengine_nat_gateway_v2.all_params_filter")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataPublicGateway_basic(name),
				Check: resource.ComposeTestCheckFunc(
					nameFilter.CheckResourceExists(),
					idFilter.CheckResourceExists(),
					allParamsFilter.CheckResourceExists(),
				),
			},
		},
	})
}

func testAccDataPublicGateway_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "flexibleengine_nat_gateway_v2" "test" {
  name                  = "%[2]s"
  spec                  = "1"
  subnet_id             = flexibleengine_vpc_subnet_v1.test.id
  vpc_id                = flexibleengine_vpc_v1.test.id
}

data "flexibleengine_nat_gateway_v2" "name_filter" {
  name = flexibleengine_nat_gateway_v2.test.name
}

data "flexibleengine_nat_gateway_v2" "id_filter" {
  id = flexibleengine_nat_gateway_v2.test.id
}

data "flexibleengine_nat_gateway_v2" "all_params_filter" {
  name                  = flexibleengine_nat_gateway_v2.test.name
  spec                  = flexibleengine_nat_gateway_v2.test.spec
  subnet_id             = flexibleengine_nat_gateway_v2.test.subnet_id
  vpc_id                = flexibleengine_nat_gateway_v2.test.vpc_id
  enterprise_project_id = "0"
}
`, testBaseNetwork(name), name)
}
