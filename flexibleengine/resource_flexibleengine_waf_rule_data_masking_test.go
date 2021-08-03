package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	rules "github.com/huaweicloud/golangsdk/openstack/waf/v1/datamasking_rules"
)

func TestAccWafRuleDataMasking_basic(t *testing.T) {
	var rule rules.DataMasking
	randName := acctest.RandString(5)
	resourceName := "flexibleengine_waf_rule_data_masking.rule_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWafRuleDataMaskingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafRuleDataMasking_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafRuleDataMaskingExists(resourceName, &rule),
					resource.TestCheckResourceAttr(resourceName, "path", "/login"),
					resource.TestCheckResourceAttr(resourceName, "field", "params"),
					resource.TestCheckResourceAttr(resourceName, "subfield", "password"),
				),
			},
			{
				Config: testAccWafRuleDataMasking_update(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafRuleDataMaskingExists(resourceName, &rule),
					resource.TestCheckResourceAttr(resourceName, "path", "/login_new"),
					resource.TestCheckResourceAttr(resourceName, "field", "params"),
					resource.TestCheckResourceAttr(resourceName, "subfield", "secret"),
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

func testAccCheckWafRuleDataMaskingDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	wafClient, err := config.WafV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_waf_rule_data_masking" {
			continue
		}

		policyID := rs.Primary.Attributes["policy_id"]
		_, err := rules.Get(wafClient, policyID, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("WAF data masking rule still exists")
		}
	}

	return nil
}

func testAccCheckWafRuleDataMaskingExists(n string, rule *rules.DataMasking) resource.TestCheckFunc {
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
			return fmt.Errorf("WAF data masking rule not found")
		}

		*rule = *found

		return nil
	}
}

func testAccWafRuleImportStateIdFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		policy, ok := s.RootModule().Resources["flexibleengine_waf_policy.policy_1"]
		if !ok {
			return "", fmt.Errorf("WAF policy not found")
		}
		rule, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("WAF rule not found")
		}

		if policy.Primary.ID == "" || rule.Primary.ID == "" {
			return "", fmt.Errorf("resource not found: %s/%s", policy.Primary.ID, rule.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", policy.Primary.ID, rule.Primary.ID), nil
	}
}

func testAccWafRuleDataMasking_basic(name string) string {
	return fmt.Sprintf(`
resource "flexibleengine_waf_policy" "policy_1" {
  name = "policy_%s"
}

resource "flexibleengine_waf_rule_data_masking" "rule_1" {
  policy_id = flexibleengine_waf_policy.policy_1.id
  path      = "/login"
  field     = "params"
  subfield  = "password"
}
`, name)
}

func testAccWafRuleDataMasking_update(name string) string {
	return fmt.Sprintf(`
resource "flexibleengine_waf_policy" "policy_1" {
  name = "policy_%s"
}

resource "flexibleengine_waf_rule_data_masking" "rule_1" {
  policy_id = flexibleengine_waf_policy.policy_1.id
  path      = "/login_new"
  field     = "params"
  subfield  = "secret"
}
`, name)
}
