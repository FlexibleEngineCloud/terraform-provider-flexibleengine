package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/pools"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccLBV2Pool_basic(t *testing.T) {
	var pool pools.Pool

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckLB(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBV2PoolDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TestAccLBV2PoolConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2PoolExists("flexibleengine_lb_pool_v2.pool_1", &pool),
				),
			},
			resource.TestStep{
				Config: TestAccLBV2PoolConfig_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("flexibleengine_lb_pool_v2.pool_1", "name", "pool_1_updated"),
				),
			},
		},
	})
}

func testAccCheckLBV2PoolDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_lb_pool_v2" {
			continue
		}

		_, err := pools.Get(networkingClient, rs.Primary.ID).Extract()
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
		networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine networking client: %s", err)
		}

		found, err := pools.Get(networkingClient, rs.Primary.ID).Extract()
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

var TestAccLBV2PoolConfig_basic = fmt.Sprintf(`
resource "flexibleengine_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = "%s"
}

resource "flexibleengine_lb_listener_v2" "listener_1" {
  name = "listener_1"
  protocol = "HTTP"
  protocol_port = 8080
  loadbalancer_id = "${flexibleengine_lb_loadbalancer_v2.loadbalancer_1.id}"
}

resource "flexibleengine_lb_pool_v2" "pool_1" {
  name = "pool_1"
  protocol = "HTTP"
  lb_method = "ROUND_ROBIN"
  listener_id = "${flexibleengine_lb_listener_v2.listener_1.id}"

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, OS_SUBNET_ID)

var TestAccLBV2PoolConfig_update = fmt.Sprintf(`
resource "flexibleengine_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = "%s"
}

resource "flexibleengine_lb_listener_v2" "listener_1" {
  name = "listener_1"
  protocol = "HTTP"
  protocol_port = 8080
  loadbalancer_id = "${flexibleengine_lb_loadbalancer_v2.loadbalancer_1.id}"
}

resource "flexibleengine_lb_pool_v2" "pool_1" {
  name = "pool_1_updated"
  protocol = "HTTP"
  lb_method = "LEAST_CONNECTIONS"
  admin_state_up = "true"
  listener_id = "${flexibleengine_lb_listener_v2.listener_1.id}"

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, OS_SUBNET_ID)
