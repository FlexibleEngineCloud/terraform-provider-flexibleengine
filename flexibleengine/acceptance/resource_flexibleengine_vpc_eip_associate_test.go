package acceptance

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/networking/v1/eips"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getEipResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.NetworkingV1Client(OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating networking client: %s", err)
	}
	return eips.Get(c, state.Primary.ID).Extract()
}

func TestAccEIPAssociate_basic(t *testing.T) {
	var eip eips.PublicIp
	rName := acceptance.RandomAccResourceName()
	associateName := "flexibleengine_vpc_eip_associate.test"
	resourceName := "flexibleengine_vpc_eip.test"
	partten := `^((25[0-5]|2[0-4]\d|(1\d{2}|[1-9]?\d))\.){3}(25[0-5]|2[0-4]\d|(1\d{2}|[1-9]?\d))$`

	// flexibleengine_vpc_eip_associate and flexibleengine_vpc_eip have the same ID
	// and call the same API to get resource
	rc := acceptance.InitResourceCheck(
		associateName,
		&eip,
		getEipResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEIPAssociate_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(associateName, "status", "BOUND"),
					resource.TestCheckResourceAttrPair(
						associateName, "public_ip", resourceName, "address"),
					resource.TestMatchOutput("public_ip_address", regexp.MustCompile(partten)),
				),
			},
			{
				ResourceName:      associateName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccEIPAssociate_port(t *testing.T) {
	var eip eips.PublicIp
	rName := acceptance.RandomAccResourceName()
	associateName := "flexibleengine_vpc_eip_associate.test"
	resourceName := "flexibleengine_vpc_eip.test"

	// flexibleengine_vpc_eip_associate and flexibleengine_vpc_eip have the same ID
	// and call the same API to get resource
	rc := acceptance.InitResourceCheck(
		associateName,
		&eip,
		getEipResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEIPAssociate_port(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(associateName, "status", "BOUND"),
					resource.TestCheckResourceAttrPtr(
						associateName, "port_id", &eip.PortID),
					resource.TestCheckResourceAttrPair(
						associateName, "public_ip", resourceName, "address"),
				),
			},
			{
				ResourceName:      associateName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccEIPAssociate_base(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_vpc_v1" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "test" {
  name       = "%[1]s"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = flexibleengine_vpc_v1.test.id
}

resource "flexibleengine_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    share_type  = "PER"
    size        = 5
    name        = "%[1]s"
    charge_mode = "traffic"
  }
}`, rName)
}

func testAccEIPAssociate_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_availability_zones" "test" {}

data "flexibleengine_compute_flavors_v2" "test" {
  availability_zone = data.flexibleengine_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core          = 2
  memory_size       = 4
}

data "flexibleengine_images_image_v2" "test" {
  name        = "OBS Ubuntu 18.04"
  most_recent = true
}

resource "flexibleengine_compute_instance_v2" "test" {
  name               = "%[2]s"
  image_id           = data.flexibleengine_images_image_v2.test.id
  flavor_id          = data.flexibleengine_compute_flavors_v2.test.flavors[0]
  availability_zone  = data.flexibleengine_availability_zones.test.names[0]
  security_groups    = ["default"]

  network {
    uuid = flexibleengine_vpc_subnet_v1.test.id
  }
}

resource "flexibleengine_vpc_eip_associate" "test" {
  public_ip  = flexibleengine_vpc_eip.test.address
  network_id = flexibleengine_compute_instance_v2.test.network[0].uuid
  fixed_ip   = flexibleengine_compute_instance_v2.test.network[0].fixed_ip_v4
}

data "flexibleengine_compute_instance_v2" "test" {
  depends_on = [flexibleengine_vpc_eip_associate.test]

  name = "%[2]s"
}

output "public_ip_address" {
  value = data.flexibleengine_compute_instance_v2.test.floating_ip
}
`, testAccEIPAssociate_base(rName), rName)
}

func testAccEIPAssociate_port(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_networking_vip_v2" "test" {
  name       = "%s"
  network_id = flexibleengine_vpc_subnet_v1.test.id
  subnet_id  = flexibleengine_vpc_subnet_v1.test.subnet_id
}

resource "flexibleengine_vpc_eip_associate" "test" {
  public_ip = flexibleengine_vpc_eip.test.address
  port_id   = flexibleengine_networking_vip_v2.test.id
}
`, testAccEIPAssociate_base(rName), rName)
}
