package acceptance

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getDmsRocketMQTopicResourceFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := OS_REGION_NAME
	// getRocketmqTopic: query DMS rocketmq topic
	var (
		getRocketmqTopicHttpUrl = "v2/{project_id}/instances/{instance_id}/topics/{topic}"
		getRocketmqTopicProduct = "dms"
	)
	getRocketmqTopicClient, err := config.NewServiceClient(getRocketmqTopicProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DmsRocketMQTopic Client: %s", err)
	}

	parts := strings.SplitN(state.Primary.ID, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, must be <instance_id>/<topic>")
	}
	instanceID := parts[0]
	topic := parts[1]
	getRocketmqTopicPath := getRocketmqTopicClient.Endpoint + getRocketmqTopicHttpUrl
	getRocketmqTopicPath = strings.ReplaceAll(getRocketmqTopicPath, "{project_id}", getRocketmqTopicClient.ProjectID)
	getRocketmqTopicPath = strings.ReplaceAll(getRocketmqTopicPath, "{instance_id}", instanceID)
	getRocketmqTopicPath = strings.ReplaceAll(getRocketmqTopicPath, "{topic}", topic)

	getRocketmqTopicOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getRocketmqTopicResp, err := getRocketmqTopicClient.Request("GET", getRocketmqTopicPath, &getRocketmqTopicOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DmsRocketMQTopic: %s", err)
	}
	return utils.FlattenResponse(getRocketmqTopicResp)
}

func TestAccDmsRocketMQTopic_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "flexibleengine_dms_rocketmq_topic.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDmsRocketMQTopicResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDmsRocketMQTopic_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "total_read_queue_num", "3"),
					resource.TestCheckResourceAttr(rName, "total_write_queue_num", "3"),
					resource.TestCheckResourceAttr(rName, "brokers.0.write_queue_num", "3"),
					resource.TestCheckResourceAttr(rName, "brokers.0.write_queue_num", "3"),
				),
			},
			{
				Config: testDmsRocketMQTopic_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "permission", "sub"),
					resource.TestCheckResourceAttr(rName, "total_read_queue_num", "4"),
					resource.TestCheckResourceAttr(rName, "total_write_queue_num", "5"),
					resource.TestCheckResourceAttr(rName, "brokers.0.read_queue_num", "4"),
					resource.TestCheckResourceAttr(rName, "brokers.0.write_queue_num", "5"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_id"},
			},
		},
	})
}

func testAccDmsRocketmqTopic_Base(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_vpc_v1" "test" {
  name        = "%[1]s"
  cidr        = "192.168.0.0/24"
  description = "Test for DMS RocketMQ"
}

resource "flexibleengine_vpc_subnet_v1" "test" {
  name       = "%[1]s"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = flexibleengine_vpc_v1.test.id
}

resource "flexibleengine_networking_secgroup_v2" "test" {
  name        = "%[1]s"
  description = "secgroup for rocketmq"
}

data "flexibleengine_availability_zones" "test" {}

resource "flexibleengine_dms_rocketmq_instance" "test" {
  name                = "%[1]s"
  engine_version      = "4.8.0"
  storage_space       = 600
  vpc_id              = flexibleengine_vpc_v1.test.id
  subnet_id           = flexibleengine_vpc_subnet_v1.test.id
  security_group_id   = flexibleengine_networking_secgroup_v2.test.id
  availability_zones  = [
    data.flexibleengine_availability_zones.test.names[0]
  ]
  flavor_id           = "c6.4u8g.cluster"
  storage_spec_code   = "dms.physical.storage.high.v2"
  broker_num          = 1
}
`, rName)
}

func testDmsRocketMQTopic_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_dms_rocketmq_topic" "test" {
  instance_id = flexibleengine_dms_rocketmq_instance.test.id
  name        = "%s"
  brokers {
    name = "broker-0"
  }
}
`, testAccDmsRocketmqTopic_Base(name), name)
}

func testDmsRocketMQTopic_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_dms_rocketmq_topic" "test" {
  instance_id           = flexibleengine_dms_rocketmq_instance.test.id
  name                  = "%s"
  permission            = "sub"
  total_read_queue_num  = "4"
  total_write_queue_num = "5"
}
`, testAccDmsRocketmqTopic_Base(name), name)
}
