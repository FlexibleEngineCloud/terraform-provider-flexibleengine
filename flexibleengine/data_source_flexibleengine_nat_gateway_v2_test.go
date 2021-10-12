package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNatGatewayDataSource_basic(t *testing.T) {
	natgateway := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNatGatewayV2DataSource_basic(natgateway),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatGatewayV2DataSourceID("data.flexibleengine_nat_gateway_v2.nat_by_name"),
					testAccCheckNatGatewayV2DataSourceID("data.flexibleengine_nat_gateway_v2.nat_by_id"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_nat_gateway_v2.nat_by_name", "name", natgateway),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_nat_gateway_v2.nat_by_id", "name", natgateway),
				),
			},
		},
	})
}

func testAccCheckNatGatewayV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find natgateway data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("NatGateway data source ID not set")
		}

		return nil
	}
}

func testAccNatGatewayV2DataSource_basic(name string) string {
	return fmt.Sprintf(`
resource "flexibleengine_vpc_v1" "vpc_1" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "subnet_1" {
  name       = "%s"
  cidr       = "192.168.199.0/24"
  gateway_ip = "192.168.199.1"
  vpc_id     = flexibleengine_vpc_v1.vpc_1.id
}

resource "flexibleengine_nat_gateway_v2" "nat_1" {
  name        = "%s"
  description = "test for terraform"
  spec        = "1"
  vpc_id      = flexibleengine_vpc_v1.vpc_1.id
  subnet_id   = flexibleengine_vpc_subnet_v1.subnet_1.id
}

data "flexibleengine_nat_gateway_v2" "nat_by_name" {
  name = flexibleengine_nat_gateway_v2.nat_1.name
}

data "flexibleengine_nat_gateway_v2" "nat_by_id" {
  id = flexibleengine_nat_gateway_v2.nat_1.id
}
`, name, name, name)
}
