package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/elb/v2/whitelists"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getWhitelistResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.LoadBalancerClient(OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating FlexibleEngine LB v2 client: %s", err)
	}
	resp, err := whitelists.Get(c, state.Primary.ID).Extract()
	if resp == nil && err == nil {
		return resp, fmt.Errorf("Unable to find the whitelist (%s)", state.Primary.ID)
	}
	return resp, err
}

func TestAccLBV2Whitelist_basic(t *testing.T) {
	var whitelist whitelists.Whitelist
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "flexibleengine_lb_whitelist.whitelist_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&whitelist,
		getWhitelistResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccLBV2WhitelistConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				Config: testAccLBV2WhitelistConfig_update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "enable_whitelist", "true"),
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

func testAccLBV2WhitelistConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_vpc_v1" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "test" {
  name       = "%[1]s"
  cidr       = "192.168.0.0/24"
  vpc_id     = flexibleengine_vpc_v1.test.id
  gateway_ip = "192.168.0.1"
}
resource "flexibleengine_lb_loadbalancer" "loadbalancer_1" {
  name          = "%[1]s"
  vip_subnet_id = flexibleengine_vpc_subnet_v1.test.ipv4_subnet_id
}

resource "flexibleengine_lb_listener" "listener_1" {
  name            = "%[1]s"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = flexibleengine_lb_loadbalancer.loadbalancer_1.id
}

resource "flexibleengine_lb_whitelist" "whitelist_1" {
  enable_whitelist = true
  whitelist        = "192.168.11.1,192.168.0.1/24"
  listener_id      = flexibleengine_lb_listener.listener_1.id
}
`, rName)
}

func testAccLBV2WhitelistConfig_update(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_vpc_v1" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "test" {
  name       = "%[1]s"
  cidr       = "192.168.0.0/24"
  vpc_id     = flexibleengine_vpc_v1.test.id
  gateway_ip = "192.168.0.1"
}

resource "flexibleengine_lb_loadbalancer" "loadbalancer_1" {
  name          = "%[1]s"
  vip_subnet_id = flexibleengine_vpc_subnet_v1.test.ipv4_subnet_id
}

resource "flexibleengine_lb_listener" "listener_1" {
  name            = "%[1]s"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = flexibleengine_lb_loadbalancer.loadbalancer_1.id
}

resource "flexibleengine_lb_whitelist" "whitelist_1" {
  enable_whitelist = true
  whitelist        = "192.168.11.1,192.168.0.1/24,192.168.201.18/8"
  listener_id      = flexibleengine_lb_listener.listener_1.id
}
`, rName)
}
