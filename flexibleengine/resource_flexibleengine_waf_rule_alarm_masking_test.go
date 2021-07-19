package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	rules "github.com/huaweicloud/golangsdk/openstack/waf/v1/falsealarmmasking_rules"
)

func TestAccWafRuleAlarmMasking_basic(t *testing.T) {
	var rule rules.AlarmMasking
	randName := acctest.RandString(5)
	resourceName := "flexibleengine_waf_rule_alarm_masking.rule_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWafRuleAlarmMaskingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafRuleAlarmMasking_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafRuleAlarmMaskingExists(resourceName, &rule),
					resource.TestCheckResourceAttr(resourceName, "path", "/a"),
					resource.TestCheckResourceAttrSet(resourceName, "event_type"),
				),
			},
			{
				Config: testAccWafRuleAlarmMasking_update(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafRuleAlarmMaskingExists(resourceName, &rule),
					resource.TestCheckResourceAttr(resourceName, "path", "/abc"),
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

func testAccCheckWafRuleAlarmMaskingDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	wafClient, err := config.WafV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_waf_rule_alarm_masking" {
			continue
		}

		policyID := rs.Primary.Attributes["policy_id"]
		_, err := rules.Get(wafClient, policyID, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("WAF alarm masking rule still exists")
		}
	}

	return nil
}

func testAccCheckWafRuleAlarmMaskingExists(n string, rule *rules.AlarmMasking) resource.TestCheckFunc {
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
			return fmt.Errorf("WAF alarm masking rule not found")
		}

		*rule = *found

		return nil
	}
}

func testAccWafRuleAlarmMasking_basic(name string) string {
	return fmt.Sprintf(`
resource "flexibleengine_waf_policy" "policy_1" {
  name = "policy_%s"
}

resource "flexibleengine_waf_rule_alarm_masking" "rule_1" {
  policy_id = flexibleengine_waf_policy.policy_1.id
  path      = "/a"
  event_id  = "3737fb122f2140f39292f597ad3b7e9a"
}
`, name)
}

func testAccWafRuleAlarmMasking_update(name string) string {
	return fmt.Sprintf(`
resource "flexibleengine_waf_policy" "policy_1" {
  name = "policy_%s"
}

resource "flexibleengine_waf_rule_alarm_masking" "rule_1" {
  policy_id = flexibleengine_waf_policy.policy_1.id
  path      = "/abc"
  event_id  = "3737fb122f2140f39292f597ad3b7e9a"
}
`, name)
}
