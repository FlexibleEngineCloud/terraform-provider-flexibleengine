package acceptance

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccRabbitMQFlavorsDataSource_basic(t *testing.T) {
	dataSourceName := "data.flexibleengine_dms_rabbitmq_flavors.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDmsRabbitMQFlavors_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "versions.#", regexp.MustCompile(`[1-9]\d*`)),
					resource.TestMatchResourceAttr(dataSourceName, "flavors.#", regexp.MustCompile(`[1-9]\d*`)),
					resource.TestCheckOutput("type_validation", "true"),
					resource.TestCheckOutput("arch_types_validation", "true"),
					resource.TestCheckOutput("storage_spec_code_validation", "true"),
					resource.TestCheckOutput("availability_zones_validation", "true"),
				),
			},
		},
	})
}

func testAccDatasourceDmsRabbitMQFlavors_basic() string {
	return `
data "flexibleengine_dms_rabbitmq_flavors" "baisc" {
  type = "cluster"
}

data "flexibleengine_dms_rabbitmq_flavors" "test" {
  type               = local.test_refer.type
  arch_type          = local.test_refer.arch_types[0]
  storage_spec_code  = local.test_refer.ios[0].storage_spec_code
  availability_zones = local.test_refer.ios[0].availability_zones
}

locals {
  test_refer   = data.flexibleengine_dms_rabbitmq_flavors.baisc.flavors[0]
  test_results = data.flexibleengine_dms_rabbitmq_flavors.test
}

output "type_validation" {
  value = contains(local.test_results.flavors[*].type, local.test_refer.type)
}

output "arch_types_validation" {
  value = alltrue([for a in local.test_results.flavors[*].arch_types : contains(a, local.test_refer.arch_types[0])])
}

output "storage_spec_code_validation" {
  value = alltrue([for ios in local.test_results.flavors[*].ios :
  alltrue([for io in ios : io.storage_spec_code == local.test_refer.ios[0].storage_spec_code])])
}

output "availability_zones_validation" {
  value = alltrue([for ios in local.test_results.flavors[*].ios :
  alltrue([for io in ios : length(setintersection(io.availability_zones,
  local.test_refer.ios[0].availability_zones))== length(local.test_refer.ios[0].availability_zones)])])
}
`
}
