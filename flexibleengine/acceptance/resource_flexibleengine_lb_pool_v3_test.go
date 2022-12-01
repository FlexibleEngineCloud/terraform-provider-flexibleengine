package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/elb/v3/pools"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getElbPoolResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.LoadBalancerClient(OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ELB v3 client: %s", err)
	}
	resp, err := pools.Get(c, state.Primary.ID).Extract()
	if resp == nil && err == nil {
		return resp, fmt.Errorf("unable to find the pool (%s)", state.Primary.ID)
	}
	return resp, err
}

func TestAccElbV3Pool_basic(t *testing.T) {
	var pool pools.Pool
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	rNameUpdate := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "flexibleengine_lb_pool_v3.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&pool,
		getElbPoolResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3PoolConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "lb_method", "ROUND_ROBIN"),
				),
			},
			{
				Config: testAccElbV3PoolConfig_update(rName, rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "lb_method", "LEAST_CONNECTIONS"),
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

func testAccPoolV3Config_base(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_vpc_v1" "test" {
  name = "vpc_%[1]s"
  cidr = "192.168.0.0/16"
}
resource "flexibleengine_vpc_subnet_v1" "test" {
  name        = "subnet_%[1]s"
  cidr        = "192.168.0.0/24"
  gateway_ip  = "192.168.0.1"
  vpc_id      = flexibleengine_vpc_v1.test.id
  ipv6_enable = true
}
resource "flexibleengine_lb_loadbalancer_v3" "test" {
  name            = "%[1]s"
  ipv4_subnet_id  = flexibleengine_vpc_subnet_v1.test.subnet_id
  ipv6_network_id = flexibleengine_vpc_subnet_v1.test.id
  availability_zone = [
    "%[2]s"
  ]
}
resource "flexibleengine_lb_listener_v3" "test" {
  name            = "%[1]s"
  description     = "test description"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = flexibleengine_lb_loadbalancer_v3.test.id
  forward_eip = true
  idle_timeout     = 60
  request_timeout  = 60
  response_timeout = 60
}
`, rName, OS_AVAILABILITY_ZONE)
}

func testAccElbV3PoolConfig_basic(rName string) string {
	return fmt.Sprintf(`
%s
resource "flexibleengine_lb_pool_v3" "test" {
  name        = "%s"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = flexibleengine_lb_listener_v3.test.id
}
`, testAccPoolV3Config_base(rName), rName)
}

func testAccElbV3PoolConfig_update(rName, rNameUpdate string) string {
	return fmt.Sprintf(`
%s
resource "flexibleengine_lb_pool_v3" "test" {
  name        = "%s"
  protocol    = "HTTP"
  lb_method   = "LEAST_CONNECTIONS"
  listener_id = flexibleengine_lb_listener_v3.test.id
}
`, testAccPoolV3Config_base(rName), rNameUpdate)
}
