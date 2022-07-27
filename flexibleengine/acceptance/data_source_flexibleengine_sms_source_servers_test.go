package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccSourceServers_basic(t *testing.T) {
	basicDataSource := "data.flexibleengine_sms_source_servers.all"
	byNameDataSource := "data.flexibleengine_sms_source_servers.byName"
	nonExistentDataSource := "data.flexibleengine_sms_source_servers.non-existent"
	basicDC := acceptance.InitDataSourceCheck(basicDataSource)
	name := OS_SMS_SOURCE_SERVER

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckSms(t)
		},
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSourceServers_basic(name),
				Check: resource.ComposeTestCheckFunc(
					basicDC.CheckResourceExists(),
					resource.TestCheckResourceAttr(basicDataSource, "servers.#", "1"),
					resource.TestCheckResourceAttr(basicDataSource, "servers.0.name", name),

					resource.TestCheckResourceAttr(byNameDataSource, "servers.#", "1"),
					resource.TestCheckResourceAttr(byNameDataSource, "servers.0.name", name),
					resource.TestCheckResourceAttrSet(byNameDataSource, "servers.0.ip"),
					resource.TestCheckResourceAttrSet(byNameDataSource, "servers.0.state"),

					resource.TestCheckResourceAttr(nonExistentDataSource, "id", "0"),
					resource.TestCheckResourceAttr(nonExistentDataSource, "servers.#", "0"),
				),
			},
		},
	})
}

func testAccSourceServers_basic(name string) string {
	return fmt.Sprintf(`
data "flexibleengine_sms_source_servers" "all" {
}

data "flexibleengine_sms_source_servers" "byName" {
  name = "%s"
}

data "flexibleengine_sms_source_servers" "non-existent" {
  name = "non-existent"
}
`, name)
}
