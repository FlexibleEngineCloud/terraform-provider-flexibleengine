package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	rules "github.com/huaweicloud/golangsdk/openstack/waf/v1/preciseprotection_rules"
)

func TestAccWafRulePreciseProtection_basic(t *testing.T) {
	var rule rules.Precise
	randName := acctest.RandString(5)
	resourceName := "flexibleengine_waf_rule_precise_protection.rule_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWafRulePreciseProtectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafRulePreciseProtection_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafRulePreciseProtectionExists(resourceName, &rule),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("rule_%s", randName)),
					resource.TestCheckResourceAttr(resourceName, "action", "block"),
					resource.TestCheckResourceAttr(resourceName, "priority", "10"),
					resource.TestCheckResourceAttr(resourceName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "conditions.0.content", "/login"),
					resource.TestCheckNoResourceAttr(resourceName, "start_time"),
					resource.TestCheckNoResourceAttr(resourceName, "end_time"),
				),
			},
			{
				Config: testAccWafRulePreciseProtection_update(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafRulePreciseProtectionExists(resourceName, &rule),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("rule_%s_update", randName)),
					resource.TestCheckResourceAttr(resourceName, "action", "block"),
					resource.TestCheckResourceAttr(resourceName, "priority", "20"),
					resource.TestCheckResourceAttr(resourceName, "conditions.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "conditions.0.content", "/login"),
					resource.TestCheckResourceAttr(resourceName, "conditions.1.content", "192.168.1.1"),
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

func TestAccWafRulePreciseProtection_time(t *testing.T) {
	var rule rules.Precise
	randName := acctest.RandString(5)
	resourceName := "flexibleengine_waf_rule_precise_protection.rule_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWafRulePreciseProtectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafRulePreciseProtection_time(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafRulePreciseProtectionExists(resourceName, &rule),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("rule_%s", randName)),
					resource.TestCheckResourceAttr(resourceName, "action", "block"),
					resource.TestCheckResourceAttr(resourceName, "priority", "10"),
					resource.TestCheckResourceAttr(resourceName, "start_time", "2021-10-01 00:00:00"),
					resource.TestCheckResourceAttr(resourceName, "end_time", "2021-12-31 23:59:59"),
					resource.TestCheckResourceAttr(resourceName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "conditions.0.content", "/login"),
				),
			},
		},
	})
}

func testAccCheckWafRulePreciseProtectionDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	wafClient, err := config.WafV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_waf_rule_precise_protection" {
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

func testAccCheckWafRulePreciseProtectionExists(n string, rule *rules.Precise) resource.TestCheckFunc {
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

func testAccWafRulePreciseProtection_basic(name string) string {
	return fmt.Sprintf(`
resource "flexibleengine_waf_policy" "policy_1" {
  name = "policy_%s"
}

resource "flexibleengine_waf_rule_precise_protection" "rule_1" {
  policy_id = flexibleengine_waf_policy.policy_1.id
  name      = "rule_%s"
  priority  = 10
  
  conditions {
    field   = "path"
    logic   = "contain"
    content = "/login"
  }
}
`, name, name)
}

func testAccWafRulePreciseProtection_update(name string) string {
	return fmt.Sprintf(`
resource "flexibleengine_waf_policy" "policy_1" {
  name = "policy_%s"
}

resource "flexibleengine_waf_rule_precise_protection" "rule_1" {
  policy_id = flexibleengine_waf_policy.policy_1.id
  name      = "rule_%s_update"
  priority  = 20
  
  conditions {
    field   = "path"
    logic   = "contain"
    content = "/login"
  }
  conditions {
    field   = "ip"
    logic   = "equal"
    content = "192.168.1.1"
  }
}
`, name, name)
}

func testAccWafRulePreciseProtection_time(name string) string {
	return fmt.Sprintf(`
resource "flexibleengine_waf_policy" "policy_1" {
  name = "policy_%s"
}

resource "flexibleengine_waf_rule_precise_protection" "rule_1" {
  policy_id  = flexibleengine_waf_policy.policy_1.id
  name       = "rule_%s"
  action     = "block"
  priority   = 10
  start_time = "2021-10-01 00:00:00"
  end_time   = "2021-12-31 23:59:59"
  
  conditions {
    field   = "path"
    logic   = "prefix"
    content = "/login"
  }
}
`, name, name)
}
