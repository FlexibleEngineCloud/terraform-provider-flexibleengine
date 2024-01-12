package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccVpcPeeringConnectionDataSource_basic(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	dataSourceName := "data.flexibleengine_vpc_peering_connection_v2.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcPeeringConnectionDataSource_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(dataSourceName, "description", "vpc1 peers to vpc2"),
				),
			},
		},
	})
}

func TestAccVpcPeeringConnectionDataSource_byVpcId(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	dataSourceName := "data.flexibleengine_vpc_peering_connection_v2.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcPeeringConnectionDataSource_byVpcId(randName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "status", "ACTIVE"),
				),
			},
		},
	})
}

func TestAccVpcPeeringConnectionDataSource_byPeerVpcId(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	dataSourceName := "data.flexibleengine_vpc_peering_connection_v2.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcPeeringConnectionDataSource_byPeerVpcId(randName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "status", "ACTIVE"),
				),
			},
		},
	})
}

func TestAccVpcPeeringConnectionDataSource_byVpcIds(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	dataSourceName := "data.flexibleengine_vpc_peering_connection_v2.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcPeeringConnectionDataSource_byVpcIds(randName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "status", "ACTIVE"),
				),
			},
		},
	})
}

func testAccVpcPeeringConnectionDataSource_base(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_vpc_v1" "vpc_1" {
  name = "%[1]s_1"
  cidr = "172.16.0.0/20"
}

resource "flexibleengine_vpc_v1" "vpc_2" {
  name = "%[1]s_2"
  cidr = "172.16.128.0/20"
}

resource "flexibleengine_vpc_peering_connection_v2" "test" {
  name        = "%[1]s"
  vpc_id      = flexibleengine_vpc_v1.vpc_1.id
  peer_vpc_id = flexibleengine_vpc_v1.vpc_2.id
  description = "vpc1 peers to vpc2"
}
`, rName)
}

func testAccVpcPeeringConnectionDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_vpc_peering_connection_v2" "test" {
  id = flexibleengine_vpc_peering_connection_v2.test.id
}
`, testAccVpcPeeringConnectionDataSource_base(rName))
}

func testAccVpcPeeringConnectionDataSource_byVpcId(rName string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_vpc_peering_connection_v2" "test" {
	vpc_id = flexibleengine_vpc_peering_connection_v2.test.vpc_id
}
`, testAccVpcPeeringConnectionDataSource_base(rName))
}

func testAccVpcPeeringConnectionDataSource_byPeerVpcId(rName string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_vpc_peering_connection_v2" "test" {
	peer_vpc_id = flexibleengine_vpc_peering_connection_v2.test.peer_vpc_id
}
`, testAccVpcPeeringConnectionDataSource_base(rName))
}

func testAccVpcPeeringConnectionDataSource_byVpcIds(rName string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_vpc_peering_connection_v2" "test" {
  vpc_id      = flexibleengine_vpc_peering_connection_v2.test.vpc_id
  peer_vpc_id = flexibleengine_vpc_peering_connection_v2.test.peer_vpc_id
}
`, testAccVpcPeeringConnectionDataSource_base(rName))
}
