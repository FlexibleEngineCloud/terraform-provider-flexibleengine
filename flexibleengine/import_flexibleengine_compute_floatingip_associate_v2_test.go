package flexibleengine

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccComputeV2FloatingIPAssociate_importBasic(t *testing.T) {
	resourceName := "flexibleengine_compute_floatingip_associate_v2.fip_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2FloatingIPAssociateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2FloatingIPAssociate_basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
