package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/huaweicloud/golangsdk/openstack/vpcep/v1/endpoints"
	"github.com/huaweicloud/golangsdk/openstack/vpcep/v1/services"
)

func TestAccVPCEndpointApproval(t *testing.T) {
	var service services.Service
	var endpoint endpoints.Endpoint

	rName := fmt.Sprintf("acc-test-%s", acctest.RandString(4))
	resourceName := "flexibleengine_vpcep_approval.approval"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVPCEPServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCEndpointApprovalBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCEPServiceExists("flexibleengine_vpcep_service.test", &service),
					testAccCheckVPCEndpointExists("flexibleengine_vpcep_endpoint.test", &endpoint),
					resource.TestCheckResourceAttrPtr(resourceName, "id", &service.ID),
					resource.TestCheckResourceAttrPtr(resourceName, "connections.0.endpoint_id", &endpoint.ID),
					resource.TestCheckResourceAttr(resourceName, "connections.0.status", "accepted"),
				),
			},
			{
				Config: testAccVPCEndpointApprovalUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr(resourceName, "connections.0.endpoint_id", &endpoint.ID),
					resource.TestCheckResourceAttr(resourceName, "connections.0.status", "rejected"),
				),
			},
		},
	})
}

func testAccVPCEndpointApprovalBasic(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_vpcep_service" "test" {
  name        = "%s"
  server_type = "VM"
  vpc_id      = "%s"
  port_id     = flexibleengine_compute_instance_v2.instance_1.network[0].port
  approval    = true

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
  lifecycle {
    ignore_changes = [enable_dns]
  }
}

resource "flexibleengine_vpcep_approval" "approval" {
  service_id = flexibleengine_vpcep_service.test.id
  endpoints  = [flexibleengine_vpcep_endpoint.test.id]
}
`, testAccVPCEndpointPrecondition(rName), rName, OS_VPC_ID, OS_VPC_ID, OS_NETWORK_ID)
}

func testAccVPCEndpointApprovalUpdate(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_vpcep_service" "test" {
  name        = "%s"
  server_type = "VM"
  vpc_id      = "%s"
  port_id     = flexibleengine_compute_instance_v2.instance_1.network[0].port
  approval    = true

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
  lifecycle {
    ignore_changes = [enable_dns]
  }
}

resource "flexibleengine_vpcep_approval" "approval" {
  service_id = flexibleengine_vpcep_service.test.id
  endpoints  = []
}
`, testAccVPCEndpointPrecondition(rName), rName, OS_VPC_ID, OS_VPC_ID, OS_NETWORK_ID)
}
