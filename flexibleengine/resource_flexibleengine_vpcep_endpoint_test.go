package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/golangsdk/openstack/vpcep/v1/endpoints"
)

func TestAccVPCEndpointBasic(t *testing.T) {
	var endpoint endpoints.Endpoint

	rName := fmt.Sprintf("acc-test-%s", acctest.RandString(4))
	resourceName := "flexibleengine_vpcep_endpoint.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVPCEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCEndpointBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCEndpointExists(resourceName, &endpoint),
					resource.TestCheckResourceAttr(resourceName, "status", "accepted"),
					resource.TestCheckResourceAttr(resourceName, "enable_dns", "true"),
					resource.TestCheckResourceAttr(resourceName, "service_type", "interface"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "tf-acc"),
					resource.TestCheckResourceAttrSet(resourceName, "service_name"),
					resource.TestCheckResourceAttrSet(resourceName, "private_domain_name"),
				),
			},
			{
				Config: testAccVPCEndpointUpdate(rName),
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

func TestAccVPCEndpointPublic(t *testing.T) {
	var endpoint endpoints.Endpoint
	resourceName := "flexibleengine_vpcep_endpoint.myendpoint"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVPCEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCEndpointPublic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCEndpointExists(resourceName, &endpoint),
					resource.TestCheckResourceAttr(resourceName, "status", "accepted"),
					resource.TestCheckResourceAttr(resourceName, "enable_dns", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_whitelist", "true"),
					resource.TestCheckResourceAttr(resourceName, "service_type", "interface"),
					resource.TestCheckResourceAttr(resourceName, "whitelist.#", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "service_name"),
					resource.TestCheckResourceAttrSet(resourceName, "private_domain_name"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
				),
			},
		},
	})
}

func testAccCheckVPCEndpointDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	vpcepClient, err := config.vpcepV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating VPC endpoint client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_vpcep_endpoint" {
			continue
		}

		_, err := endpoints.Get(vpcepClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("VPC endpoint still exists")
		}
	}

	return nil
}

func testAccCheckVPCEndpointExists(n string, endpoint *endpoints.Endpoint) resource.TestCheckFunc {
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

		found, err := endpoints.Get(vpcepClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("VPC endpoint not found")
		}

		*endpoint = *found

		return nil
	}
}

func testAccVPCEndpointBasic(rName string) string {
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

resource "flexibleengine_vpcep_endpoint" "test" {
  service_id  = flexibleengine_vpcep_service.test.id
  vpc_id      = "%s"
  network_id  = "%s"
  enable_dns  = true

  tags = {
    owner = "tf-acc"
  }
}
`, testAccVPCEndpointPrecondition(rName), rName, OS_VPC_ID, OS_VPC_ID, OS_NETWORK_ID)
}

func testAccVPCEndpointUpdate(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_vpcep_service" "test" {
  name        = "tf-%s"
  server_type = "VM"
  vpc_id      = "%s"
  port_id     = flexibleengine_compute_instance_v2.instance_1.network[0].port
  approval    = false

  port_mapping {
    service_port  = 8088
    terminal_port = 80
  }
  tags = {
    owner = "tf-acc"
  }
}

resource "flexibleengine_vpcep_endpoint" "test" {
  service_id  = flexibleengine_vpcep_service.test.id
  vpc_id      = "%s"
  network_id  = "%s"
  enable_dns  = true

  tags = {
    owner = "tf-acc-update"
    foo   = "bar"
  }
}
`, testAccVPCEndpointPrecondition(rName), rName, OS_VPC_ID, OS_VPC_ID, OS_NETWORK_ID)
}

var testAccVPCEndpointPublic string = fmt.Sprintf(`
data "flexibleengine_vpcep_public_services" "cloud_service" {
  service_name = "dns"
}

resource "flexibleengine_vpcep_endpoint" "myendpoint" {
  service_id       = data.flexibleengine_vpcep_public_services.cloud_service.services[0].id
  vpc_id           = "%s"
  network_id       = "%s"
  enable_dns       = true
  enable_whitelist = true
  whitelist        = ["192.168.0.0/24", "10.10.10.10"]
}
`, OS_VPC_ID, OS_NETWORK_ID)
