package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/elbaas/backendmember"
	"log"
)

func TestAccELBBackend_basic(t *testing.T) {
	var backend backendmember.Backend

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckELBBackendDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccELBBackendConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckELBBackendExists("flexibleengine_elb_backend.backend_flexibleengine_acctest", &backend),
				),
			},
		},
	})
}

func testAccCheckELBBackendDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.otcV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_elb_backendmember" {
			continue
		}

		_, err := backendmember.Get(client, rs.Primary.Attributes["listener_id"], rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Backend member still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckELBBackendExists(n string, backend *backendmember.Backend) resource.TestCheckFunc {
	return func(s *terraform.State) error {
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

		founds, err := backendmember.Get(client, rs.Primary.Attributes["listener_id"], rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		log.Printf("[DEBUG] testAccCheckELBBackendExists found %+v.\n", founds)
		found := founds[0]
		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Backend member not found")
		}

		*backend = found

		return nil
	}
}

var TestAccELBBackendConfig_basic = fmt.Sprintf(`
resource "flexibleengine_elb_loadbalancer" "lb_flexibleengine_acctest" {
  name = "lb_flexibleengine_acctest"
  vpc_id = "%s"
  type = "External"
  bandwidth = 5
}

resource "flexibleengine_elb_listener" "ls_flexibleengine_acctest" {
  name = "ls_flexibleengine_acctest"
  protocol = "TCP"
  protocol_port = 8080
  backend_protocol = "TCP"
  backend_port = 8080
  lb_algorithm = "roundrobin"
  loadbalancer_id = "${flexibleengine_elb_loadbalancer.lb_flexibleengine_acctest.id}"
}


resource "flexibleengine_elb_health" "health_flexibleengine_acctest" {
  listener_id = "${flexibleengine_elb_listener.ls_flexibleengine_acctest.id}"
  healthcheck_protocol = "HTTP"
  healthy_threshold = 3
  healthcheck_timeout = 10
  healthcheck_interval = 5

  timeouts {
    create = "5m"
    delete = "5m"
  }
}

resource "flexibleengine_networking_secgroup_v2" "secgroup_flexibleengine_acc_test" {
  name = "secgroup_flexibleengine_acc_test"
  description = "flexibleengine security group acceptance test"
}

resource "flexibleengine_compute_instance_v2" "instance_flexibleengine_backend_test" {
  name = "instance_flexibleengine_backend_test"
  security_groups = ["${flexibleengine_networking_secgroup_v2.secgroup_flexibleengine_acc_test.name}"]
  network {
    uuid = "%s"
  }
}

resource "flexibleengine_elb_backend" "backend_flexibleengine_acctest" {
  address = "${flexibleengine_compute_instance_v2.instance_flexibleengine_backend_test.access_ip_v4}"
  listener_id = "${flexibleengine_elb_listener.ls_flexibleengine_acctest.id}"
  server_id = "${flexibleengine_compute_instance_v2.instance_flexibleengine_backend_test.id}"
  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`, OS_VPC_ID, OS_NETWORK_ID)
