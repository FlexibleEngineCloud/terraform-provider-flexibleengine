package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	rules "github.com/huaweicloud/golangsdk/openstack/waf/v1/webtamperprotection_rules"
)

func TestAccWafRuleWebTamperProtection_basic(t *testing.T) {
	var rule rules.WebTamper
	randName := acctest.RandString(5)
	resourceName := "flexibleengine_waf_rule_web_tamper_protection.rule_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWafWafRuleWebTamperProtectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafWafRuleWebTamperProtection_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafRuleWebTamperProtectionExists(resourceName, &rule),
					resource.TestCheckResourceAttr(resourceName, "domain", "www.abc.com"),
					resource.TestCheckResourceAttr(resourceName, "path", "/a"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccWafRuleImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccCheckWafWafRuleWebTamperProtectionDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	wafClient, err := config.WafV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_waf_rule_web_tamper_protection" {
			continue
		}

		policyID := rs.Primary.Attributes["policy_id"]
		_, err := rules.Get(wafClient, policyID, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("WAF rule still exists")
		}
	}

	return nil
}

func testAccCheckWafRuleWebTamperProtectionExists(n string, rule *rules.WebTamper) resource.TestCheckFunc {
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

		policyID := rs.Primary.Attributes["policy_id"]
		found, err := rules.Get(wafClient, policyID, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("WAF web tamper protection rule not found")
		}

		*rule = *found

		return nil
	}
}

func testAccWafWafRuleWebTamperProtection_basic(name string) string {
	return fmt.Sprintf(`
resource "flexibleengine_waf_policy" "policy_1" {
  name = "policy_%s"
}

resource "flexibleengine_waf_rule_web_tamper_protection" "rule_1" {
  policy_id = flexibleengine_waf_policy.policy_1.id
  domain    = "www.abc.com"
  path      = "/a"
}
`, name)
}
