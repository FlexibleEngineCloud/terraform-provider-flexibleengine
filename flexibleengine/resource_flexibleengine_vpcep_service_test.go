package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/golangsdk/openstack/vpcep/v1/services"
)

func TestAccVPCEPServiceBasic(t *testing.T) {
	var service services.Service

	rName := fmt.Sprintf("acc-test-%s", acctest.RandString(4))
	resourceName := "flexibleengine_vpcep_service.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVPCEPServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCEPServiceBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCEPServiceExists(resourceName, &service),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "available"),
					resource.TestCheckResourceAttr(resourceName, "approval", "false"),
					resource.TestCheckResourceAttr(resourceName, "server_type", "VM"),
					resource.TestCheckResourceAttr(resourceName, "service_type", "interface"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "tf-acc"),
					resource.TestCheckResourceAttr(resourceName, "port_mapping.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(resourceName, "port_mapping.0.service_port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "port_mapping.0.terminal_port", "80"),
				),
			},
			{
				Config: testAccVPCEPServiceUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "tf-"+rName),
					resource.TestCheckResourceAttr(resourceName, "status", "available"),
					resource.TestCheckResourceAttr(resourceName, "approval", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "tf-acc-update"),
					resource.TestCheckResourceAttr(resourceName, "port_mapping.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(resourceName, "port_mapping.0.service_port", "8088"),
					resource.TestCheckResourceAttr(resourceName, "port_mapping.0.terminal_port", "80"),
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

func TestAccVPCEPServicePermission(t *testing.T) {
	var service services.Service

	rName := fmt.Sprintf("acc-test-%s", acctest.RandString(4))
	resourceName := "flexibleengine_vpcep_service.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVPCEPServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCEPServicePermission(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCEPServiceExists(resourceName, &service),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "available"),
					resource.TestCheckResourceAttr(resourceName, "permissions.#", "2"),
				),
			},
			{
				Config: testAccVPCEPServicePermissionUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "available"),
					resource.TestCheckResourceAttr(resourceName, "permissions.#", "1"),
				),
			},
		},
	})
}

func testAccCheckVPCEPServiceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	vpcepClient, err := config.vpcepV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating VPC endpoint client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_vpcep_service" {
			continue
		}

		_, err := services.Get(vpcepClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("VPC endpoint service still exists")
		}
	}

	return nil
}

func testAccCheckVPCEPServiceExists(n string, service *services.Service) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		vpcepClient, err := config.vpcepV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating VPC endpoint client: %s", err)
		}

		found, err := services.Get(vpcepClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("VPC endpoint service not found")
		}

		*service = *found

		return nil
	}
}

func testAccVPCEndpointPrecondition(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_compute_instance_v2" "instance_1" {
  name = "%s"
  security_groups = ["default"]
  availability_zone = "%s"

  network {
    uuid = "%s"
  }
  tags = {
    owner   = "terraform"
    service = "vpc-endpoint"
  }
}
`, rName, OS_AVAILABILITY_ZONE, OS_NETWORK_ID)
}

func testAccVPCEPServiceBasic(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_vpcep_service" "test" {
  name        = "%s"
  server_type = "VM"
  vpc_id      = "%s"
  port_id     = flexibleengine_compute_instance_v2.instance_1.network[0].port
  approval    = false

  port_mapping {
    service_port  = 8080
    terminal_port = 80
  }
  tags = {
    owner = "tf-acc"
  }
}
`, testAccVPCEndpointPrecondition(rName), rName, OS_VPC_ID)
}

func testAccVPCEPServiceUpdate(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_vpcep_service" "test" {
  name        = "tf-%s"
  server_type = "VM"
  vpc_id      = "%s"
  port_id     = flexibleengine_compute_instance_v2.instance_1.network[0].port
  approval    = true

  port_mapping {
    service_port  = 8088
    terminal_port = 80
  }
  tags = {
    owner = "tf-acc-update"
  }
}
`, testAccVPCEndpointPrecondition(rName), rName, OS_VPC_ID)
}

func testAccVPCEPServicePermission(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_vpcep_service" "test" {
  name        = "%s"
  server_type = "VM"
  vpc_id      = "%s"
  port_id     = flexibleengine_compute_instance_v2.instance_1.network[0].port
  approval    = false
  permissions = ["iam:domain::1234", "iam:domain::5678"]

  port_mapping {
    service_port  = 8080
    terminal_port = 80
  }
}
`, testAccVPCEndpointPrecondition(rName), rName, OS_VPC_ID)
}

func testAccVPCEPServicePermissionUpdate(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_vpcep_service" "test" {
  name        = "%s"
  server_type = "VM"
  vpc_id      = "%s"
  port_id     = flexibleengine_compute_instance_v2.instance_1.network[0].port
  approval    = false
  permissions = ["iam:domain::abcd"]

  port_mapping {
    service_port  = 8080
    terminal_port = 80
  }
}
`, testAccVPCEndpointPrecondition(rName), rName, OS_VPC_ID)
}
