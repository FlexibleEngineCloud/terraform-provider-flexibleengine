package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCCEAddonTemplateDataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCEClusterV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEAddonTemplateDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.flexibleengine_cce_addon_template.autoscaler", "spec"),
					resource.TestCheckResourceAttrSet("data.flexibleengine_cce_addon_template.metrics-server", "spec"),
				),
			},
		},
	})
}

func testAccCCEAddonTemplateDataSource_basic(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_cce_cluster_v3" "test" {
  name         = "%s"
  cluster_type = "VirtualMachine"
  flavor_id    = "cce.s1.small"
  vpc_id       = "%s"
  subnet_id    = "%s"
  container_network_type = "overlay_l2"
}

data "flexibleengine_cce_addon_template" "autoscaler" {
  cluster_id = flexibleengine_cce_cluster_v3.test.id
  name       = "autoscaler"
  version    = "1.19.1"
}

data "flexibleengine_cce_addon_template" "metrics-server" {
  cluster_id = flexibleengine_cce_cluster_v3.test.id
  name       = "metrics-server"
  version    = "1.0.6"
}
`, rName, OS_VPC_ID, OS_NETWORK_ID)
}
