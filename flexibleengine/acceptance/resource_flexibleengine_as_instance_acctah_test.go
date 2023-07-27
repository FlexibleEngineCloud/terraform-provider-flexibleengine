package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/autoscaling/v1/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getASInstanceAttachResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := OS_REGION_NAME
	client, err := cfg.AutoscalingV1Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating autoscaling client: %s", err)
	}

	groupID := state.Primary.Attributes["scaling_group_id"]
	instanceID := state.Primary.Attributes["instance_id"]
	page, err := instances.List(client, groupID, nil).AllPages()
	if err != nil {
		return nil, err
	}

	allInstances, err := page.(instances.InstancePage).Extract()
	if err != nil {
		return nil, fmt.Errorf("failed to fetching instances in AS group %s: %s", groupID, err)
	}

	for _, ins := range allInstances {
		if ins.ID == instanceID {
			return &ins, nil
		}
	}

	return nil, fmt.Errorf("can not find the instance %s in AS group %s", instanceID, groupID)
}

func TestAccASInstanceAttach_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "flexibleengine_as_instance_attach.test0"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getASInstanceAttachResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testASInstanceAttach_conf(name, "false", "false"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "scaling_group_id", "flexibleengine_as_group.acc_as_group", "id"),
					resource.TestCheckResourceAttrPair(rName, "instance_id", "flexibleengine_compute_instance_v2.test.0", "id"),
					resource.TestCheckResourceAttr(rName, "protected", "false"),
					resource.TestCheckResourceAttr(rName, "standby", "false"),
					resource.TestCheckResourceAttr(rName, "status", "INSERVICE"),
				),
			},
			{
				Config: testASInstanceAttach_conf(name, "true", "false"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "protected", "true"),
					resource.TestCheckResourceAttr(rName, "standby", "false"),
					resource.TestCheckResourceAttr(rName, "status", "INSERVICE"),
				),
			},
			{
				Config: testASInstanceAttach_conf(name, "true", "true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "protected", "true"),
					resource.TestCheckResourceAttr(rName, "standby", "true"),
					resource.TestCheckResourceAttr(rName, "status", "STANDBY"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"append_instance"},
			},
		},
	})
}

func testASGroup_Base(rName string) string {
	return fmt.Sprintf(`
data "flexibleengine_availability_zones" "test" {}

data "flexibleengine_compute_flavors_v2" "test" {
  availability_zone = data.flexibleengine_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core          = 2
  memory_size       = 4
}

data "flexibleengine_images_image" "test" {
  name = "OBS Ubuntu 18.04"
}

resource "flexibleengine_vpc_v1" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "test" {
  name       = "%[1]s"
  vpc_id     = flexibleengine_vpc_v1.test.id
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
}

resource "flexibleengine_networking_secgroup_v2" "test" {
  name                 = "%[1]s"
  delete_default_rules = true
}

resource "flexibleengine_compute_keypair_v2" "acc_key" {
  name       = "%[1]s"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB jrp-hp-pc"
}

resource "flexibleengine_lb_loadbalancer_v3" "loadbalancer_1" {
  availability_zone = [data.flexibleengine_availability_zones.test.names[0]]
  name              = "%[1]s"
  ipv4_subnet_id    = flexibleengine_vpc_subnet_v1.test.ipv4_subnet_id
}

resource "flexibleengine_lb_listener_v3" "listener_1" {
  name            = "%[1]s"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = flexibleengine_lb_loadbalancer_v3.loadbalancer_1.id
}

resource "flexibleengine_lb_pool_v3" "pool_1" {
  name        = "%[1]s"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = flexibleengine_lb_listener_v3.listener_1.id
}

resource "flexibleengine_as_configuration" "acc_as_config"{
  scaling_configuration_name = "%[1]s"
  instance_config {
	image    = data.flexibleengine_images_image.test.id
	flavor   = data.flexibleengine_compute_flavors_v2.test.flavors[0]
    key_name = flexibleengine_compute_keypair_v2.acc_key.id
    disk {
      size        = 40
      volume_type = "SSD"
      disk_type   = "SYS"
    }
  }
}`, rName)
}

func testASGroup_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_as_group" "acc_as_group"{
  scaling_group_name       = "%s"
  scaling_configuration_id = flexibleengine_as_configuration.acc_as_config.id
  vpc_id                   = flexibleengine_vpc_v1.test.id
  max_instance_number      = 5

  networks {
    id = flexibleengine_vpc_subnet_v1.test.id
  }
  security_groups {
    id = flexibleengine_networking_secgroup_v2.test.id
  }
  lbaas_listeners {
    pool_id       = flexibleengine_lb_pool_v3.pool_1.id
    protocol_port = flexibleengine_lb_listener_v3.listener_1.protocol_port
  }
  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testASGroup_Base(rName), rName)
}

func testASInstanceAttach_conf(name, protection, standby string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_compute_instance_v2" "test" {
  count           = 2
  name            = "%s-${count.index}"
  image_id        = data.flexibleengine_images_image.test.id
  flavor_id       = data.flexibleengine_compute_flavors_v2.test.flavors[0]
  security_groups = [flexibleengine_networking_secgroup_v2.test.name]

  network {
    uuid = flexibleengine_vpc_subnet_v1.test.id
  }
}

resource "flexibleengine_as_instance_attach" "test0" {
  scaling_group_id = flexibleengine_as_group.acc_as_group.id
  instance_id      = flexibleengine_compute_instance_v2.test[0].id
  protected        = %[3]s
  standby          = %[4]s
}

resource "flexibleengine_as_instance_attach" "test1" {
  scaling_group_id = flexibleengine_as_group.acc_as_group.id
  instance_id      = flexibleengine_compute_instance_v2.test[1].id
  protected        = %[3]s
  standby          = %[4]s
}
`, testASGroup_basic(name), name, protection, standby)
}
