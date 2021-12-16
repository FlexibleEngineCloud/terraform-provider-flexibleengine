package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var zoneName = fmt.Sprintf("acpttest%s.com.", acctest.RandString(5))

func TestAccDNSZoneV2DataSource_basic(t *testing.T) {
	resourceName := "data.flexibleengine_dns_zone_v2.z1"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSZoneV2DataSource_zone,
			},
			{
				Config: testAccDNSZoneV2DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSZoneV2DataSourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", zoneName),
					resource.TestCheckResourceAttr(resourceName, "zone_type", "public"),
					resource.TestCheckResourceAttr(resourceName, "ttl", "7200"),
				),
			},
		},
	})
}

func testAccCheckDNSZoneV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find DNS Zone data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("DNS Zone data source ID not set")
		}

		return nil
	}
}

var testAccDNSZoneV2DataSource_zone = fmt.Sprintf(`
resource "flexibleengine_dns_zone_v2" "z1" {
  name        = "%s"
  description = "dns data source test"
  email       = "terraform-dns-zone-v2-test-name@example.com"
  zone_type   = "public"
  ttl         = 7200
}`, zoneName)

var testAccDNSZoneV2DataSource_basic = fmt.Sprintf(`
%s

data "flexibleengine_dns_zone_v2" "z1" {
  name = flexibleengine_dns_zone_v2.z1.name
}
`, testAccDNSZoneV2DataSource_zone)
