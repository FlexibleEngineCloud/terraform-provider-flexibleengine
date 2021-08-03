package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/golangsdk/openstack/cloudeyeservice/alarmrule"
)

func TestCESAlarmRule_basic(t *testing.T) {
	var ar alarmrule.AlarmRule
	resourceName := "flexibleengine_ces_alarmrule.alarmrule_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCESAlarmRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCESAlarmRule_basic,
				Check: resource.ComposeTestCheckFunc(
					testCESAlarmRuleExists(resourceName, &ar),
					resource.TestCheckResourceAttr(resourceName, "alarm_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "alarm_action_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "alarm_level", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testCESAlarmRule_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "alarm_enabled", "false"),
				),
			},
		},
	})
}

func testCESAlarmRuleDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	cesClient, err := config.CesV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine ces client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_ces_alarmrule" {
			continue
		}

		id := rs.Primary.ID
		_, err := alarmrule.Get(cesClient, id).Extract()
		if err == nil {
			return fmt.Errorf("Alarm rule still exists")
		}
	}

	return nil
}

func testCESAlarmRuleExists(n string, ar *alarmrule.AlarmRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		cesClient, err := config.CesV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine ces client: %s", err)
		}

		id := rs.Primary.ID
		found, err := alarmrule.Get(cesClient, id).Extract()
		if err != nil {
			return err
		}

		*ar = *found

		return nil
	}
}

var testCESAlarmRule_basic = fmt.Sprintf(`
resource "flexibleengine_compute_instance_v2" "vm_1" {
  name = "instance_1"
  network {
    uuid = "%s"
  }
}

resource "flexibleengine_smn_topic_v2" "topic_1" {
  name         = "topic_1"
  display_name = "The display name of topic_1"
}

resource "flexibleengine_ces_alarmrule" "alarmrule_1" {
  alarm_name           = "alarm_rule1"
  alarm_action_enabled = true

  metric {
    namespace   = "SYS.ECS"
    metric_name = "network_outgoing_bytes_rate_inband"
    dimensions {
        name  = "instance_id"
        value = flexibleengine_compute_instance_v2.vm_1.id
    }
  }
  condition  {
    period              = 300
    filter              = "average"
    comparison_operator = ">"
    value               = 6
    unit                = "B/s"
    count               = 1
  } 

  alarm_actions {
    type = "notification"
    notification_list = [
      flexibleengine_smn_topic_v2.topic_1.topic_urn
    ]
  }
}
`, OS_NETWORK_ID)

var testCESAlarmRule_update = fmt.Sprintf(`
resource "flexibleengine_compute_instance_v2" "vm_1" {
  name = "instance_1"
  network {
    uuid = "%s"
  }
}

resource "flexibleengine_smn_topic_v2" "topic_1" {
  name         = "topic_1"
  display_name = "The display name of topic_1"
}

resource "flexibleengine_ces_alarmrule" "alarmrule_1" {
  alarm_name           = "alarm_rule1"
  alarm_enabled        = false
  alarm_action_enabled = true

  metric {
    namespace   = "SYS.ECS"
    metric_name = "network_outgoing_bytes_rate_inband"
    dimensions {
        name  = "instance_id"
        value = flexibleengine_compute_instance_v2.vm_1.id
    }
  }
  condition  {
    period              = 300
    filter              = "average"
    comparison_operator = ">"
    value               = 6
    unit                = "B/s"
    count               = 1
  }

  alarm_actions {
    type = "notification"
    notification_list = [
      flexibleengine_smn_topic_v2.topic_1.topic_urn
    ]
  }
}
`, OS_NETWORK_ID)
