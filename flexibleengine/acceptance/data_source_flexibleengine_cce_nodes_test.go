package acceptance

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNodesDataSource_basic(t *testing.T) {
	dataSourceName := "data.flexibleengine_cce_nodes.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNodesDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "nodes.0.name", rName),
				),
			},
		},
	})
}

func testAccNodesDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_cce_nodes" "test" {
  cluster_id = flexibleengine_cce_cluster_v3.test.id
  name       = flexibleengine_cce_node_v3.test.name

  depends_on = [flexibleengine_cce_node_v3.test]
}
`, testAccCceCluster_config(rName))
}
