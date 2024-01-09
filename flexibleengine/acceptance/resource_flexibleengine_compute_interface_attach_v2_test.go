package acceptance

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/compute/v2/extensions/attachinterfaces"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccComputeInterfaceAttach_Basic(t *testing.T) {
	var ai attachinterfaces.Interface
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "flexibleengine_compute_interface_attach_v2.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeInterfaceAttachDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInterfaceAttach_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInterfaceAttachExists(resourceName, &ai),
					testAccCheckComputeInterfaceAttachIP(&ai, "192.168.0.199"),
					resource.TestCheckResourceAttr(resourceName, "source_dest_check", "true"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_ids.0",
						"flexibleengine_networking_secgroup_v2.test", "id"),
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

func computeInterfaceAttachParseID(id string) (instanceID, portID string, err error) {
	idParts := strings.Split(id, "/")
	if len(idParts) < 2 {
		err = fmt.Errorf("unable to parse the resource ID, must be <instance_id>/<port_id> format")
		return
	}

	instanceID = idParts[0]
	portID = idParts[1]
	return
}

func testAccCheckComputeInterfaceAttachDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	computeClient, err := cfg.ComputeV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_compute_interface_attach_v2" {
			continue
		}

		instanceId, portId, err := computeInterfaceAttachParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = attachinterfaces.Get(computeClient, instanceId, portId).Extract()
		if err == nil {
			return fmt.Errorf("interface attachment still exists")
		}
	}

	return nil
}

func testAccCheckComputeInterfaceAttachExists(n string, ai *attachinterfaces.Interface) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		computeClient, err := cfg.ComputeV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating compute client: %s", err)
		}

		instanceId, portId, err := computeInterfaceAttachParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		found, err := attachinterfaces.Get(computeClient, instanceId, portId).Extract()
		if err != nil {
			return err
		}
		if found.PortID != portId {
			return fmt.Errorf("interface attachment not found")
		}

		*ai = *found

		return nil
	}
}

func testAccCheckComputeInterfaceAttachIP(
	ai *attachinterfaces.Interface, ip string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, i := range ai.FixedIPs {
			if i.IPAddress == ip {
				return nil
			}
		}
		return fmt.Errorf("requested ip (%s) does not exist on port", ip)
	}
}

func testAccComputeInterfaceAttach_basic(rName string) string {
	return fmt.Sprintf(`
data "flexibleengine_availability_zones" "test" {}

resource "flexibleengine_vpc_v1" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "test" {
  vpc_id     = flexibleengine_vpc_v1.test.id
  name       = "%[1]s"
  cidr       = cidrsubnet(flexibleengine_vpc_v1.test.cidr, 4, 0)
  gateway_ip = cidrhost(cidrsubnet(flexibleengine_vpc_v1.test.cidr, 4, 0), 1)
}

resource "flexibleengine_networking_secgroup_v2" "test" {
  name = "%[1]s"
}

resource "flexibleengine_apig_instance" "test" {
  name                  = "%[1]s"
  edition               = "BASIC"
  vpc_id                = flexibleengine_vpc_v1.test.id
  subnet_id             = flexibleengine_vpc_subnet_v1.test.id
  security_group_id     = flexibleengine_networking_secgroup_v2.test.id
  enterprise_project_id = "0"
  availability_zones    = try(slice(data.flexibleengine_availability_zones.test.names, 0, 1), null)
}

data "flexibleengine_compute_flavors_v2" "test" {
  availability_zone = data.flexibleengine_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core          = 2
  memory_size       = 4
}

data "flexibleengine_images_images" "test" {
  flavor_id  = data.flexibleengine_compute_flavors_v2.test.flavors[0]
  os         = "Ubuntu"
  visibility = "public"
}

resource "flexibleengine_compute_instance_v2" "test" {
  name              = "%[1]s"
  image_id          = data.flexibleengine_images_images.test.images[0].id
  flavor_id         = data.flexibleengine_compute_flavors_v2.test.flavors[0]
  availability_zone = data.flexibleengine_availability_zones.test.names[0]
  security_groups   = [flexibleengine_networking_secgroup_v2.test.name]

  network {
    uuid = flexibleengine_vpc_subnet_v1.test.id
  }
}

resource "flexibleengine_compute_interface_attach_v2" "test" {
  instance_id        = flexibleengine_compute_instance_v2.test.id
  network_id         = flexibleengine_vpc_subnet_v1.test.id
  fixed_ip           = cidrhost(cidrsubnet(flexibleengine_vpc_v1.test.cidr, 4, 0), 199)
  security_group_ids = [flexibleengine_networking_secgroup_v2.test.id]
}
`, rName)
}
