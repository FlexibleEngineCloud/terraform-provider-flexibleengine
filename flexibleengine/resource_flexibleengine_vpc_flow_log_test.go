package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/golangsdk/openstack/networking/v1/flowlogs"
)

func TestAccVpcFlowLogV1_basic(t *testing.T) {
	var flowlog flowlogs.FlowLog
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	rNameUpdate := rName + "-updated"
	resourceName := "flexibleengine_vpc_flow_log_v1.flow_log"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcFlowLogV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcFlowLogV1_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcFlowLogV1Exists(resourceName, &flowlog),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "created by terraform testacc"),
					resource.TestCheckResourceAttr(resourceName, "resource_type", "port"),
					resource.TestCheckResourceAttr(resourceName, "traffic_type", "all"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
				),
			},
			{
				Config: testAccVpcFlowLogV1_update(rName, rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "description", "updated by terraform testacc"),
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

func testAccCheckVpcFlowLogV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	vpcClient, err := config.NetworkingV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating FlexibleEngine vpc client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_vpc_flow_log_v1" {
			continue
		}

		_, err := flowlogs.Get(vpcClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("VPC flow log still exists")
		}
	}

	return nil
}

func testAccCheckVpcFlowLogV1Exists(n string, flowlog *flowlogs.FlowLog) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		vpcClient, err := config.NetworkingV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating FlexibleEngine Vpc client: %s", err)
		}

		found, err := flowlogs.Get(vpcClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("VPC flow log not found")
		}

		*flowlog = *found

		return nil
	}
}

func testAccVpcFlowLogConfigBase(name string) string {
	return fmt.Sprintf(`
data "flexibleengine_images_image_v2" "ubuntu" {
  name = "OBS Ubuntu 18.04"
}

resource "flexibleengine_lts_group" "log_group1" {
  group_name = "group-%s"
}

resource "flexibleengine_lts_topic" "log_topic1" {
  group_id   = flexibleengine_lts_group.log_group1.id
  topic_name = "topic-%s"
}

resource "flexibleengine_vpc_v1" "vpc_1" {
  name = "vpc-%s"
  cidr = "172.16.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "subnet_1" {
  vpc_id     = flexibleengine_vpc_v1.vpc_1.id
  name       = "subnet-%s"
  cidr       = "172.16.0.0/24"
  gateway_ip = "172.16.0.1"
}

resource "flexibleengine_compute_instance_v2" "instance_1" {
  name              = "ecs-%s"
  image_id          = data.flexibleengine_images_image_v2.ubuntu.id
  flavor_name       = "s3.small.1"
  security_groups   = ["default"]
  availability_zone = "%s"

  network {
    uuid = flexibleengine_vpc_subnet_v1.subnet_1.id
  }
}
`, name, name, name, name, name, OS_AVAILABILITY_ZONE)
}

func testAccVpcFlowLogV1_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_vpc_flow_log_v1" "flow_log" {
  name          = "%s"
  description   = "created by terraform testacc"
  resource_id   = flexibleengine_compute_instance_v2.instance_1.network[0].port
  log_group_id  = flexibleengine_lts_group.log_group1.id
  log_topic_id  = flexibleengine_lts_topic.log_topic1.id
}
`, testAccVpcFlowLogConfigBase(name), name)
}

func testAccVpcFlowLogV1_update(base, name string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_vpc_flow_log_v1" "flow_log" {
  name          = "%s"
  description   = "updated by terraform testacc"
  resource_id   = flexibleengine_compute_instance_v2.instance_1.network[0].port
  log_group_id  = flexibleengine_lts_group.log_group1.id
  log_topic_id  = flexibleengine_lts_topic.log_topic1.id
}
`, testAccVpcFlowLogConfigBase(base), name)
}
