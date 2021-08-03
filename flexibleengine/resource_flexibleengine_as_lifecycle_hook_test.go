package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/golangsdk/openstack/autoscaling/v1/lifecyclehooks"
)

func TestAccASLifecycleHook_basic(t *testing.T) {
	var hook lifecyclehooks.Hook
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceGroupName := "flexibleengine_as_group_v1.hth_as_group"
	resourceHookName := "flexibleengine_as_lifecycle_hook_v1.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckASLifecycleHookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testASLifecycleHook_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASLifecycleHookExists(resourceGroupName, resourceHookName, &hook),
					resource.TestCheckResourceAttr(resourceHookName, "name", rName),
					resource.TestCheckResourceAttr(resourceHookName, "type", "ADD"),
					resource.TestCheckResourceAttr(resourceHookName, "default_result", "ABANDON"),
					resource.TestCheckResourceAttr(resourceHookName, "timeout", "3600"),
					resource.TestCheckResourceAttr(resourceHookName, "notification_message", "This is a test message"),
					resource.TestCheckResourceAttrSet(resourceHookName, "notification_topic_urn"),
				),
			},
			{
				Config: testASLifecycleHook_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASLifecycleHookExists(resourceGroupName, resourceHookName, &hook),
					resource.TestCheckResourceAttr(resourceHookName, "name", rName),
					resource.TestCheckResourceAttr(resourceHookName, "type", "REMOVE"),
					resource.TestCheckResourceAttr(resourceHookName, "default_result", "CONTINUE"),
					resource.TestCheckResourceAttr(resourceHookName, "timeout", "600"),
					resource.TestCheckResourceAttr(resourceHookName, "notification_message",
						"This is a update message"),
				),
			},
			{
				ResourceName:      resourceHookName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccASLifecycleHookImportStateIdFunc(resourceGroupName, resourceHookName),
			},
		},
	})
}

func testAccCheckASLifecycleHookDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	asClient, err := config.AutoscalingV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating autoscaling client: %s", err)
	}

	var groupID string
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "flexibleengine_as_group_v1" {
			groupID = rs.Primary.ID
			break
		}
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_as_lifecycle_hook_v1" {
			continue
		}

		_, err := lifecyclehooks.Get(asClient, groupID, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("AS lifecycle hook still exists")
		}
	}

	return nil
}

func testAccCheckASLifecycleHookExists(resGroup, resHook string, hook *lifecyclehooks.Hook) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resGroup]
		if !ok {
			return fmt.Errorf("Not found: %s", resGroup)
		}
		groupID := rs.Primary.ID

		rs, ok = s.RootModule().Resources[resHook]
		if !ok {
			return fmt.Errorf("Not found: %s", resHook)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		asClient, err := config.AutoscalingV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating autoscaling client: %s", err)
		}
		found, err := lifecyclehooks.Get(asClient, groupID, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		hook = found

		return nil
	}
}

func testAccASLifecycleHookImportStateIdFunc(groupRes, hookRes string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		group, ok := s.RootModule().Resources[groupRes]
		if !ok {
			return "", fmt.Errorf("Auto Scaling group not found: %s", group)
		}
		hook, ok := s.RootModule().Resources[hookRes]
		if !ok {
			return "", fmt.Errorf("Auto Scaling lifecycle hook not found: %s", hook)
		}
		if group.Primary.ID == "" || hook.Primary.ID == "" {
			return "", fmt.Errorf("resource not found: %s/%s", group.Primary.ID, hook.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", group.Primary.ID, hook.Primary.ID), nil
	}
}

func testASLifecycleHook_base(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_smn_topic_v2" "test" {
  name = "%s"
}

resource "flexibleengine_smn_topic_v2" "update" {
  name = "%s-update"
}
`, testASV1Group_basic, rName, rName)
}

func testASLifecycleHook_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_as_lifecycle_hook_v1" "test" {
  name                   = "%s"
  type                   = "ADD"
  scaling_group_id       = flexibleengine_as_group_v1.hth_as_group.id
  notification_topic_urn = flexibleengine_smn_topic_v2.test.topic_urn
  notification_message   = "This is a test message"
}
`, testASLifecycleHook_base(rName), rName)
}

func testASLifecycleHook_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_as_lifecycle_hook_v1" "test" {
  name                   = "%s"
  type                   = "REMOVE"
  scaling_group_id       = flexibleengine_as_group_v1.hth_as_group.id
  notification_topic_urn = flexibleengine_smn_topic_v2.update.topic_urn
  notification_message   = "This is a update message"
  default_result         = "CONTINUE"
  timeout                = 600
}
`, testASLifecycleHook_base(rName), rName)
}
