package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/elb/v3/logtanks"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getElbLogTankResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.ElbV3Client(OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ELB client: %s", err)
	}
	return logtanks.Get(client, state.Primary.ID).Extract()
}

func TestAccElbLogTank_basic(t *testing.T) {
	var logTanks logtanks.LogTank
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "flexibleengine_elb_logtank.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&logTanks,
		getElbLogTankResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbLogTankConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "log_group_id",
						"flexibleengine_lts_group.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "log_topic_id",
						"flexibleengine_lts_topic.test", "id"),
				),
			},
			{
				Config: testAccElbLogTankConfig_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "log_group_id",
						"flexibleengine_lts_group.test_update", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "log_topic_id",
						"flexibleengine_lts_topic.test_update", "id"),
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

func testAccElbLogTankConfig_base(rName, updateName string) string {
	return fmt.Sprintf(`
data "flexibleengine_availability_zones" "test" {}

resource "flexibleengine_vpc_v1" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "test" {
  name        = "%[1]s"
  cidr        = "192.168.0.0/24"
  gateway_ip  = "192.168.0.1"
  vpc_id      = flexibleengine_vpc_v1.test.id
  ipv6_enable = true
}

resource "flexibleengine_lb_loadbalancer_v3" "test" {
  name            = "%[1]s"
  ipv4_subnet_id  = flexibleengine_vpc_subnet_v1.test.ipv4_subnet_id
  ipv6_network_id = flexibleengine_vpc_subnet_v1.test.id

  availability_zone = [
    data.flexibleengine_availability_zones.test.names[0]
  ]
}

resource "flexibleengine_lts_group" "%[2]s" {
  group_name = "%[2]s"
}

resource "flexibleengine_lts_topic" "%[2]s" {
  group_id   = flexibleengine_lts_group.%[2]s.id
  topic_name = "%[1]s"
}
`, rName, updateName)
}

func testAccElbLogTankConfig_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_elb_logtank" "test" {
  loadbalancer_id = flexibleengine_lb_loadbalancer_v3.test.id
  log_group_id    = flexibleengine_lts_group.test.id
  log_topic_id    = flexibleengine_lts_topic.test.id
}
`, testAccElbLogTankConfig_base(rName, "test"))
}

func testAccElbLogTankConfig_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_elb_logtank" "test" {
  loadbalancer_id = flexibleengine_lb_loadbalancer_v3.test.id
  log_group_id    = flexibleengine_lts_group.test_update.id
  log_topic_id    = flexibleengine_lts_topic.test_update.id
}
`, testAccElbLogTankConfig_base(rName, "test_update"))
}
