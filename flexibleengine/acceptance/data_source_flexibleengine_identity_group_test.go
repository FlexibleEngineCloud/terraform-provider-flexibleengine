package acceptance

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIdentityGroupDataSource_basic(t *testing.T) {
	resourceName := "data.flexibleengine_identity_group.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(resourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityGroupDataSource_by_name(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
		},
	})
}

func testAccIdentityGroupDataSource_by_name(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_identity_group_v3" "group_1" {
  name        = "%s"
  description = "A ACC test group"
}

data "flexibleengine_identity_group" "test" {
  name = flexibleengine_identity_group_v3.group_1.name
  
  depends_on = [
    flexibleengine_identity_group_v3.group_1
  ]
}
`, rName)
}
