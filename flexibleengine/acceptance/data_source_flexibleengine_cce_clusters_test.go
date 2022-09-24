package acceptance

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCCEClustersDataSource_basic(t *testing.T) {
	dataSourceName := "data.flexibleengine_cce_clusters.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEClusterV3DataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "clusters.0.name", rName),
					resource.TestCheckResourceAttr(dataSourceName, "clusters.0.status", "Available"),
					resource.TestCheckResourceAttr(dataSourceName, "clusters.0.cluster_type", "VirtualMachine"),
				),
			},
		},
	})
}

func testAccCCEClusterV3DataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_cce_clusters" "test" {
  name = flexibleengine_cce_cluster_v3.test.name

  depends_on = [flexibleengine_cce_cluster_v3.test]
}
`, testAccCceCluster_config(rName))
}
