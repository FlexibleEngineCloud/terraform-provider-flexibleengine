package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccRTSStackV1_importBasic(t *testing.T) {
	resourceName := "flexibleengine_rts_stack_v1.stack_1"
	var stackName = fmt.Sprintf("terra-test-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRTSStackV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRTSStackV1_basic(stackName),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
