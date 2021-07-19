package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var keyAlias = fmt.Sprintf("key_alias_%s", acctest.RandString(5))

func TestAccKmsKeyV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKeyV1DataSource_key,
			},
			{
				Config: testAccKmsKeyV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKmsKeyV1DataSourceID("data.flexibleengine_kms_key_v1.key1"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_kms_key_v1.key1", "key_alias", keyAlias),
				),
			},
		},
	})
}

func testAccCheckKmsKeyV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Kms key data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Kms key data source ID not set")
		}

		return nil
	}
}

var testAccKmsKeyV1DataSource_key = fmt.Sprintf(`
resource "flexibleengine_kms_key_v1" "key1" {
  key_alias       = "%s"
  key_description = "test description"
  pending_days    = "7"
  is_enabled      = true
}`, keyAlias)

var testAccKmsKeyV1DataSource_basic = fmt.Sprintf(`
%s

data "flexibleengine_kms_key_v1" "key1" {
  key_alias       = flexibleengine_kms_key_v1.key1.key_alias
  key_id          = flexibleengine_kms_key_v1.key1.id
  key_description = "test description"
  key_state       = "2"
}
`, testAccKmsKeyV1DataSource_key)
