package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccBMSV2KeyPairDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccBmsKeyPairPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBMSV2KeyPairDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBMSV2KeyPairDataSourceID("data.flexibleengine_compute_bms_keypairs_v2.keypair"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_compute_bms_keypairs_v2.keypair", "name", OS_KEYPAIR_NAME),
				),
			},
		},
	})
}

func testAccCheckBMSV2KeyPairDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find keypair data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Keypair data source ID not set")
		}

		return nil
	}
}

var testAccBMSV2KeyPairDataSource_basic = fmt.Sprintf(`
data "flexibleengine_compute_bms_keypairs_v2" "keypair" {
  name = "%s"
}
`, OS_KEYPAIR_NAME)
