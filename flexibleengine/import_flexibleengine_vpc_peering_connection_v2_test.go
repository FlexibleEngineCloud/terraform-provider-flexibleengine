package flexibleengine

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

// PASS
func TestAccFlexibleEngineVpcPeeringConnectionV1_importBasic(t *testing.T) {
	resourceName := "flexibleengine_vpc_peering_connection_v2.peering_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFlexibleEngineVpcPeeringConnectionV2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccFlexibleEngineVpcPeeringConnectionV2_basic,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
