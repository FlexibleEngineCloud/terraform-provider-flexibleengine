package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/waf_hw/v1/policies"
)

func TestAccWafDedicatedPolicyV1_basic(t *testing.T) {
	var policy policies.Policy
	randName := acceptance.RandomAccResourceName()
	resourceName := "flexibleengine_waf_dedicated_policy.policy_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckWafDedicatedPolicyV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafDedicatedPolicyV1_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDedicatedPolicyV1Exists(resourceName, &policy),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "level", "1"),
					resource.TestCheckResourceAttr(resourceName, "full_detection", "false"),
				),
			},
			{
				Config: testAccWafDedicatedPolicyV1_update(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDedicatedPolicyV1Exists(resourceName, &policy),
					resource.TestCheckResourceAttr(resourceName, "name", randName+"_updated"),
					resource.TestCheckResourceAttr(resourceName, "protection_mode", "block"),
					resource.TestCheckResourceAttr(resourceName, "level", "3"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckWafDedicatedPolicyV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	wafClient, err := wafDedicatedv1Client(config, OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_waf_dedicated_policy" {
			continue
		}
		_, err := policies.Get(wafClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Waf policy still exists")
		}
	}
	return nil
}

func testAccCheckWafDedicatedPolicyV1Exists(n string, policy *policies.Policy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		wafClient, err := wafDedicatedv1Client(config, OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating Flexibleengine WAF client: %s", err)
		}

		found, err := policies.Get(wafClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("Waf policy not found")
		}

		*policy = *found
		return nil
	}
}

func testAccWafDedicatedPolicyV1_basic(name string) string {
	return fmt.Sprintf(`

resource "flexibleengine_waf_dedicated_policy" "policy_1" {
  name  = "%s"
  level = 1
}
`, name)
}

func testAccWafDedicatedPolicyV1_update(name string) string {
	return fmt.Sprintf(`

resource "flexibleengine_waf_dedicated_policy" "policy_1" {
  name            = "%s_updated"
  protection_mode = "block"
  level           = 3
}
`, name)
}
