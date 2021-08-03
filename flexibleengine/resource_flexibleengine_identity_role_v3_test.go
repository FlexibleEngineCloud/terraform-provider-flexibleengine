package flexibleengine

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/golangsdk/openstack/identity/v3.0/policies"
)

func TestAccIdentityRole_basic(t *testing.T) {
	var role policies.Role
	var roleName = fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))
	var roleNameUpdate = roleName + "update"
	resourceName := "flexibleengine_identity_role_v3.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIdentityRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityRole_basic(roleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityRoleExists(resourceName, &role),
					resource.TestCheckResourceAttrPtr(
						resourceName, "name", &role.Name),
					resource.TestCheckResourceAttrPtr(
						resourceName, "description", &role.Description),
				),
			},
			{
				Config: testAccIdentityRole_update(roleNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityRoleExists(resourceName, &role),
					resource.TestCheckResourceAttrPtr(
						resourceName, "name", &role.Name),
					resource.TestCheckResourceAttrPtr(
						resourceName, "description", &role.Description),
				),
			},
		},
	})
}

func testAccCheckIdentityRoleDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	identityClient, err := config.identityV3Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine identity client: %s", err)
	}

	identityClient.Endpoint = strings.Replace(identityClient.Endpoint, "v3", "v3.0", 1)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_identity_role_v3" {
			continue
		}

		_, err := policies.Get(identityClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Role still exists")
		}
	}

	return nil
}

func testAccCheckIdentityRoleExists(n string, role *policies.Role) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		identityClient, err := config.identityV3Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine identity client: %s", err)
		}

		identityClient.Endpoint = strings.Replace(identityClient.Endpoint, "v3", "v3.0", 1)

		found, err := policies.Get(identityClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Role not found")
		}

		*role = *found

		return nil
	}
}

func testAccIdentityRole_basic(roleName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_identity_role_v3" "test" {
  name = "%s"
  description = "created by terraform"
  type = "AX"
    policy = <<EOF
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
  `, roleName)
}

func testAccIdentityRole_update(roleName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_identity_role_v3" "test" {
  name = "%s"
  description = "created by terraform"
  type = "AX"
    policy = <<EOF
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
      ],
      "Condition": {
        "StringStartWith": {
          "g:ProjectName": [
            "eu-west-0"
          ]
        }
      }
    }
  ]
}
EOF
}
  `, roleName)
}
