package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccTurbosDataSource_basic(t *testing.T) {
	var (
		rName         = acceptance.RandomAccResourceNameWithDash()
		dcByName      = acceptance.InitDataSourceCheck("data.flexibleengine_sfs_turbos.by_name")
		dcBySize      = acceptance.InitDataSourceCheck("data.flexibleengine_sfs_turbos.by_size")
		dcByShareType = acceptance.InitDataSourceCheck("data.flexibleengine_sfs_turbos.by_share_type")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTurbosDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("name_query_result_validation", "true"),
					dcBySize.CheckResourceExists(),
					resource.TestCheckOutput("size_query_result_validation", "true"),
					dcByShareType.CheckResourceExists(),
					resource.TestCheckOutput("share_type_query_result_validation", "true"),
				),
			},
		},
	})
}

func testAccTurbosDataSource_basic(rName string) string {
	return fmt.Sprintf(`
variable "turbo_configuration" {
  type = list(object({
    size        = number
    share_type  = string
  }))

  default = [
    {size = 500, share_type = "PERFORMANCE"},
    {size = 600, share_type = "STANDARD"},
    {size = 600, share_type = "PERFORMANCE"},
  ]
}

data "flexibleengine_availability_zones" "test" {}

resource "flexibleengine_vpc_v1" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "test" {
  vpc_id = flexibleengine_vpc_v1.test.id

  name       = "%[1]s"
  cidr       = cidrsubnet(flexibleengine_vpc_v1.test.cidr, 4, 1)
  gateway_ip = cidrhost(cidrsubnet(flexibleengine_vpc_v1.test.cidr, 4, 1), 1)
}

resource "flexibleengine_networking_secgroup_v2" "test" {
  name = "%[1]s"
}

resource "flexibleengine_sfs_turbo" "test" {
  count = length(var.turbo_configuration)

  vpc_id            = flexibleengine_vpc_v1.test.id
  subnet_id         = flexibleengine_vpc_subnet_v1.test.id
  security_group_id = flexibleengine_networking_secgroup_v2.test.id
  availability_zone = data.flexibleengine_availability_zones.test.names[0]

  name        = "%[1]s-${count.index}"
  size        = var.turbo_configuration[count.index]["size"]
  share_proto = "NFS"
  share_type  = var.turbo_configuration[count.index]["share_type"]
}

data "flexibleengine_sfs_turbos" "by_name" {
  depends_on = [flexibleengine_sfs_turbo.test]

  name = flexibleengine_sfs_turbo.test[0].name
}

data "flexibleengine_sfs_turbos" "by_size" {
  depends_on = [flexibleengine_sfs_turbo.test]

  size = var.turbo_configuration[0]["size"]
}

data "flexibleengine_sfs_turbos" "by_share_type" {
  depends_on = [flexibleengine_sfs_turbo.test]

  share_type = var.turbo_configuration[1]["share_type"]
}

output "name_query_result_validation" {
  value = contains(data.flexibleengine_sfs_turbos.by_name.turbos[*].id,
  flexibleengine_sfs_turbo.test[0].id) && !contains(data.flexibleengine_sfs_turbos.by_name.turbos[*].id,
  flexibleengine_sfs_turbo.test[1].id) && !contains(data.flexibleengine_sfs_turbos.by_name.turbos[*].id,
  flexibleengine_sfs_turbo.test[2].id)
}

output "size_query_result_validation" {
  value = contains(data.flexibleengine_sfs_turbos.by_size.turbos[*].id,
  flexibleengine_sfs_turbo.test[0].id) && !contains(data.flexibleengine_sfs_turbos.by_size.turbos[*].id,
  flexibleengine_sfs_turbo.test[1].id) && !contains(data.flexibleengine_sfs_turbos.by_size.turbos[*].id,
  flexibleengine_sfs_turbo.test[2].id)
}

output "share_type_query_result_validation" {
  value = contains(data.flexibleengine_sfs_turbos.by_share_type.turbos[*].id,
  flexibleengine_sfs_turbo.test[1].id) && !contains(data.flexibleengine_sfs_turbos.by_share_type.turbos[*].id,
  flexibleengine_sfs_turbo.test[0].id) && !contains(data.flexibleengine_sfs_turbos.by_share_type.turbos[*].id,
  flexibleengine_sfs_turbo.test[2].id)
}
`, rName)
}
