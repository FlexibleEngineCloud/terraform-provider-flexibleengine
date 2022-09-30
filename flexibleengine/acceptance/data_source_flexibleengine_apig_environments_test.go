package acceptance

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataEnvironments_basic(t *testing.T) {
	var (
		dataSourceName = "data.flexibleengine_apig_environments.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
		rName          = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataEnvironments_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "environments.#", regexp.MustCompile(`[1-9]\d*`)),
				),
			},
		},
	})
}

func testAccDataEnvironments_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_apig_environment" "test" {
  name        = "%s"
  instance_id = flexibleengine_apig_instance.test.id
  description = "Created by script"
}

data "flexibleengine_apig_environments" "test" {
  instance_id = flexibleengine_apig_instance.test.id
  name        = flexibleengine_apig_environment.test.name
}
`, testAccApigApplication_base(rName), rName)
}
