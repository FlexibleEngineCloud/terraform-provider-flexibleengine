package acceptance

import (
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccEnterpriseProjectDataSource_basic(t *testing.T) {
	dataSourceName := "data.flexibleengine_enterprise_project.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEnterpriseProjectDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", "default"),
					resource.TestCheckResourceAttr(dataSourceName, "id", "0"),
				),
			},
		},
	})
}

const testAccEnterpriseProjectDataSource_basic = `
data "flexibleengine_enterprise_project" "test" {
  name = "default"
}
`
