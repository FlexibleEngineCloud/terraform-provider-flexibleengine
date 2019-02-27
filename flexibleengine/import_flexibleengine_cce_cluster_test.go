package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccCCEClusterV3_importBasic(t *testing.T) {
	resourceName := "flexibleengine_cce_cluster_v3.cluster_1"
	var cceName = fmt.Sprintf("terra_test_%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCEClusterV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEClusterV3_basic(cceName),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
