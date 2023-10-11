package acceptance

import (
	"fmt"
	"testing"
	"time"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	aomservice "github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/aom"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"

	aom "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/aom/v2/model"
)

func getServiceDiscoveryRuleResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.HcAomV2Client(OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating AOM client: %s", err)
	}

	// wait 30 seconds before listing the rule, to avoid error
	// lintignore:R018
	time.Sleep(30 * time.Second)

	response, err := c.ListServiceDiscoveryRules(&aom.ListServiceDiscoveryRulesRequest{})
	if err != nil {
		return nil, fmt.Errorf("error retrieving AOM service discovery rule: %s", state.Primary.ID)
	}

	allRules := *response.AppRules

	return aomservice.FilterRules(allRules, state.Primary.ID)
}

func TestAccAOMServiceDiscoveryRule_basic(t *testing.T) {
	var ar aom.QueryAlarmResult
	rName := acceptance.RandomAccResourceNameWithDash()
	rNameUpdate := rName + "-update"
	resourceName := "flexibleengine_aom_service_discovery_rule.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&ar,
		getServiceDiscoveryRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAOMServiceDiscoveryRule_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "priority", "9999"),
					resource.TestCheckResourceAttr(resourceName, "detect_log_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "discovery_rule_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "is_default_rule", "true"),
					resource.TestCheckResourceAttr(resourceName, "log_file_suffix.0", "log"),
					resource.TestCheckResourceAttr(resourceName, "service_type", "Python"),
					resource.TestCheckResourceAttr(resourceName, "discovery_rules.0.check_content.0", "python"),
					resource.TestCheckResourceAttr(resourceName, "log_path_rules.0.args.0", "python"),
					resource.TestCheckResourceAttr(
						resourceName, "name_rules.0.service_name_rule.0.args.0", "python"),
					resource.TestCheckResourceAttr(
						resourceName, "name_rules.0.application_name_rule.0.args.0", "python"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAOMServiceDiscoveryRule_update(rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "priority", "9998"),
					resource.TestCheckResourceAttr(resourceName, "detect_log_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "discovery_rule_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_default_rule", "false"),
					resource.TestCheckResourceAttr(resourceName, "log_file_suffix.0", "trace"),
					resource.TestCheckResourceAttr(resourceName, "service_type", "Java"),
					resource.TestCheckResourceAttr(resourceName, "discovery_rules.0.check_content.0", "java"),
					resource.TestCheckResourceAttr(resourceName, "log_path_rules.0.args.0", "java"),
					resource.TestCheckResourceAttr(
						resourceName, "name_rules.0.service_name_rule.0.args.0", "java"),
					resource.TestCheckResourceAttr(
						resourceName, "name_rules.0.application_name_rule.0.args.0", "java"),
				),
			},
			{
				Config: testAOMServiceDiscoveryRule_update2(rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "priority", "9998"),
					resource.TestCheckResourceAttr(resourceName, "detect_log_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "discovery_rule_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_default_rule", "false"),
					resource.TestCheckResourceAttr(resourceName, "log_file_suffix.0", "out"),
					resource.TestCheckResourceAttr(resourceName, "service_type", "Java"),
					resource.TestCheckResourceAttr(resourceName, "discovery_rules.0.check_content.0", "java"),
					resource.TestCheckResourceAttr(resourceName, "log_path_rules.0.args.0", "java"),
					resource.TestCheckResourceAttr(
						resourceName, "name_rules.0.service_name_rule.0.args.0", "java"),
					resource.TestCheckResourceAttr(
						resourceName, "name_rules.0.application_name_rule.0.args.0", "java"),
				),
			},
		},
	})
}

func testAOMServiceDiscoveryRule_basic(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_aom_service_discovery_rule" "test" {
  name                   = "%s"
  priority               = 9999
  detect_log_enabled     = true
  discovery_rule_enabled = true
  is_default_rule        = true
  log_file_suffix        = ["log"]
  service_type           = "Python"

  discovery_rules {
    check_content = ["python"]
    check_mode    = "contain"
    check_type    = "cmdLine"
  }

  log_path_rules {
    name_type = "cmdLineHash"
    args      = ["python"]
    value     = ["/tmp/log"]
  }

  name_rules {
    service_name_rule {
      name_type = "str"
      args      = ["python"]
    }
    application_name_rule {
      name_type = "str"
      args      = ["python"]
    }
  }
}
`, rName)
}

func testAOMServiceDiscoveryRule_update(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_aom_service_discovery_rule" "test" {
  name                   = "%s"
  priority               = 9998
  detect_log_enabled     = false
  discovery_rule_enabled = false
  is_default_rule        = false
  log_file_suffix        = ["trace"]
  service_type           = "Java"

  discovery_rules {
    check_content = ["java"]
    check_mode    = "contain"
    check_type    = "cmdLine"
  }

  log_path_rules {
    name_type = "cmdLineHash"
    args      = ["java"]
    value     = ["/tmp/log"]
  }

  name_rules {
    service_name_rule {
      name_type = "str"
      args      = ["java"]
    }
    application_name_rule {
      name_type = "str"
      args      = ["java"]
    }
  }
}
`, rName)
}

func testAOMServiceDiscoveryRule_update2(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_aom_service_discovery_rule" "test" {
  name                   = "%s"
  priority               = 9998
  detect_log_enabled     = false
  discovery_rule_enabled = false
  is_default_rule        = false
  log_file_suffix        = ["out"]
  service_type           = "Java"

  discovery_rules {
    check_content = ["java"]
    check_mode    = "contain"
    check_type    = "cmdLine"
  }

  log_path_rules {
    name_type = "cmdLineHash"
    args      = ["java"]
    value     = ["/tmp/log"]
  }

  name_rules {
    service_name_rule {
      name_type = "str"
      args      = ["java"]
    }
    application_name_rule {
      name_type = "str"
      args      = ["java"]
    }
  }
}
`, rName)
}
