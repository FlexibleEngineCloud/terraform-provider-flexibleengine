package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/huaweicloud/golangsdk/openstack/autoscaling/v1/groups"
	"log"
)

func TestAccASV1Group_basic(t *testing.T) {
	var asGroup groups.Group

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckASV1GroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testASV1Group_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASV1GroupExists("flexibleengine_as_group_v1.hth_as_group", &asGroup),
					resource.TestCheckResourceAttr(
						"flexibleengine_as_group_v1.hth_as_group", "lbaas_listeners.0.protocol_port", "8080"),
				),
			},
		},
	})
}

func testAccCheckASV1GroupDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	asClient, err := config.autoscalingV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating flexibleengine autoscaling client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_as_group_v1" {
			continue
		}

		_, err := groups.Get(asClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("AS group still exists")
		}
	}

	log.Printf("[DEBUG] testCheckASV1GroupDestroy success!")

	return nil
}

func testAccCheckASV1GroupExists(n string, group *groups.Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		asClient, err := config.autoscalingV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating flexibleengine autoscaling client: %s", err)
		}

		found, err := groups.Get(asClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Autoscaling Group not found")
		}
		log.Printf("[DEBUG] test found is: %#v", found)
		group = &found

		return nil
	}
}

var testASV1Group_basic = fmt.Sprintf(`
resource "flexibleengine_networking_secgroup_v2" "secgroup" {
  name        = "terraform"
  description = "This is a terraform test security group"
}

resource "flexibleengine_compute_keypair_v2" "hth_key" {
  name = "hth_key"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB jrp-hp-pc"
}

resource "flexibleengine_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = "2c0a74a9-4395-4e62-a17b-e3e86fbf66b7"
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
}

resource "flexibleengine_as_configuration_v1" "hth_as_config"{
  scaling_configuration_name = "hth_as_config"
  instance_config {
    image = "%s"
    disk {
      size = 40
      volume_type = "SATA"
      disk_type = "SYS"
    }
    key_name = "${flexibleengine_compute_keypair_v2.hth_key.id}"
  }
}

resource "flexibleengine_as_group_v1" "hth_as_group"{
  scaling_group_name = "hth_as_group"
  scaling_configuration_id = "${flexibleengine_as_configuration_v1.hth_as_config.id}"
  networks {
    id = "%s"
  }
  security_groups {
    id = "${flexibleengine_networking_secgroup_v2.secgroup.id}"
  }
  lbaas_listeners {
    pool_id = "${flexibleengine_lb_pool_v2.pool_1.id}"
    protocol_port = "${flexibleengine_lb_listener_v2.listener_1.protocol_port}"
  }
  vpc_id = "%s"
}
`, OS_IMAGE_ID, OS_NETWORK_ID, OS_VPC_ID)
