package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkingV2PortDataSource_basic(t *testing.T) {
	resourceName := "data.flexibleengine_networking_port.gw_port"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "all_fixed_ips.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "mac_address"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
		},
	})
}

func testAccNetworkingV2PortDataSource_basic() string {
	return `
data "flexibleengine_vpc_subnet_v1" "mynet" {
  name = "subnet-default"
}

data "flexibleengine_networking_port" "gw_port" {
  network_id = data.flexibleengine_vpc_subnet_v1.mynet.id
  fixed_ip   = data.flexibleengine_vpc_subnet_v1.mynet.gateway_ip
}
`
}
