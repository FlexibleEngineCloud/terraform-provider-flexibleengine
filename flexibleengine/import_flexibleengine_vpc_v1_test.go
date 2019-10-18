package flexibleengine

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccFlexibleEngineVpcV1_importBasic(t *testing.T) {
	resourceName := "flexibleengine_vpc_v1.vpc_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFlexibleEngineVpcV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcV1_basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
