package acceptance

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/networking/v1/routetables"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getVpcRTBRouteResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	vpcClient, err := conf.NetworkingV1Client(OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("Error creating VPC client: %s", err)
	}

	parts := strings.SplitN(state.Primary.ID, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("the format of resource ID %s is invalid", state.Primary.ID)
	}

	routeTableID := parts[0]
	destination := parts[1]
	routeTable, err := routetables.Get(vpcClient, routeTableID).Extract()
	if err != nil {
		return nil, fmt.Errorf("Error retrieving VPC route table %s: %s", routeTableID, err)
	}

	var route *routetables.Route
	for _, item := range routeTable.Routes {
		if item.DestinationCIDR == destination {
			route = &item
			break
		}
	}
	if route == nil {
		return nil, fmt.Errorf("can not find the vpc route %s with %s", routeTableID, destination)
	}

	return route, nil
}

func TestAccVpcRTBRoute_basic(t *testing.T) {
	var route routetables.Route
	randName := acceptance.RandomAccResourceName()
	resourceName := "flexibleengine_vpc_route.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&route,
		getVpcRTBRouteResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcRTBRoute_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "type", "peering"),
					resource.TestCheckResourceAttr(resourceName, "description", "peering route"),
					resource.TestCheckResourceAttrSet(resourceName, "route_table_id"),
					resource.TestCheckResourceAttrSet(resourceName, "route_table_name"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "nexthop",
						"${flexibleengine_vpc_peering_connection_v2.test.id}"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "destination",
						"${flexibleengine_vpc_v1.test2.cidr}"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "vpc_id",
						"${flexibleengine_vpc_v1.test1.id}"),
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

func testAccVpcRTBRoute_basic(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_vpc_v1" "test1" {
  name = "%s_1"
  cidr = "172.16.0.0/16"
}

resource "flexibleengine_vpc_v1" "test2" {
  name = "%s_2"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_peering_connection_v2" "test" {
  name        = "%s"
  vpc_id      = flexibleengine_vpc_v1.test1.id
  peer_vpc_id = flexibleengine_vpc_v1.test2.id
}

resource "flexibleengine_vpc_route" "test" {
  vpc_id      = flexibleengine_vpc_v1.test1.id
  destination = flexibleengine_vpc_v1.test2.cidr
  type        = "peering"
  nexthop     = flexibleengine_vpc_peering_connection_v2.test.id
  description = "peering route"
}
`, rName, rName, rName)
}
