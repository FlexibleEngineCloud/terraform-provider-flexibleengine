package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/snatrules"
)

func TestAccNatSnatRule_basic(t *testing.T) {
	randSuffix := acctest.RandString(5)
	resourceName := "flexibleengine_nat_snat_rule_v2.snat_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNatV2SnatRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNatV2SnatRule_basic(randSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatV2GatewayExists("flexibleengine_nat_gateway_v2.nat_1"),
					testAccCheckNatV2SnatRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
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

func testAccCheckNatV2SnatRuleDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	natClient, err := config.natV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine nat client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_nat_snat_rule_v2" {
			continue
		}

		_, err := snatrules.Get(natClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Snat rule still exists")
		}
	}

	return nil
}

func testAccCheckNatV2SnatRuleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		natClient, err := config.natV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine nat client: %s", err)
		}

		found, err := snatrules.Get(natClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Snat rule not found")
		}

		return nil
	}
}

func testAccNatV2SnatRule_basic(suffix string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_networking_floatingip_v2" "fip_1" {
}

resource "flexibleengine_nat_gateway_v2" "nat_1" {
  name        = "natgw-test-%s"
  description = "test for terraform"
  spec        = "1"
  vpc_id      = flexibleengine_vpc_v1.vpc_1.id
  subnet_id   = flexibleengine_vpc_subnet_v1.subnet_1.id
}

resource "flexibleengine_nat_snat_rule_v2" "snat_1" {
  nat_gateway_id = flexibleengine_nat_gateway_v2.nat_1.id
  subnet_id      = flexibleengine_vpc_subnet_v1.subnet_1.id
  floating_ip_id = flexibleengine_networking_floatingip_v2.fip_1.id
}
`, testAccNatPreCondition(suffix), suffix)
}
