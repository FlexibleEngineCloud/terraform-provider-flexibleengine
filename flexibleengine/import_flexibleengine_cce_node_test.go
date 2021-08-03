package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCCENodeV3_importBasic(t *testing.T) {
	resourceName := "flexibleengine_cce_node_v3.node_1"
	var cceName = fmt.Sprintf("terra-test-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCENodeV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodeV3_basic(cceName),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCCENodeImportStateIdFunc(),
			},
		},
	})
}

func testAccCCENodeImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		cluster, ok := s.RootModule().Resources["flexibleengine_cce_cluster_v3.cluster_1"]
		if !ok {
			return "", fmt.Errorf("Cluster not found: %s", cluster)
		}
		node, ok := s.RootModule().Resources["flexibleengine_cce_node_v3.node_1"]
		if !ok {
			return "", fmt.Errorf("Node not found: %s", node)
		}

		if cluster.Primary.ID == "" || node.Primary.ID == "" {
			return "", fmt.Errorf("resource not found: %s/%s", cluster.Primary.ID, node.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", cluster.Primary.ID, node.Primary.ID), nil
	}
}
