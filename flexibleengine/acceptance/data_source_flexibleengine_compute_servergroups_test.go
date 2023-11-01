package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccComputeServerGroupsDataSource_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()
	dataSourceName := "data.flexibleengine_compute_servergroups.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeServerGroupsDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", rName),
					resource.TestCheckResourceAttr(dataSourceName, "servergroups.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "servergroups.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "servergroups.0.name"),
				),
			},
		},
	})
}

func testAccComputeServerGroupsDataSource_basic(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_compute_servergroup_v2" "test" {
  name     = "%s"
  policies = ["anti-affinity"]
}

data "flexibleengine_compute_servergroups" "test" {
  name = flexibleengine_compute_servergroup_v2.test.name
}
`, rName)
}
