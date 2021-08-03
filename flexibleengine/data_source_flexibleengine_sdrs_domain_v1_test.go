package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSdrsDomainV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckSdrs(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSdrsDomainV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSdrsDomainV1DataSourceID("data.flexibleengine_sdrs_domain_v1.domain_1"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_sdrs_domain_v1.domain_1", "name", "SDRS_HypeDomain01"),
				),
			},
		},
	})
}

func testAccCheckSdrsDomainV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find SDRS domain data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("SDRS domain data source ID not set")
		}

		return nil
	}
}

const testAccSdrsDomainV1DataSource_basic = `
data "flexibleengine_sdrs_domain_v1" "domain_1" {
	name = "SDRS_HypeDomain01"
}
`
