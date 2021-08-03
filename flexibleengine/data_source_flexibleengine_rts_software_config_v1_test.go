package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccRtsConfigV1DataSource_basic(t *testing.T) {
	var rtsName = fmt.Sprintf("terra-test-%s", acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRtsSoftwareConfigV1DataSource_basic(rtsName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRtsConfigV1DataSourceID("data.flexibleengine_rts_software_config_v1.configs"),
					resource.TestCheckResourceAttr("data.flexibleengine_rts_software_config_v1.configs", "name", rtsName),
					resource.TestCheckResourceAttr("data.flexibleengine_rts_software_config_v1.configs", "group", "script"),
				),
			},
		},
	})
}

func testAccCheckRtsConfigV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find software config data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("RTS software config data source ID not set ")
		}

		return nil
	}
}

func testAccRtsSoftwareConfigV1DataSource_basic(rtsName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_rts_software_config_v1" "config_1" {
  name = "%s"
  output_values = [{
    type = "String"
    name = "result"
    error_output = "false"
    description = "value1"
  }]
  input_values = [{
    default = "0"
    type = "String"
    name = "foo"
    description = "value2"
  }]
  group = "script"
}

data "flexibleengine_rts_software_config_v1" "configs" {
  id = "${flexibleengine_rts_software_config_v1.config_1.id}"
}
`, rtsName)
}
