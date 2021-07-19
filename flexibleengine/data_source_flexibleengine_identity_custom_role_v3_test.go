package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccFlexibleEngineIdentityCustomRoleV3DataSource_basic(t *testing.T) {
	var rName = fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFlexibleEngineIdentityCustomRoleDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityCustomDataSourceID("data.flexibleengine_identity_custom_role_v3.role_1"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_identity_custom_role_v3.role_1", "name", rName),
				),
			},
		},
	})
}

func testAccCheckIdentityCustomDataSourceID(n string) resource.TestCheckFunc {
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

func testAccFlexibleEngineIdentityCustomRoleDataSource_basic(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_identity_role_v3" test {
  name        = "%s"
  description = "created by terraform"
  type        = "AX"
  policy      = <<EOF
{
  "Version": "1.1",
  "Statement": [
    {
      "Action": [
        "obs:bucket:GetBucketAcl"
      ],
      "Effect": "Allow",
      "Resource": [
        "obs:*:*:bucket:*"
      ]
    }
  ]
}
EOF
}

data "flexibleengine_identity_custom_role_v3" "role_1" {
  name = flexibleengine_identity_role_v3.test.name
}
`, rName)
}
