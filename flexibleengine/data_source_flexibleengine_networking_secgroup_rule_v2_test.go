package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccFlexibleEngineNetworkingSecGroupRuleV2DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFlexibleEngineNetworkingSecGroupRuleV2DataSource_group,
			},
			{
				Config: testAccFlexibleEngineNetworkingSecGroupRuleV2DataSource_rule,
			},
			{
				Config: testAccFlexibleEngineNetworkingSecGroupRuleV2DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupRuleV2DataSourceID("data.flexibleengine_networking_secgroup_rule_v2.test_provider_secgroup_rule_1"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_networking_secgroup_rule_v2.test_provider_secgroup_rule_1", "direction", "ingress"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_networking_secgroup_rule_v2.test_provider_secgroup_rule_1", "ethertype", "IPv4"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_networking_secgroup_rule_v2.test_provider_secgroup_rule_1", "port_range_max", "22"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_networking_secgroup_rule_v2.test_provider_secgroup_rule_1", "port_range_min", "22"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_networking_secgroup_rule_v2.test_provider_secgroup_rule_1", "protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_networking_secgroup_rule_v2.test_provider_secgroup_rule_1", "remote_ip_prefix", "192.168.0.1/32"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_networking_secgroup_rule_v2.test_provider_secgroup_rule_1", "description", "allow SSH"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_networking_secgroup_rule_v2.test_provider_secgroup_rule_1", "protocol", "tcp"),
				),
			},
		},
	})
}

func testAccCheckNetworkingSecGroupRuleV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find security group rule data source: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Security group rule data source ID not set")
		}
		return nil
	}
}

const testAccFlexibleEngineNetworkingSecGroupRuleV2DataSource_group = `
resource "flexibleengine_networking_secgroup_v2" "test_provider_secgroup_1" {
	name        = "flexibleengine_acctest_secgroup"
	description = "My neutron security group for flexibleengine acctest"
}
`

var testAccFlexibleEngineNetworkingSecGroupRuleV2DataSource_rule = fmt.Sprintf(` 
%s

resource "flexibleengine_networking_secgroup_rule_v2" "test_provider_secgroup_rule_1" {
	direction         = "ingress"
	ethertype         = "IPv4"
	port_range_max    = 22
	port_range_min    = 22
	protocol          = "tcp"
	remote_ip_prefix  = "192.168.0.1/32"
	security_group_id = flexibleengine_networking_secgroup_v2.test_provider_secgroup_1.id
	description       = "allow SSH"
  }`, testAccFlexibleEngineNetworkingSecGroupRuleV2DataSource_group)

var testAccFlexibleEngineNetworkingSecGroupRuleV2DataSource_basic = fmt.Sprintf(`
%s
  data "flexibleengine_networking_secgroup_rule_v2" "test_provider_secgroup_rule_1" {
	security_group_id = flexibleengine_networking_secgroup_v2.test_provider_secgroup_1.id
	direction = "ingress"
}`, testAccFlexibleEngineNetworkingSecGroupRuleV2DataSource_rule)
