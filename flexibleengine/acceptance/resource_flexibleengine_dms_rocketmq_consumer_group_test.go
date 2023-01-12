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

func getDmsRocketMQConsumerGroupResourceFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := OS_REGION_NAME
	// getRocketmqConsumerGroup: query DMS rocketmq consumer group
	var (
		getRocketmqConsumerGroupHttpUrl = "v2/{project_id}/instances/{instance_id}/groups/{group}"
		getRocketmqConsumerGroupProduct = "dms"
	)
	getRocketmqConsumerGroupClient, err := config.NewServiceClient(getRocketmqConsumerGroupProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DmsRocketMQConsumerGroup Client: %s", err)
	}

	// Split instance_id and group from resource id
	parts := strings.SplitN(state.Primary.ID, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, must be <instance_id>/<consumerGroup>")
	}
	instanceID := parts[0]
	name := parts[1]
	getRocketmqConsumerGroupPath := getRocketmqConsumerGroupClient.Endpoint + getRocketmqConsumerGroupHttpUrl
	getRocketmqConsumerGroupPath = strings.ReplaceAll(getRocketmqConsumerGroupPath, "{project_id}", getRocketmqConsumerGroupClient.ProjectID)
	getRocketmqConsumerGroupPath = strings.ReplaceAll(getRocketmqConsumerGroupPath, "{instance_id}", instanceID)
	getRocketmqConsumerGroupPath = strings.ReplaceAll(getRocketmqConsumerGroupPath, "{group}", name)

	getRocketmqConsumerGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getRocketmqConsumerGroupResp, err := getRocketmqConsumerGroupClient.Request("GET", getRocketmqConsumerGroupPath, &getRocketmqConsumerGroupOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DmsRocketMQConsumerGroup: %s", err)
	}
	return utils.FlattenResponse(getRocketmqConsumerGroupResp)
}

func TestAccDmsRocketMQConsumerGroup_basic(t *testing.T) {
	var obj interface{}

	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "flexibleengine_dms_rocketmq_consumer_group.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDmsRocketMQConsumerGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDmsRocketMQConsumerGroup_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "broadcast", "true"),
					resource.TestCheckResourceAttr(resourceName, "retry_max_times", "3"),
				),
			},
			{
				Config: testDmsRocketMQConsumerGroup_basic_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "broadcast", "false"),
					resource.TestCheckResourceAttr(resourceName, "retry_max_times", "5"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
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

func testAccDmsRocketmqConsumerGroup_Base(rName string) string {
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

func testDmsRocketMQConsumerGroup_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_dms_rocketmq_consumer_group" "test" {
  instance_id    = flexibleengine_dms_rocketmq_instance.test.id
  broadcast      = true
  brokers        = [
    "broker-0"
  ]
  name            = "%s"
  retry_max_times = "3"
}
`, testAccDmsRocketmqConsumerGroup_Base(name), name)
}

func testDmsRocketMQConsumerGroup_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_dms_rocketmq_consumer_group" "test" {
  instance_id    = flexibleengine_dms_rocketmq_instance.test.id
  broadcast      = false
  brokers        = [
    "broker-0"
  ]
  name            = "%s"
  retry_max_times = "5"
  enabled         = false
}
`, testAccDmsRocketmqConsumerGroup_Base(name), name)
}
