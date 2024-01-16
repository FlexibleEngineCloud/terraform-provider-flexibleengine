package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/vpcep/v1/endpoints"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccVPCEndpoint_Basic(t *testing.T) {
	var endpoint endpoints.Endpoint

	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "flexibleengine_vpcep_endpoint.test"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&endpoint,
		getVpcepEndpointResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVPCEndpoint_Basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "status", "accepted"),
					resource.TestCheckResourceAttr(resourceName, "enable_dns", "true"),
					resource.TestCheckResourceAttr(resourceName, "service_type", "interface"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "tf-acc"),
					resource.TestCheckResourceAttrSet(resourceName, "service_name"),
					resource.TestCheckResourceAttrSet(resourceName, "private_domain_name"),
				),
			},
			{
				Config: testAccVPCEndpoint_Update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", "accepted"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "tf-acc-update"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
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

func getVpcepEndpointResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	vpcepClient, err := conf.VPCEPClient(OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating VPCEP client: %s", err)
	}

	return endpoints.Get(vpcepClient, state.Primary.ID).Extract()
}

const testAccCompute_data = `
data "flexibleengine_availability_zones" "test" {}

data "flexibleengine_compute_flavors_v2" "test" {
  availability_zone = data.flexibleengine_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core          = 2
  memory_size       = 4
}

data "flexibleengine_images_image_v2" "test" {
  name = "OBS Ubuntu 20.04"
}

data "flexibleengine_networking_secgroup_v2" "test" {
  name = "default"
}
`

func testAccVPCEndpoint_Precondition(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "flexibleengine_vpc_v1" "test" {
  name = "%[2]s"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "test" {
  name       = "%[2]s"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = flexibleengine_vpc_v1.test.id
}

resource "flexibleengine_compute_instance_v2" "ecs" {
  name              = "%[2]s"
  image_id          = data.flexibleengine_images_image_v2.test.id
  flavor_id         = data.flexibleengine_compute_flavors_v2.test.flavors[0]
  security_groups   = [data.flexibleengine_networking_secgroup_v2.test.name]
  availability_zone = data.flexibleengine_availability_zones.test.names[0]

  network {
    uuid = flexibleengine_vpc_subnet_v1.test.id
  }
}

resource "flexibleengine_vpcep_service" "test" {
  name        = "%[2]s"
  server_type = "VM"
  vpc_id      = flexibleengine_vpc_v1.test.id
  port_id     = flexibleengine_compute_instance_v2.ecs.network[0].port
  approval    = false

  port_mapping {
    service_port  = 8080
    terminal_port = 80
  }
  tags = {
    owner = "tf-acc"
  }
}
`, testAccCompute_data, rName)
}

func testAccVPCEndpoint_Basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_vpcep_endpoint" "test" {
  service_id  = flexibleengine_vpcep_service.test.id
  vpc_id      = flexibleengine_vpc_v1.test.id
  network_id  = flexibleengine_vpc_subnet_v1.test.id
  enable_dns  = true
  description = "test description"

  enable_whitelist = true
  whitelist        = ["192.168.0.0/24", "10.10.10.12"]

  tags = {
    owner = "tf-acc"
  }
}
`, testAccVPCEndpoint_Precondition(rName))
}

func testAccVPCEndpoint_Update(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_vpcep_endpoint" "test" {
  service_id  = flexibleengine_vpcep_service.test.id
  vpc_id      = flexibleengine_vpc_v1.test.id
  network_id  = flexibleengine_vpc_subnet_v1.test.id
  enable_dns  = true
  description = "test description2"

  enable_whitelist = true
  whitelist        = ["192.168.0.0/24", "10.10.10.13"]

  tags = {
    owner = "tf-acc-update"
    foo   = "bar"
  }
}
`, testAccVPCEndpoint_Precondition(rName))
}
