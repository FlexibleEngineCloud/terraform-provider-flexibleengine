package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	rules "github.com/huaweicloud/golangsdk/openstack/waf/v1/whiteblackip_rules"
)

func TestAccWafRuleBlackList_basic(t *testing.T) {
	var rule rules.WhiteBlackIP
	randName := acctest.RandString(5)
	resourceName := "flexibleengine_waf_rule_blacklist.rule_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWafRuleBlackListDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafRuleBlackList_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafRuleBlackListExists(resourceName, &rule),
					resource.TestCheckResourceAttr(resourceName, "address", "192.168.0.0/24"),
					resource.TestCheckResourceAttr(resourceName, "action", "0"),
				),
			},
			{
				Config: testAccWafRuleBlackList_update(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafRuleBlackListExists(resourceName, &rule),
					resource.TestCheckResourceAttr(resourceName, "address", "192.168.0.125"),
					resource.TestCheckResourceAttr(resourceName, "action", "1"),
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

func testAccCheckWafRuleBlackListDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	wafClient, err := config.WafV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_waf_rule_blacklist" {
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

func testAccCheckWafRuleBlackListExists(n string, rule *rules.WhiteBlackIP) resource.TestCheckFunc {
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
			return fmt.Errorf("WAF black list rule not found")
		}

		*rule = *found

		return nil
	}
}

func testAccWafRuleBlackList_basic(name string) string {
	return fmt.Sprintf(`
resource "flexibleengine_waf_policy" "policy_1" {
  name = "policy_%s"
}

resource "flexibleengine_waf_rule_blacklist" "rule_1" {
  policy_id = flexibleengine_waf_policy.policy_1.id
  address   = "192.168.0.0/24"
}
`, name)
}

func testAccWafRuleBlackList_update(name string) string {
	return fmt.Sprintf(`
resource "flexibleengine_waf_policy" "policy_1" {
  name = "policy_%s"
}

resource "flexibleengine_waf_rule_blacklist" "rule_1" {
  policy_id = flexibleengine_waf_policy.policy_1.id
  address   = "192.168.0.125"
  action    = 1
}
`, name)
}
