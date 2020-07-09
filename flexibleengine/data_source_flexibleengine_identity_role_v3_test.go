package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccIdentityV3RoleDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackIdentityV3RoleDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3DataSourceID("data.flexibleengine_identity_role_v3.role_1"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_identity_role_v3.role_1", "name", "rds_adm"),
				),
			},
		},
	})
}

func testAccCheckIdentityV3DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find role data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Role data source ID not set")
		}

		return nil
	}
}

const testAccOpenStackIdentityV3RoleDataSource_basic = `
data "flexibleengine_identity_role_v3" "role_1" {
    name = "rds_adm"
}
`
