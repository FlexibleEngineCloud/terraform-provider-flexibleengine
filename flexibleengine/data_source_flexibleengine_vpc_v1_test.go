package flexibleengine

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"time"
)

func TestAccFlexibleEngineVpcV1DataSource_basic(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	rInt := rand.Intn(50)
	cidr := fmt.Sprintf("172.16.%d.0/24", rInt)
	name := fmt.Sprintf("terraform-testacc-vpc-data-source-%d", rInt)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceFlexibleEngineVpcV1Config(name, cidr),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceFlexibleEngineVpcV1Check("data.flexibleengine_vpc_v1.by_id", name, cidr),
					testAccDataSourceFlexibleEngineVpcV1Check("data.flexibleengine_vpc_v1.by_cidr", name, cidr),
					testAccDataSourceFlexibleEngineVpcV1Check("data.flexibleengine_vpc_v1.by_name", name, cidr),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_vpc_v1.by_id", "shared", "false"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_vpc_v1.by_id", "status", "OK"),
				),
			},
		},
	})
}

func testAccDataSourceFlexibleEngineVpcV1Check(n, name, cidr string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", n)
		}

		vpcRs, ok := s.RootModule().Resources["flexibleengine_vpc_v1.vpc_1"]
		if !ok {
			return fmt.Errorf("can't find flexibleengine_vpc_v1.vpc_1 in state")
		}

		attr := rs.Primary.Attributes

		if attr["id"] != vpcRs.Primary.Attributes["id"] {
			return fmt.Errorf(
				"id is %s; want %s",
				attr["id"],
				vpcRs.Primary.Attributes["id"],
			)
		}

		if attr["cidr"] != cidr {
			return fmt.Errorf("bad vpc cidr %s, expected: %s", attr["cidr"], cidr)
		}
		if attr["name"] != name {
			return fmt.Errorf("bad vpc name %s", attr["name"])
		}

		return nil
	}
}

func testAccDataSourceFlexibleEngineVpcV1Config(name, cidr string) string {
	return fmt.Sprintf(`
resource "flexibleengine_vpc_v1" "vpc_1" {
	name = "%s"
	cidr= "%s"
}

data "flexibleengine_vpc_v1" "by_id" {
  id = "${flexibleengine_vpc_v1.vpc_1.id}"
}

data "flexibleengine_vpc_v1" "by_cidr" {
  cidr = "${flexibleengine_vpc_v1.vpc_1.cidr}"
}

data "flexibleengine_vpc_v1" "by_name" {
	name = "${flexibleengine_vpc_v1.vpc_1.name}"
}
`, name, cidr)
}
