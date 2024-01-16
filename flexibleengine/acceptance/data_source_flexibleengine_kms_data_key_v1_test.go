package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccKmsDataKeyV1DataSource_basic(t *testing.T) {
	keyAlias := acceptance.RandomAccResourceName()
	datasourceName := "data.flexibleengine_kms_data_key_v1.test"
	dc := acceptance.InitDataSourceCheck(datasourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckKms(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsDataKeyV1DataSource_basic(keyAlias),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(
						datasourceName, "plain_text"),
					resource.TestCheckResourceAttrSet(
						datasourceName, "cipher_text"),
				),
			},
		},
	})
}

func testAccKmsDataKeyV1DataSource_basic(keyAlias string) string {
	return fmt.Sprintf(`
resource "flexibleengine_kms_key_v1" "test" {
  key_alias    = "%s"
  pending_days = "7"
}

data "flexibleengine_kms_data_key_v1" "test" {
  key_id         = flexibleengine_kms_key_v1.test.id
  datakey_length = "512"
}
`, keyAlias)
}
