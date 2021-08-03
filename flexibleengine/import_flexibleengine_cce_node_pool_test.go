package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCCENodePool_importBasic(t *testing.T) {
	resourceName := "flexibleengine_cce_node_pool_v3.test"
	var cceName = fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCENodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodePool_basic(cceName),
			},

			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateIdFunc:       nodePoolImportStateIdFunc(resourceName),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"max_pods"},
			},
		},
	})
}

func nodePoolImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		cluster_id := rs.Primary.Attributes["cluster_id"]
		return fmt.Sprintf("%s/%s", cluster_id, rs.Primary.Attributes["id"]), nil
	}
}
