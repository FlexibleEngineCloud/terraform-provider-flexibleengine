package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/lts/huawei/logstreams"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccLTSTopicV2_basic(t *testing.T) {
	var topic logstreams.LogStream
	rand := acctest.RandString(5)
	resourceName := "flexibleengine_lts_topic.topic_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLTSTopicV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLTSTopicV2_basic(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLTSTopicV2Exists(resourceName, &topic),
					resource.TestCheckResourceAttr(resourceName, "topic_name", fmt.Sprintf("testacc_topic-%s", rand)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccLTSTopicV2ImportStateIDFunc(resourceName),
			},
		},
	})
}

func testAccCheckLTSTopicV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	ltsclient, err := config.LtsV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine LTS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_lts_topic" {
			continue
		}

		groupID := rs.Primary.Attributes["group_id"]
		_, err = logstreams.List(ltsclient, groupID).Extract()
		if err == nil {
			return fmt.Errorf("LTS topic %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckLTSTopicV2Exists(n string, topic *logstreams.LogStream) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		ltsclient, err := config.LtsV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine LTS client: %s", err)
		}

		groupID := rs.Primary.Attributes["group_id"]
		streams, err := logstreams.List(ltsclient, groupID).Extract()
		if err != nil {
			return fmt.Errorf("LTS topic get list err: %s", err)
		}
		for _, logstream := range streams.LogStreams {
			if logstream.ID == rs.Primary.ID {
				*topic = logstream
				return nil
			}
		}

		return fmt.Errorf("resource not found: lts topic %s", rs.Primary.ID)
	}
}

func testAccLTSTopicV2ImportStateIDFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		ltsGroup, ok := s.RootModule().Resources["flexibleengine_lts_group.group_1"]
		if !ok {
			return "", fmt.Errorf("LTS group not found")
		}
		ltsTopic, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("LTS topic not found")
		}

		if ltsGroup.Primary.ID == "" || ltsTopic.Primary.ID == "" {
			return "", fmt.Errorf("resource not found: %s/%s", ltsGroup.Primary.ID, ltsTopic.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", ltsGroup.Primary.ID, ltsTopic.Primary.ID), nil
	}
}

func testAccLTSTopicV2_basic(name string) string {
	return fmt.Sprintf(`
resource "flexibleengine_lts_group" "group_1" {
  group_name = "testacc_group-%s"
}

resource "flexibleengine_lts_topic" "topic_1" {
  group_id   = flexibleengine_lts_group.group_1.id
  topic_name = "testacc_topic-%s"
}
`, name, name)
}
