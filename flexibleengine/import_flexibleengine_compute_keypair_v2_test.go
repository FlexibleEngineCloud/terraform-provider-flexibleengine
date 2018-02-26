package flexibleengine

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

// PASS
func TestAccComputeV2Keypair_importBasic(t *testing.T) {
	resourceName := "flexibleengine_compute_keypair_v2.orange_acctest_kp"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2KeypairDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeV2Keypair_basic,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
