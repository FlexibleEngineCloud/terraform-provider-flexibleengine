package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccVpcSubnetDataSource_ipv4Basic(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr, randGatewayIp := acceptance.RandomCidrAndGatewayIp()
	dataSourceName := "data.flexibleengine_vpc_subnet_v1.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetDataSource_ipv4Basic(randName, randCidr, randGatewayIp),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "gateway_ip", randGatewayIp),
					resource.TestCheckResourceAttr(dataSourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(dataSourceName, "dhcp_enable", "true"),
					resource.TestCheckResourceAttrSet(dataSourceName, "primary_dns"),
					resource.TestCheckResourceAttrSet(dataSourceName, "secondary_dns"),
					resource.TestCheckResourceAttrSet(dataSourceName, "ipv4_subnet_id"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "vpc_id",
						"${flexibleengine_vpc_v1.test.id}"),
				),
			},
		},
	})
}

func TestAccVpcSubnetDataSource_ipv4ByCidr(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr, randGatewayIp := acceptance.RandomCidrAndGatewayIp()
	dataSourceName := "data.flexibleengine_vpc_subnet_v1.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetDataSource_ipv4ByCidr(randName, randCidr, randGatewayIp),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "gateway_ip", randGatewayIp),
					resource.TestCheckResourceAttr(dataSourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(dataSourceName, "dhcp_enable", "true"),
					resource.TestCheckResourceAttrSet(dataSourceName, "primary_dns"),
					resource.TestCheckResourceAttrSet(dataSourceName, "secondary_dns"),
					resource.TestCheckResourceAttrSet(dataSourceName, "ipv4_subnet_id"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "vpc_id",
						"${flexibleengine_vpc_v1.test.id}"),
				),
			},
		},
	})
}

func TestAccVpcSubnetDataSource_ipv4ByName(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr, randGatewayIp := acceptance.RandomCidrAndGatewayIp()
	dataSourceName := "data.flexibleengine_vpc_subnet_v1.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetDataSource_ipv4ByName(randName, randCidr, randGatewayIp),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "gateway_ip", randGatewayIp),
					resource.TestCheckResourceAttr(dataSourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(dataSourceName, "dhcp_enable", "true"),
					resource.TestCheckResourceAttrSet(dataSourceName, "primary_dns"),
					resource.TestCheckResourceAttrSet(dataSourceName, "secondary_dns"),
					resource.TestCheckResourceAttrSet(dataSourceName, "ipv4_subnet_id"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "vpc_id",
						"${flexibleengine_vpc_v1.test.id}"),
				),
			},
		},
	})
}

func TestAccVpcSubnetDataSource_ipv4ByVpcId(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr, randGatewayIp := acceptance.RandomCidrAndGatewayIp()
	dataSourceName := "data.flexibleengine_vpc_subnet_v1.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetDataSource_ipv4ByVpcId(randName, randCidr, randGatewayIp),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "gateway_ip", randGatewayIp),
					resource.TestCheckResourceAttr(dataSourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(dataSourceName, "dhcp_enable", "true"),
					resource.TestCheckResourceAttrSet(dataSourceName, "primary_dns"),
					resource.TestCheckResourceAttrSet(dataSourceName, "secondary_dns"),
					resource.TestCheckResourceAttrSet(dataSourceName, "ipv4_subnet_id"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "vpc_id",
						"${flexibleengine_vpc_v1.test.id}"),
				),
			},
		},
	})
}

func TestAccVpcSubnetDataSource_ipv6Basic(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr, randGatewayIp := acceptance.RandomCidrAndGatewayIp()
	dataSourceName := "data.flexibleengine_vpc_subnet_v1.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetDataSource_ipv6Basic(randName, randCidr, randGatewayIp),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "gateway_ip", randGatewayIp),
					resource.TestCheckResourceAttr(dataSourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(dataSourceName, "dhcp_enable", "true"),
					resource.TestCheckResourceAttrSet(dataSourceName, "primary_dns"),
					resource.TestCheckResourceAttrSet(dataSourceName, "secondary_dns"),
					resource.TestCheckResourceAttrSet(dataSourceName, "ipv4_subnet_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "ipv6_subnet_id"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "vpc_id",
						"${flexibleengine_vpc_v1.test.id}"),
				),
			},
		},
	})
}

func testAccVpcSubnetDataSource_ipv4Base(rName, cidr, gatewayIp string) string {
	return fmt.Sprintf(`
resource "flexibleengine_vpc_v1" "test" {
  name = "%s"
  cidr = "%s"
}

resource "flexibleengine_vpc_subnet_v1" "test" {
  name       = "%s"
  vpc_id     = flexibleengine_vpc_v1.test.id
  cidr       = "%s"
  gateway_ip = "%s"
}`, rName, cidr, rName, cidr, gatewayIp)
}

func testAccVpcSubnetDataSource_ipv4Basic(rName, cidr, gatewayIp string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_vpc_subnet_v1" "test" {
  id = flexibleengine_vpc_subnet_v1.test.id
}
`, testAccVpcSubnetDataSource_ipv4Base(rName, cidr, gatewayIp))
}

func testAccVpcSubnetDataSource_ipv4ByCidr(rName, cidr, gatewayIp string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_vpc_subnet_v1" "test" {
  cidr = flexibleengine_vpc_subnet_v1.test.cidr
}
`, testAccVpcSubnetDataSource_ipv4Base(rName, cidr, gatewayIp))
}

func testAccVpcSubnetDataSource_ipv4ByName(rName, cidr, gatewayIp string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_vpc_subnet_v1" "test" {
  name = flexibleengine_vpc_subnet_v1.test.name
}
`, testAccVpcSubnetDataSource_ipv4Base(rName, cidr, gatewayIp))
}

func testAccVpcSubnetDataSource_ipv4ByVpcId(rName, cidr, gatewayIp string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_vpc_subnet_v1" "test" {
  vpc_id = flexibleengine_vpc_subnet_v1.test.vpc_id
}
`, testAccVpcSubnetDataSource_ipv4Base(rName, cidr, gatewayIp))
}

func testAccVpcSubnetDataSource_ipv6Base(rName, cidr, gatewayIp string) string {
	return fmt.Sprintf(`
resource "flexibleengine_vpc_v1" "test" {
  name = "%s"
  cidr = "%s"
}

resource "flexibleengine_vpc_subnet_v1" "test" {
  name        = "%s"
  cidr        = "%s"
  gateway_ip  = "%s"
  vpc_id      = flexibleengine_vpc_v1.test.id
  ipv6_enable = true
}`, rName, cidr, rName, cidr, gatewayIp)
}

func testAccVpcSubnetDataSource_ipv6Basic(rName, cidr, gatewayIp string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_vpc_subnet_v1" "test" {
  id = flexibleengine_vpc_subnet_v1.test.id
}
`, testAccVpcSubnetDataSource_ipv6Base(rName, cidr, gatewayIp))
}
