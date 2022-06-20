package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/swr/v2/namespaces"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getResourcePermissions(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	swrClient, err := conf.SwrV2Client(OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("Error creating SWR client: %s", err)
	}

	return namespaces.GetAccess(swrClient, state.Primary.ID).Extract()
}

func TestAccSwrOrganizationPermissions_basic(t *testing.T) {
	var permissions namespaces.Access
	organizationName := acceptance.RandomAccResourceName()
	userName1 := acceptance.RandomAccResourceName()
	userName2 := acceptance.RandomAccResourceName()
	resourceName := "flexibleengine_swr_organization_users.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&permissions,
		getResourcePermissions,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccswrOrganizationPermissions_basic(organizationName, userName1, userName2),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "users.0.user_name", userName1),
					resource.TestCheckResourceAttr(resourceName, "users.0.permission", "Read"),
					resource.TestCheckResourceAttrPair(resourceName, "organization", "flexibleengine_swr_organization.test", "name"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccswrOrganizationPermissions_update(organizationName, userName1, userName2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "users.0.user_name", userName1),
					resource.TestCheckResourceAttr(resourceName, "users.0.permission", "Write"),
					resource.TestCheckResourceAttr(resourceName, "users.1.user_name", userName2),
					resource.TestCheckResourceAttr(resourceName, "users.1.permission", "Read"),
					resource.TestCheckResourceAttrPair(resourceName, "organization", "flexibleengine_swr_organization.test", "name"),
				),
			},
		},
	})
}

func testAccswrOrganizationPermissions_basic(organizationName, userName1, userName2 string) string {
	return fmt.Sprintf(`
resource "flexibleengine_swr_organization" "test" {
  name = "%s"
}

resource "flexibleengine_identity_user_v3" "user_1" {
  name     = "%s"
  enabled  = true
  password = "password12345!"
}

resource "flexibleengine_swr_organization_users" "test" {
  organization = flexibleengine_swr_organization.test.name

  users {
    user_name  = flexibleengine_identity_user_v3.user_1.name
    user_id    = flexibleengine_identity_user_v3.user_1.id
    permission = "Read"
  }
}
`, organizationName, userName1)
}

func testAccswrOrganizationPermissions_update(organizationName, userName1, userName2 string) string {
	return fmt.Sprintf(`
resource "flexibleengine_swr_organization" "test" {
  name = "%s"
}

resource "flexibleengine_identity_user_v3" "user_1" {
  name     = "%s"
  enabled  = true
  password = "password12345!"
}

resource "flexibleengine_identity_user_v3" "user_2" {
  name     = "%s"
  enabled  = true
  password = "password12345!"
}

resource "flexibleengine_swr_organization_users" "test" {
  organization = flexibleengine_swr_organization.test.name

  users {
    user_name  = flexibleengine_identity_user_v3.user_1.name
    user_id    = flexibleengine_identity_user_v3.user_1.id
    permission = "Write"
  }

  users {
    user_name  = flexibleengine_identity_user_v3.user_2.name
    user_id    = flexibleengine_identity_user_v3.user_2.id
    permission = "Read"
  }
}
`, organizationName, userName1, userName2)
}
