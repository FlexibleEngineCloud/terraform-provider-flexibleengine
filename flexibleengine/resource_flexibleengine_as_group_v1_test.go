package flexibleengine

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/autoscaling/v1/groups"
)

func TestAccASV1Group_basic(t *testing.T) {
	var asGroup groups.Group
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "flexibleengine_as_group_v1.as_group"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckASV1GroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testASV1Group_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASV1GroupExists(resourceName, &asGroup),
					resource.TestCheckResourceAttr(resourceName, "desire_instance_number", "0"),
					resource.TestCheckResourceAttr(resourceName, "min_instance_number", "0"),
					resource.TestCheckResourceAttr(resourceName, "max_instance_number", "0"),
					resource.TestCheckResourceAttr(resourceName, "lbaas_listeners.0.protocol_port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "status", "INSERVICE"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
		},
	})
}

func TestAccASV1Group_forceDelete(t *testing.T) {
	var asGroup groups.Group
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "flexibleengine_as_group_v1.as_group"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckASV1GroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testASGroup_forceDelete(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASV1GroupExists(resourceName, &asGroup),
					resource.TestCheckResourceAttr(resourceName, "desire_instance_number", "2"),
					resource.TestCheckResourceAttr(resourceName, "min_instance_number", "2"),
					resource.TestCheckResourceAttr(resourceName, "max_instance_number", "5"),
					resource.TestCheckResourceAttr(resourceName, "instances.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "status", "INSERVICE"),
				),
			},
		},
	})
}

func testAccCheckASV1GroupDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	asClient, err := config.AutoscalingV1Client(OS_REGION_NAME)
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
		asClient, err := config.AutoscalingV1Client(OS_REGION_NAME)
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

func testASV1Group_base(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_networking_secgroup_v2" "secgroup" {
  name        = "sg-%[1]s"
  description = "This is a terraform test security group"
}

resource "flexibleengine_compute_keypair_v2" "test_key" {
  name       = "key-%[1]s"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB jrp-hp-pc"
}

data "flexibleengine_images_image_v2" "ubuntu" {
  name = "OBS Ubuntu 18.04"
}

resource "flexibleengine_as_configuration_v1" "test_as_config"{
  scaling_configuration_name = "cfg-%[1]s"
  instance_config {
    flavor   = "%s"
    image    = data.flexibleengine_images_image_v2.ubuntu.id
    key_name = flexibleengine_compute_keypair_v2.test_key.id
    disk {
      size        = 40
      volume_type = "SATA"
      disk_type   = "SYS"
    }
  }
}
`, rName, OS_FLAVOR_NAME)
}

func testASV1Group_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_lb_loadbalancer" "loadbalancer_1" {
  name          = "lb-%s"
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

resource "flexibleengine_as_group_v1" "as_group"{
  scaling_group_name       = "as-%s"
  scaling_configuration_id = flexibleengine_as_configuration_v1.test_as_config.id
  vpc_id                   = "%s"

  networks {
    id = "%s"
  }
  security_groups {
    id = flexibleengine_networking_secgroup_v2.secgroup.id
  }
  lbaas_listeners {
    pool_id       = flexibleengine_lb_pool.pool_1.id
    protocol_port = flexibleengine_lb_listener.listener_1.protocol_port
  }
  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testASV1Group_base(rName), rName, OS_SUBNET_ID, rName, OS_VPC_ID, OS_NETWORK_ID)
}

func testASGroup_forceDelete(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_as_group_v1" "as_group"{
  scaling_group_name       = "as-%s"
  scaling_configuration_id = flexibleengine_as_configuration_v1.test_as_config.id
  vpc_id                   = "%s"
  min_instance_number      = 2
  desire_instance_number   = 2
  max_instance_number      = 5
  force_delete             = true

  networks {
    id = "%s"
  }
  security_groups {
    id = flexibleengine_networking_secgroup_v2.secgroup.id
  }
  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testASV1Group_base(rName), rName, OS_VPC_ID, OS_NETWORK_ID)
}
