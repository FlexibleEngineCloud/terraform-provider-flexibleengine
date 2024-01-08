package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/lbaas_v2/pools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccLBV2Pool_basic(t *testing.T) {
	var pool pools.Pool

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBV2PoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBV2PoolConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2PoolExists("flexibleengine_lb_pool.pool_1", &pool),
					resource.TestCheckResourceAttr("flexibleengine_lb_pool.pool_1", "name", "pool_1"),
				),
			},
			{
				Config: testAccLBV2PoolConfig_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("flexibleengine_lb_pool.pool_1", "name", "pool_1_updated"),
				),
			},
		},
	})
}

func testAccCheckLBV2PoolDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	lbClient, err := config.ElbV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine ELB v2.0 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_lb_pool" {
			continue
		}

		_, err := pools.Get(lbClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Pool still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBV2PoolExists(n string, pool *pools.Pool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		lbClient, err := config.ElbV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine ELB v2.0 client: %s", err)
		}

		found, err := pools.Get(lbClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Member not found")
		}

		*pool = *found

		return nil
	}
}

var testAccLBV2PoolConfig_basic = fmt.Sprintf(`
resource "flexibleengine_lb_loadbalancer" "loadbalancer_1" {
  name          = "loadbalancer_1"
  vip_subnet_id = "%s"
}

resource "flexibleengine_lb_listener" "listener_1" {
  name            = "listener_1"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = flexibleengine_lb_loadbalancer.loadbalancer_1.id
}

resource "flexibleengine_lb_pool" "pool_1" {
  name        = "pool_1"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = flexibleengine_lb_listener.listener_1.id
}
`, OS_SUBNET_ID)

var testAccLBV2PoolConfig_update = fmt.Sprintf(`
resource "flexibleengine_lb_loadbalancer" "loadbalancer_1" {
  name          = "loadbalancer_1"
  vip_subnet_id = "%s"
}

resource "flexibleengine_lb_listener" "listener_1" {
  name            = "listener_1"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = flexibleengine_lb_loadbalancer.loadbalancer_1.id
}

resource "flexibleengine_lb_pool" "pool_1" {
  name           = "pool_1_updated"
  protocol       = "HTTP"
  lb_method      = "LEAST_CONNECTIONS"
  admin_state_up = "true"
  listener_id    = flexibleengine_lb_listener.listener_1.id
}
`, OS_SUBNET_ID)
