package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	rules "github.com/huaweicloud/golangsdk/openstack/waf/v1/ccattackprotection_rules"
)

func TestAccWafRuleCCAttackProtection_basic(t *testing.T) {
	var rule rules.CcAttack
	randName := acctest.RandString(5)
	resourceName := "flexibleengine_waf_rule_cc_protection.rule_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWafRuleCCAttackProtectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafRuleCCAttackProtection_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafRuleCCAttackProtectionExists(resourceName, &rule),
					resource.TestCheckResourceAttr(resourceName, "path", "/abc"),
					resource.TestCheckResourceAttr(resourceName, "limit_num", "10"),
					resource.TestCheckResourceAttr(resourceName, "limit_period", "60"),
					resource.TestCheckResourceAttr(resourceName, "block_time", "10"),
					resource.TestCheckResourceAttr(resourceName, "mode", "cookie"),
					resource.TestCheckResourceAttr(resourceName, "action", "block"),
				),
			},
			{
				Config: testAccWafRuleCCAttackProtection_update(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafRuleCCAttackProtectionExists(resourceName, &rule),
					resource.TestCheckResourceAttr(resourceName, "path", "/abcd"),
					resource.TestCheckResourceAttr(resourceName, "limit_num", "30"),
					resource.TestCheckResourceAttr(resourceName, "limit_period", "100"),
					resource.TestCheckResourceAttr(resourceName, "block_time", "20"),
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

func testAccCheckWafRuleCCAttackProtectionDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	wafClient, err := config.WafV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_waf_rule_cc_protection" {
			continue
		}

		policyID := rs.Primary.Attributes["policy_id"]
		_, err := rules.Get(wafClient, policyID, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Waf rule still exists")
		}
	}

	return nil
}

func testAccCheckWafRuleCCAttackProtectionExists(n string, rule *rules.CcAttack) resource.TestCheckFunc {
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
			return fmt.Errorf("WAF rule not found")
		}

		*rule = *found

		return nil
	}
}

func testAccWafRuleCCAttackProtection_basic(name string) string {
	return fmt.Sprintf(`
resource "flexibleengine_waf_policy" "policy_1" {
  name = "policy_%s"
}

resource "flexibleengine_waf_rule_cc_protection" "rule_1" {
  policy_id    = flexibleengine_waf_policy.policy_1.id
  path         = "/abc"
  limit_num    = 10
  limit_period = 60
  mode         = "cookie"
  cookie       = "sessionid"

  action             = "block"
  block_time         = 10
  block_page_type = "application/json"
  block_page_content      = "{\"error\":\"forbidden\"}"
}
`, name)
}

func testAccWafRuleCCAttackProtection_update(name string) string {
	return fmt.Sprintf(`
resource "flexibleengine_waf_policy" "policy_1" {
  name = "policy_%s"
}

resource "flexibleengine_waf_rule_cc_protection" "rule_1" {
  policy_id    = flexibleengine_waf_policy.policy_1.id
  path         = "/abcd"
  limit_num    = 30
  limit_period = 100
  mode         = "cookie"
  cookie       = "sessionid"

  action             = "block"
  block_time         = 20
  block_page_type    = "application/json"
  block_page_content = "{\"error\":\"forbidden\"}"
}
`, name)
}
