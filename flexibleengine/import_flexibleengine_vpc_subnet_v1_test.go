package flexibleengine

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

// PASS
func TestAccFlexibleEngineVpcSubnetV1_importBasic(t *testing.T) {
	resourceName := "flexibleengine_vpc_subnet_v1.subnet_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFlexibleEngineVpcSubnetV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccFlexibleEngineVpcSubnetV1_basic,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
