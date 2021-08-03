package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVPCEPEndpointsDataSourceBasic(t *testing.T) {

	rName := fmt.Sprintf("acc-test-%s", acctest.RandString(4))
	endpointByNameResourceName := "data.flexibleengine_vpcep_endpoints.by_name"
	endpointByEndpointIdResourceName := "data.flexibleengine_vpcep_endpoints.by_endpoint_id"
	endpointByVpcIdResourceName := "data.flexibleengine_vpcep_endpoints.by_vpc_id"

	fmt.Sprintf(testAccVPCEPEndpointsDataSourceBasic(rName))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCEPEndpointsDataSourceBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(endpointByNameResourceName, "endpoints.0.vpc_id"),
					resource.TestCheckResourceAttrSet(endpointByNameResourceName, "endpoints.0.service_name"),
					resource.TestCheckResourceAttrSet(endpointByNameResourceName, "endpoints.0.service_id"),
					resource.TestCheckResourceAttrSet(endpointByEndpointIdResourceName, "endpoints.0.vpc_id"),
					resource.TestCheckResourceAttrSet(endpointByEndpointIdResourceName, "endpoints.0.service_name"),
					resource.TestCheckResourceAttrSet(endpointByEndpointIdResourceName, "endpoints.0.service_id"),
					resource.TestCheckResourceAttrSet(endpointByVpcIdResourceName, "endpoints.0.vpc_id"),
					resource.TestCheckResourceAttrSet(endpointByVpcIdResourceName, "endpoints.0.service_name"),
					resource.TestCheckResourceAttrSet(endpointByVpcIdResourceName, "endpoints.0.service_id"),
				),
			},
		},
	})
}

func testAccVPCEndpointDataSourcePrecondition(rName string) string {
	return fmt.Sprintf(`
resource flexibleengine_networking_network_v2 test {
  name           = "%[1]s"
  admin_state_up = "true"
}

resource flexibleengine_networking_subnet_v2 test {
  name            = "%[1]s"
  cidr            = "192.168.0.0/24"
  gateway_ip      = "192.168.0.1"
  network_id      = flexibleengine_networking_network_v2.test.id
}

resource flexibleengine_vpc_v1 test {
  name = "%[1]s"
  cidr = "192.168.0.0/24"
  tags = {
    owner = "terraform-test"
  }
}

resource flexibleengine_networking_router_interface_v2 test { 
  router_id = flexibleengine_vpc_v1.test.id 
  subnet_id = flexibleengine_networking_subnet_v2.test.id 
} 

resource flexibleengine_compute_instance_v2 test {
  name = "%[1]s"
  security_groups = ["default"]
  availability_zone = "%[2]s"

  network {
    uuid = flexibleengine_networking_network_v2.test.id
  }

  tags = {
    owner   = "terraform-test"
  }
  depends_on = [ flexibleengine_networking_router_interface_v2.test ]
}

resource "flexibleengine_vpcep_service" "test" {
  name        = "%[1]s"
  server_type = "VM"
  vpc_id      = flexibleengine_vpc_v1.test.id
  port_id     = flexibleengine_compute_instance_v2.test.network[0].port
  approval    = false

  port_mapping {
    service_port  = 22
    terminal_port = 22
  }
  tags = {
    owner = "terraform-test"
  }
}

resource flexibleengine_vpcep_endpoint test {
  service_id = flexibleengine_vpcep_service.test.id
  vpc_id = flexibleengine_vpc_v1.test.id
  network_id = flexibleengine_networking_network_v2.test.id
  tags = {
    owner = "terraform-test"
  }
}

`, rName, OS_AVAILABILITY_ZONE)
}

func testAccVPCEPEndpointsDataSourceBasic(rName string) string {
	return fmt.Sprintf(`
%s

data flexibleengine_vpcep_endpoints by_name {
   service_name = flexibleengine_vpcep_endpoint.test.service_name
}

data flexibleengine_vpcep_endpoints by_endpoint_id {
   endpoint_id = flexibleengine_vpcep_endpoint.test.id
}

data flexibleengine_vpcep_endpoints by_vpc_id {
   vpc_id = flexibleengine_vpcep_endpoint.test.vpc_id
}

`, testAccVPCEndpointDataSourcePrecondition(rName))
}
