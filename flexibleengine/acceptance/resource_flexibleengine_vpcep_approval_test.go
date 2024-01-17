package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/chnsz/golangsdk/openstack/vpcep/v1/endpoints"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccVPCEndpointApproval_Basic(t *testing.T) {
	var endpoint endpoints.Endpoint

	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "flexibleengine_vpcep_approval.approval"

	rc := acceptance.InitResourceCheck(
		"flexibleengine_vpcep_endpoint.test",
		&endpoint,
		getVpcepEndpointResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVPCEndpointApproval_Basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "id", "flexibleengine_vpcep_service.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "connections.0.endpoint_id",
						"flexibleengine_vpcep_endpoint.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "connections.0.status", "accepted"),
				),
			},
			{
				Config: testAccVPCEndpointApproval_Update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "connections.0.endpoint_id",
						"flexibleengine_vpcep_endpoint.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "connections.0.status", "rejected"),
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

func testAccVPCEndpointApproval_Base(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_vpcep_endpoint" "test" {
  service_id = flexibleengine_vpcep_service.test.id
  vpc_id     = flexibleengine_vpc_v1.test.id
  network_id = flexibleengine_vpc_subnet_v1.test.id
  enable_dns = true

  tags = {
    owner = "tf-acc"
  }
  lifecycle {
    ignore_changes = [enable_dns]
  }
}
`, testAccVPCEndpoint_Precondition(rName))
}

func testAccVPCEndpointApproval_Basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_vpcep_approval" "approval" {
  service_id = flexibleengine_vpcep_service.test.id
  endpoints  = [flexibleengine_vpcep_endpoint.test.id]
}
`, testAccVPCEndpointApproval_Base(rName))
}

func testAccVPCEndpointApproval_Update(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_vpcep_approval" "approval" {
  service_id = flexibleengine_vpcep_service.test.id
  endpoints  = []
}
`, testAccVPCEndpointApproval_Base(rName))
}
