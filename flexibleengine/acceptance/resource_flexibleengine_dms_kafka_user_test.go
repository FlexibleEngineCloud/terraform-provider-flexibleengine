package acceptance

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/kafka/v2/model"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getDmsKafkaUserFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.HcDmsV2Client(OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DMS client: %s", err)
	}

	// Split instance_id and user from resource id
	parts := strings.SplitN(state.Primary.ID, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, must be <instance_id>/<user>")
	}
	instanceId := parts[0]
	instanceUser := parts[1]

	// List all instance users
	request := &model.ShowInstanceUsersRequest{
		InstanceId: instanceId,
	}

	response, err := client.ShowInstanceUsers(request)
	if err != nil {
		return nil, fmt.Errorf("error listing DMS kafka users in %s, error: %s", instanceId, err)
	}
	if response.Users != nil && len(*response.Users) != 0 {
		users := *response.Users
		for _, user := range users {
			if *user.UserName == instanceUser {
				return user, nil
			}
		}
	}

	return nil, fmt.Errorf("can not found DMS kafka user")
}

func TestAccDmsKafkaUser_basic(t *testing.T) {
	var user model.ShowInstanceUsersEntity
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "flexibleengine_dms_kafka_user.test"
	password := acceptance.RandomPassword()
	passwordUpdate := password + "update"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&user,
		getDmsKafkaUserFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDmsKafkaUser_basic(rName, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
			{
				Config: testAccDmsKafkaUser_basic(rName, passwordUpdate),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testAccDmsKafkaInstance_base(resName string) string {
	return fmt.Sprintf(`
data "flexibleengine_dms_product" "product_1" {
  bandwidth = "300MB"
}

resource "flexibleengine_vpc_v1" "vpc_1" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "vpc_subnet_1" {
  name       = "%s"
  cidr       = "192.168.10.0/24"
  gateway_ip = "192.168.10.1"
  vpc_id     = flexibleengine_vpc_v1.vpc_1.id
}

resource "flexibleengine_networking_secgroup_v2" "secgroup_1" {
  name        = "%s"
  description = "secgroup for DMS"
}`, resName, resName, resName)
}

func testAccDmsKafkaInstance_basic(resName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_dms_kafka_instance" "instance_1" {
  name               = "%s"
  manager_user       = "admin"
  manager_password   = "Dmstest@123"
  access_user        = "user"
  password           = "Dmstest@123"
  vpc_id             = flexibleengine_vpc_v1.vpc_1.id
  network_id         = flexibleengine_vpc_subnet_v1.vpc_subnet_1.id
  security_group_id  = flexibleengine_networking_secgroup_v2.secgroup_1.id
  availability_zones = data.flexibleengine_dms_product.product_1.availability_zones
  bandwidth          = data.flexibleengine_dms_product.product_1.bandwidth
  product_id         = data.flexibleengine_dms_product.product_1.id
  storage_space      = data.flexibleengine_dms_product.product_1.storage_space
  engine_version     = data.flexibleengine_dms_product.product_1.engine_version
}
`, testAccDmsKafkaInstance_base(resName), resName)
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

func testAccDmsKafkaUser_basic(rName, password string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_dms_kafka_user" "test" {
  instance_id = flexibleengine_dms_kafka_instance.instance_1.id
  name        = "%s"
  password    = "%s"
}
`, testAccDmsKafkaTopic_basic(rName), rName, password)
}
