package flexibleengine

import (
	"fmt"
	"testing"

	"log"

	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/elbaas/healthcheck"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccELBHealth_basic(t *testing.T) {
	var health healthcheck.Health

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckELBHealthDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccELBHealthConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckELBHealthExists(t, "flexibleengine_elb_health.health_1", &health),
					resource.TestCheckResourceAttr("flexibleengine_elb_health.health_1", "healthy_threshold", "3"),
					resource.TestCheckResourceAttr("flexibleengine_elb_health.health_1", "healthcheck_timeout", "10"),
				),
			},
			{
				Config: TestAccELBHealthConfig_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("flexibleengine_elb_health.health_1", "healthy_threshold", "5"),
					resource.TestCheckResourceAttr("flexibleengine_elb_health.health_1", "healthcheck_timeout", "15"),
				),
			},
		},
	})
}

func testAccCheckELBHealthDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.otcV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_elb_healthcheck" {
			continue
		}

		_, err := healthcheck.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Health still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckELBHealthExists(t *testing.T, n string, health *healthcheck.Health) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		log.Printf("[DEBUG] testAccCheckELBHealthExists resources %+v.\n", s.RootModule().Resources)
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		client, err := config.otcV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine networking client: %s", err)
		}

		found, err := healthcheck.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		log.Printf("[DEBUG] testAccCheckELBHealthExists found %+v.\n", found)

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Health not found")
		}

		*health = *found

		return nil
	}
}

var TestAccELBHealthConfig_basic = fmt.Sprintf(`
resource "flexibleengine_elb_loadbalancer" "loadbalancer_1" {
  name = "loadbalancer_1"
  vpc_id = "%s"
  type = "External"
  bandwidth = 5
}

resource "flexibleengine_elb_listener" "listener_1" {
  name = "listener_1"
  protocol = "TCP"
  protocol_port = 8080
  backend_protocol = "TCP"
  backend_port = 8080
  lb_algorithm = "roundrobin"
  loadbalancer_id = "${flexibleengine_elb_loadbalancer.loadbalancer_1.id}"
}


resource "flexibleengine_elb_health" "health_1" {
  listener_id = "${flexibleengine_elb_listener.listener_1.id}"
  healthcheck_protocol = "HTTP"
  healthy_threshold = 3
  healthcheck_timeout = 10
  healthcheck_interval = 5

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, OS_VPC_ID)

var TestAccELBHealthConfig_update = fmt.Sprintf(`
resource "flexibleengine_elb_loadbalancer" "loadbalancer_1" {
  name = "loadbalancer_1"
  vpc_id = "%s"
  type = "External"
  bandwidth = 5
}

resource "flexibleengine_elb_listener" "listener_1" {
  name = "listener_1"
  protocol = "TCP"
  protocol_port = 8080
  backend_protocol = "TCP"
  backend_port = 8080
  lb_algorithm = "roundrobin"
  loadbalancer_id = "${flexibleengine_elb_loadbalancer.loadbalancer_1.id}"
}


resource "flexibleengine_elb_health" "health_1" {
  listener_id = "${flexibleengine_elb_listener.listener_1.id}"
  healthcheck_protocol = "HTTP"
  healthy_threshold = 5
  healthcheck_timeout = 15
  healthcheck_interval = 3

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, OS_VPC_ID)
