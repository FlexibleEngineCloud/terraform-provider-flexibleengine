package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/lbaas_v2/monitors"
)

func TestAccLBV2Monitor_basic(t *testing.T) {
	var monitor monitors.Monitor
	rand := acctest.RandString(5)
	resourceName := "flexibleengine_lb_monitor_v2.monitor_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBV2MonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBV2MonitorConfig_basic(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2MonitorExists(t, resourceName, &monitor),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("monitor_%s", rand)),
					resource.TestCheckResourceAttr(resourceName, "delay", "20"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "10"),
					resource.TestCheckResourceAttr(resourceName, "max_retries", "5"),
					resource.TestCheckResourceAttr(resourceName, "port", "8888"),
				),
			},
			{
				Config: testAccLBV2MonitorConfig_update(rand),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("monitor_%s_updated", rand)),
					resource.TestCheckResourceAttr(resourceName, "delay", "30"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "15"),
					resource.TestCheckResourceAttr(resourceName, "max_retries", "10"),
					resource.TestCheckResourceAttr(resourceName, "port", "9999"),
				),
			},
		},
	})
}

func testAccCheckLBV2MonitorDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_lb_monitor_v2" {
			continue
		}

		_, err := monitors.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Monitor still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBV2MonitorExists(t *testing.T, n string, monitor *monitors.Monitor) resource.TestCheckFunc {
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

		found, err := monitors.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Monitor not found")
		}

		*monitor = *found

		return nil
	}
}

func testAccLBV2MonitorConfig_basic(rand string) string {
	return fmt.Sprintf(`
resource "flexibleengine_lb_loadbalancer_v2" "loadbalancer_1" {
  name          = "loadbalancer_%s"
  vip_subnet_id = "%s"
}

resource "flexibleengine_lb_listener_v2" "listener_1" {
  name            = "listener_%s"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = flexibleengine_lb_loadbalancer_v2.loadbalancer_1.id
}

resource "flexibleengine_lb_pool_v2" "pool_1" {
  name        = "pool_%s"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = flexibleengine_lb_listener_v2.listener_1.id
}

resource "flexibleengine_lb_monitor_v2" "monitor_1" {
  name        = "monitor_%s"
  type        = "PING"
  delay       = 20
  timeout     = 10
  max_retries = 5
  pool_id     = flexibleengine_lb_pool_v2.pool_1.id
  port        = 8888
}
`, rand, OS_SUBNET_ID, rand, rand, rand)
}

func testAccLBV2MonitorConfig_update(rand string) string {
	return fmt.Sprintf(`
resource "flexibleengine_lb_loadbalancer_v2" "loadbalancer_1" {
  name          = "loadbalancer_%s"
  vip_subnet_id = "%s"
}

resource "flexibleengine_lb_listener_v2" "listener_1" {
  name            = "listener_%s"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = flexibleengine_lb_loadbalancer_v2.loadbalancer_1.id
}

resource "flexibleengine_lb_pool_v2" "pool_1" {
  name        = "pool_%s"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = flexibleengine_lb_listener_v2.listener_1.id
}

resource "flexibleengine_lb_monitor_v2" "monitor_1" {
  name           = "monitor_%s_updated"
  type           = "PING"
  delay          = 30
  timeout        = 15
  max_retries    = 10
  admin_state_up = "true"
  pool_id        = flexibleengine_lb_pool_v2.pool_1.id
  port           = 9999
}
`, rand, OS_SUBNET_ID, rand, rand, rand)
}
