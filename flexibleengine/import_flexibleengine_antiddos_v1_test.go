package flexibleengine

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAntiDdosV1_importBasic(t *testing.T) {
	resourceName := "flexibleengine_antiddos_v1.antiddos_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAntiDdosV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAntiDdosV1_basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
