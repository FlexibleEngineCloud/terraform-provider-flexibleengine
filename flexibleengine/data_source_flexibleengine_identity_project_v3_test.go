package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccIdentityV3ProjectDataSource_basic(t *testing.T) {
	projectName := "eu-west-0_open_source"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProjectV3DataSource_basic(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProjectDataSourceID("data.flexibleengine_identity_project_v3.project_1"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_identity_project_v3.project_1", "name", projectName),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_identity_project_v3.project_1", "enabled", "true"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_identity_project_v3.project_1", "is_domain", "false"),
				),
			},
		},
	})
}

func testAccCheckIdentityV3ProjectDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find project data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Project data source ID not set")
		}

		return nil
	}
}

func testAccIdentityProjectV3DataSource_basic(name string) string {
	return fmt.Sprintf(`
data "flexibleengine_identity_project_v3" "project_1" {
  name = "%s"
}
	`, name)
}
