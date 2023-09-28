package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/smn/v2/topics"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getResourceSMNTopic(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	smnClient, err := conf.SmnV2Client(OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SMN client: %s", err)
	}

	return topics.Get(smnClient, state.Primary.ID).Extract()
}

func TestAccSMNV2Topic_basic(t *testing.T) {
	var topic topics.TopicGet
	resourceName := "flexibleengine_smn_topic_v2.topic_1"
	rName := acceptance.RandomAccResourceNameWithDash()
	displayName := fmt.Sprintf("The display name of %s", rName)
	update_displayName := fmt.Sprintf("The update display name of %s", rName)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&topic,
		getResourceSMNTopic,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSMNV2TopicConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "display_name", displayName),
				),
			},
			{
				Config: testAccSMNV2TopicConfig_update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "display_name", update_displayName),
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

func TestAccSMNV2Topic_policies(t *testing.T) {
	var topic topics.TopicGet
	resourceName := "flexibleengine_smn_topic_v2.topic_1"
	rName := acceptance.RandomAccResourceNameWithDash()
	displayName := fmt.Sprintf("The display name of %s", rName)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&topic,
		getResourceSMNTopic,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSMNV2TopicConfig_policies_step1(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "display_name", displayName),
					resource.TestCheckResourceAttr(resourceName, "users_publish_allowed", "*"),
					resource.TestCheckResourceAttr(resourceName, "services_publish_allowed", "obs,vod"),
					resource.TestCheckResourceAttr(resourceName, "introduction", "created by terraform"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSMNV2TopicConfig_policies_step2(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "display_name", displayName),
					resource.TestCheckResourceAttr(resourceName, "users_publish_allowed",
						"urn:csp:iam::0970d7b7d400f2470fbec00316a03560:root,urn:csp:iam::0970d7b7d400f2470fbec00316a03561:root"),
					resource.TestCheckResourceAttr(resourceName, "services_publish_allowed", "obs,vod,cce"),
					resource.TestCheckResourceAttr(resourceName, "introduction", "created by terraform update"),
				),
			},
			{
				Config: testAccSMNV2TopicConfig_policies_step3(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "display_name", displayName),
					resource.TestCheckResourceAttr(resourceName, "users_publish_allowed", ""),
					resource.TestCheckResourceAttr(resourceName, "services_publish_allowed", ""),
					resource.TestCheckResourceAttr(resourceName, "introduction", ""),
				),
			},
		},
	})
}

func testAccSMNV2TopicConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_smn_topic_v2" "topic_1" {
  name                  = "%s"
  display_name          = "The display name of %s"
  enterprise_project_id = "0"
}
`, rName, rName)
}

func testAccSMNV2TopicConfig_update(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_smn_topic_v2" "topic_1" {
  name         = "%s"
  display_name = "The update display name of %s"
}
`, rName, rName)
}

func testAccSMNV2TopicConfig_policies_step1(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_smn_topic_v2" "topic_1" {
  name                     = "%s"
  display_name             = "The display name of %s"
  users_publish_allowed    = "*"
  services_publish_allowed = "obs,vod"
  introduction             = "created by terraform"
}
`, rName, rName)
}

func testAccSMNV2TopicConfig_policies_step2(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_smn_topic_v2" "topic_1" {
  name                     = "%s"
  display_name             = "The display name of %s"
  users_publish_allowed    = "urn:csp:iam::0970d7b7d400f2470fbec00316a03560:root,urn:csp:iam::0970d7b7d400f2470fbec00316a03561:root"
  services_publish_allowed = "obs,vod,cce"
  introduction             = "created by terraform update"
}
`, rName, rName)
}

func testAccSMNV2TopicConfig_policies_step3(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_smn_topic_v2" "topic_1" {
  name         = "%s"
  display_name = "The display name of %s"
}
`, rName, rName)
}
