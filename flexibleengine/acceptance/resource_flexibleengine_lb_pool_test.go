package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/elb/v2/pools"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getPoolResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.LoadBalancerClient(OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating flexibleengine LB v2 client: %s", err)
	}
	resp, err := pools.Get(c, state.Primary.ID).Extract()
	if resp == nil && err == nil {
		return resp, fmt.Errorf("unable to find the pool (%s)", state.Primary.ID)
	}
	return resp, err
}

func TestAccLBV2Pool_basic(t *testing.T) {
	var pool pools.Pool
	rName := acceptance.RandomAccResourceNameWithDash()
	rNameUpdate := acceptance.RandomAccResourceNameWithDash()
	resourceName := "flexibleengine_lb_pool.pool_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&pool,
		getPoolResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccLBV2PoolConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "lb_method", "ROUND_ROBIN"),
				),
			},
			{
				Config: testAccLBV2PoolConfig_update(rName, rNameUpdate),
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

func testAccLBV2PoolConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_lb_loadbalancer" "loadbalancer_1" {
  name          = "%[1]s"
  vip_subnet_id = "%[2]s"
}

resource "flexibleengine_lb_listener" "listener_1" {
  name            = "%[1]s"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = flexibleengine_lb_loadbalancer.loadbalancer_1.id
}

resource "flexibleengine_lb_pool" "pool_1" {
  name        = "%[1]s"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = flexibleengine_lb_listener.listener_1.id
}
`, rName, OS_SUBNET_ID)
}

func testAccLBV2PoolConfig_update(rName, rNameUpdate string) string {
	return fmt.Sprintf(`
resource "flexibleengine_lb_loadbalancer" "loadbalancer_1" {
  name          = "%[1]s"
  vip_subnet_id = "%[3]s"
}

resource "flexibleengine_lb_listener" "listener_1" {
  name            = "%[1]s"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = flexibleengine_lb_loadbalancer.loadbalancer_1.id
}

resource "flexibleengine_lb_pool" "pool_1" {
  name           = "%[2]s"
  protocol       = "HTTP"
  lb_method      = "LEAST_CONNECTIONS"
  admin_state_up = "true"
  listener_id    = flexibleengine_lb_listener.listener_1.id
}
`, rName, rNameUpdate, OS_SUBNET_ID)
}
