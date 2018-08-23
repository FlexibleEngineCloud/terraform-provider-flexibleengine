package flexibleengine

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccFlexibleEngineVpcRouteV2_importBasic(t *testing.T) {
	resourceName := "flexibleengine_vpc_route_v2.route_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFlexibleEngineRouteV2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRouteV2_basic,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
