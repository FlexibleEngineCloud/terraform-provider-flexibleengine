package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/golangsdk/openstack/waf/v1/policies"
)

func TestAccWafPolicyV1_basic(t *testing.T) {
	var policy policies.Policy
	randName := acctest.RandString(5)
	resourceName := "flexibleengine_waf_policy.policy_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWafPolicyV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafPolicyV1_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafPolicyV1Exists(resourceName, &policy),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("policy-%s", randName)),
					resource.TestCheckResourceAttr(resourceName, "protection_mode", "log"),
					resource.TestCheckResourceAttr(resourceName, "level", "2"),
					resource.TestCheckResourceAttr(resourceName, "full_detection", "false"),
					resource.TestCheckResourceAttr(resourceName, "protection_status.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "protection_status.0.basic_web_protection", "true"),
				),
			},
			{
				Config: testAccWafPolicyV1_update(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafPolicyV1Exists(resourceName, &policy),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("policy-%s-updated", randName)),
					resource.TestCheckResourceAttr(resourceName, "protection_mode", "block"),
					resource.TestCheckResourceAttr(resourceName, "level", "1"),
					resource.TestCheckResourceAttr(resourceName, "full_detection", "true"),
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

func testAccCheckWafPolicyV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	wafClient, err := config.WafV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_waf_policy" {
			continue
		}

		_, err := policies.Get(wafClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Waf policy still exists")
		}
	}

	return nil
}

func testAccCheckWafPolicyV1Exists(n string, policy *policies.Policy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		wafClient, err := config.WafV1Client(OS_REGION_NAME)
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

func testAccWafPolicyV1_basic(name string) string {
	return fmt.Sprintf(`
resource "flexibleengine_waf_policy" "policy_1" {
  name = "policy-%s"
}
`, name)
}

func testAccWafPolicyV1_update(name string) string {
	return fmt.Sprintf(`
resource "flexibleengine_waf_policy" "policy_1" {
  name            = "policy-%s-updated"
  level           = 1
  protection_mode = "block"
  full_detection  = true
}
`, name)
}
