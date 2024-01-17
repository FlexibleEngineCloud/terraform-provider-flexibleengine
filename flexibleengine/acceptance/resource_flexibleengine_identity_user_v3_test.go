package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/identity/v3.0/users"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getIdentityUserResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.IAMV3Client(OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}
	return users.Get(client, state.Primary.ID).Extract()
}

func TestAccIdentityUser_basic(t *testing.T) {
	var user users.User
	userName := acceptance.RandomAccResourceName()
	initPassword := acceptance.RandomPassword()
	newPassword := acceptance.RandomPassword()
	resourceName := "flexibleengine_identity_user_v3.user_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&user,
		getIdentityUserResourceFunc,
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
				Config: testAccIdentityUser_basic(userName, initPassword),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", userName),
					resource.TestCheckResourceAttr(resourceName, "description", "tested by terraform"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "pwd_reset", "true"),
					resource.TestCheckResourceAttr(resourceName, "email", "user_1@abc.com"),
					resource.TestCheckResourceAttr(resourceName, "password_strength", "Strong"),
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
			{
				Config: testAccIdentityUser_update(userName, newPassword),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", userName),
					resource.TestCheckResourceAttr(resourceName, "description", "updated by terraform"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "pwd_reset", "false"),
					resource.TestCheckResourceAttr(resourceName, "email", "user_1@abcd.com"),
				),
			},
			{
				Config: testAccIdentityUser_no_desc(userName, newPassword),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", userName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
				),
			},
		},
	})
}

func TestAccIdentityUser_external(t *testing.T) {
	var user users.User
	userName := acceptance.RandomAccResourceName()
	password := acceptance.RandomPassword()
	initXUserID := "123456789-abcdefg"
	newXUserID := "abcdefg-123456789"
	resourceName := "flexibleengine_identity_user_v3.user_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&user,
		getIdentityUserResourceFunc,
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
				Config: testAccIdentityUser_external(userName, password, initXUserID),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", userName),
					resource.TestCheckResourceAttr(resourceName, "description", "IAM user with external identity id"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "pwd_reset", "true"),
					resource.TestCheckResourceAttr(resourceName, "password_strength", "Strong"),
					resource.TestCheckResourceAttr(resourceName, "external_identity_id", initXUserID),
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
			{
				Config: testAccIdentityUser_external(userName, password, newXUserID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", userName),
					resource.TestCheckResourceAttr(resourceName, "external_identity_id", newXUserID),
				),
			},
		},
	})
}

func testAccIdentityUser_basic(name, password string) string {
	return fmt.Sprintf(`
resource "flexibleengine_identity_user_v3" "user_1" {
  name        = "%s"
  password    = "%s"
  enabled     = true
  email       = "user_1@abc.com"
  description = "tested by terraform"
}
`, name, password)
}

func testAccIdentityUser_update(name, password string) string {
	return fmt.Sprintf(`
resource "flexibleengine_identity_user_v3" "user_1" {
  name        = "%s"
  password    = "%s"
  pwd_reset   = false
  enabled     = false
  email       = "user_1@abcd.com"
  description = "updated by terraform"
}
`, name, password)
}

func testAccIdentityUser_no_desc(name, password string) string {
	return fmt.Sprintf(`
resource "flexibleengine_identity_user_v3" "user_1" {
  name      = "%s"
  password  = "%s"
  pwd_reset = false
  enabled   = false
  email     = "user_1@abcd.com"
}
`, name, password)
}

func testAccIdentityUser_external(name, password, xUserID string) string {
	return fmt.Sprintf(`
resource "flexibleengine_identity_user_v3" "user_1" {
  name                 = "%s"
  password             = "%s"
  description          = "IAM user with external identity id"
  external_identity_id = "%s"
}
`, name, password, xUserID)
}

