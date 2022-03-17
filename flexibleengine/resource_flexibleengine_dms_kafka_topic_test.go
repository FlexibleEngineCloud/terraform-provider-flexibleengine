package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/dms/v1/topics"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDmsKafkaTopic_basic(t *testing.T) {
	var rName = fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))
	topicName := "flexibleengine_dms_kafka_topic.topic"
	instanceName := "flexibleengine_dms_kafka_instance.instance_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckKafkaTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsKafkaTopic_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKafkaTopicExists(topicName),
					resource.TestCheckResourceAttr(topicName, "name", rName),
					resource.TestCheckResourceAttr(topicName, "partitions", "10"),
					resource.TestCheckResourceAttr(topicName, "replicas", "3"),
					resource.TestCheckResourceAttr(topicName, "aging_time", "36"),
					resource.TestCheckResourceAttr(topicName, "sync_replication", "false"),
					resource.TestCheckResourceAttr(topicName, "sync_flushing", "false"),
				),
			},
			{
				ResourceName:      topicName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccKafkaTopicImportStateFunc(instanceName, topicName),
			},
		},
	})
}

func testAccCheckKafkaTopicDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.DmsV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating DMS client, err=%s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_dms_kafka_topic" {
			continue
		}

		instanceID := rs.Primary.Attributes["instance_id"]
		allTopics, err := topics.List(client, instanceID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil
			}
			return fmt.Errorf("Error listing DMS kafka topics in %s, error: %s", instanceID, err)
		}

		topicID := rs.Primary.ID
		for _, item := range allTopics {
			if item.Name == topicID {
				return fmt.Errorf("DMS kafka topic %s still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckKafkaTopicExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		client, err := config.DmsV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating DMS client, err=%s", err)
		}

		instanceID := rs.Primary.Attributes["instance_id"]
		allTopics, err := topics.List(client, instanceID).Extract()
		if err != nil {
			return fmt.Errorf("Error listing DMS kafka topics in %s, error: %s", instanceID, err)
		}

		for _, item := range allTopics {
			if item.Name == rs.Primary.ID {
				return nil
			}
		}

		return fmt.Errorf("DMS kafka topic %s does not exist", rs.Primary.ID)
	}
}

// testAccKafkaTopicImportStateFunc is used to import the resource
func testAccKafkaTopicImportStateFunc(instance, topic string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		instanceState, ok := s.RootModule().Resources[instance]
		if !ok {
			return "", fmt.Errorf("DMS kafka instance not found")
		}

		topicState, ok := s.RootModule().Resources[topic]
		if !ok {
			return "", fmt.Errorf("DMS kafka topic not found")
		}

		return fmt.Sprintf("%s/%s", instanceState.Primary.ID, topicState.Primary.ID), nil
	}
}

func testAccDmsKafkaTopic_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_dms_kafka_topic" "topic" {
  instance_id = flexibleengine_dms_kafka_instance.instance_1.id
  name        = "%s"
  partitions  = 10
  aging_time  = 36
}
`, testAccDmsKafkaInstance_basic(rName), rName)
}
