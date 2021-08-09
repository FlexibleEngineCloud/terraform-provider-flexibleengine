package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/golangsdk/openstack/identity/v3.0/users"
)

func TestAccIdentityV3User_basic(t *testing.T) {
	var user users.User
	var userName = fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))
	resourceName := "flexibleengine_identity_user_v3.user_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckIdentityV3UserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3User_basic(userName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3UserExists(resourceName, &user),
					resource.TestCheckResourceAttr(resourceName, "name", userName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "email", "foo123@orange-business.com"),
					resource.TestCheckResourceAttr(resourceName, "description", "created by terraform"),
					resource.TestCheckResourceAttrSet(resourceName, "password_strength"),
				),
			},
			{
				Config: testAccIdentityV3User_update(userName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3UserExists(resourceName, &user),
					resource.TestCheckResourceAttr(resourceName, "name", userName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "email", "bar123@orange-business.com"),
					resource.TestCheckResourceAttr(resourceName, "description", "updated by terraform"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
				},
			},
		},
	})
}

func testAccCheckIdentityV3UserDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	iamClient, err := config.IAMV3Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine identity client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_identity_user_v3" {
			continue
		}

		_, err := users.Get(iamClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("User still exists")
		}
	}

	return nil
}

func testAccCheckIdentityV3UserExists(n string, user *users.User) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		iamClient, err := config.IAMV3Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine identity client: %s", err)
		}

		found, err := users.Get(iamClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("User not found")
		}

		*user = *found

		return nil
	}
}

func testAccIdentityV3User_basic(userName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_identity_user_v3" "user_1" {
  name        = "%s"
  password    = "password123@!"
  enabled     = true
  email       = "foo123@orange-business.com"
  description = "created by terraform"
}
`, userName)
}

func testAccIdentityV3User_update(userName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_identity_user_v3" "user_1" {
  name        = "%s"
  enabled     = false
  password    = "password123@!"
  email       = "bar123@orange-business.com"
  description = "updated by terraform"
}
`, userName)
}
