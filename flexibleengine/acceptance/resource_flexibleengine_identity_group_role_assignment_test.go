package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/identity/v3.0/eps_permissions"
	"github.com/chnsz/golangsdk/openstack/identity/v3/roles"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/iam"
)

func getIdentityGroupRoleAssignmentResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	identityClient, err := c.IdentityV3Client(OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM v3 client: %s", err)
	}

	iamClient, err := c.IAMV3Client(OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM v3.0 client: %s", err)
	}

	groupID := state.Primary.Attributes["group_id"]
	roleID := state.Primary.Attributes["role_id"]
	domainID := state.Primary.Attributes["domain_id"]
	projectID := state.Primary.Attributes["project_id"]
	enterpriseProjectID := state.Primary.Attributes["enterprise_project_id"]

	if domainID != "" {
		return iam.GetGroupRoleAssignmentWithDomainID(identityClient, groupID, roleID, domainID)
	}

	if projectID != "" {
		if projectID == "all" {
			specifiedRole := roles.Role{
				ID: roleID,
			}
			err = roles.CheckAllResourcesPermission(identityClient, c.DomainID, groupID, roleID).ExtractErr()
			return specifiedRole, err
		}

		return iam.GetGroupRoleAssignmentWithProjectID(identityClient, groupID, roleID, projectID)
	}

	if enterpriseProjectID != "" {
		return iam.GetGroupRoleAssignmentWithEpsID(iamClient, groupID, roleID, enterpriseProjectID)
	}

	return nil, golangsdk.ErrDefault404{}
}

func TestAccIdentityGroupRoleAssignment_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "flexibleengine_identity_group_role_assignment.test"
	var role roles.Role

	rc := acceptance.InitResourceCheck(
		resourceName,
		&role,
		getIdentityGroupRoleAssignmentResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
			testAccPrecheckDomainId(t)
		},
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityGroupRoleAssignment_domain(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "group_id",
						"flexibleengine_identity_group_v3.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "role_id",
						"flexibleengine_identity_role_v3.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "domain_id",
						OS_DOMAIN_ID),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccIdentityGroupRoleAssignmentDomainImportStateFunc(resourceName),
			},
		},
	})
}

func TestAccIdentityGroupRoleAssignment_project(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "flexibleengine_identity_group_role_assignment.test"
	var role roles.Role

	rc := acceptance.InitResourceCheck(
		resourceName,
		&role,
		getIdentityGroupRoleAssignmentResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
			testAccPreCheckProjectID(t)
		},
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityGroupRoleAssignment_project(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),

					resource.TestCheckResourceAttrPair(resourceName, "group_id",
						"flexibleengine_identity_group_v3.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "role_id",
						"flexibleengine_identity_role_v3.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "project_id",
						OS_PROJECT_ID),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccIdentityGroupRoleAssignmentProjectImportStateFunc(resourceName),
			},
		},
	})
}

func TestAccIdentityGroupRoleAssignment_allProjects(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "flexibleengine_identity_group_role_assignment.test"
	var role roles.Role

	rc := acceptance.InitResourceCheck(
		resourceName,
		&role,
		getIdentityGroupRoleAssignmentResourceFunc,
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
				Config: testAccIdentityGroupRoleAssignment_allProjects(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "project_id", "all"),
					resource.TestCheckResourceAttrPair(resourceName, "group_id",
						"flexibleengine_identity_group_v3.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "role_id",
						"flexibleengine_identity_role_v3.test", "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccIdentityGroupRoleAssignmentProjectImportStateFunc(resourceName),
			},
		},
	})
}

func TestAccIdentityGroupRoleAssignment_epsID(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "flexibleengine_identity_group_role_assignment.test"
	var role eps_permissions.Role

	rc := acceptance.InitResourceCheck(
		resourceName,
		&role,
		getIdentityGroupRoleAssignmentResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
			testAccPreCheckEpsID(t)
		},
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityGroupRoleAssignment_epsID(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "group_id",
						"flexibleengine_identity_group_v3.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "role_id",
						"flexibleengine_identity_role_v3.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id",
						OS_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccIdentityGroupRoleAssignmentEpsImportStateFunc(resourceName),
			},
		},
	})
}

func testAccIdentityGroupRoleAssignmentDomainImportStateFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceName, rs)
		}
		if rs.Primary.Attributes["group_id"] == "" ||
			rs.Primary.Attributes["role_id"] == "" || rs.Primary.Attributes["domain_id"] == "" {
			return "", fmt.Errorf("invalid format specified for import ID,"+
				" want '<group_id>/<role_id>/<domain_id>', but got '%s/%s/%s'",
				rs.Primary.Attributes["group_id"], rs.Primary.Attributes["role_id"],
				rs.Primary.Attributes["domain_id"])
		}
		return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["group_id"],
			rs.Primary.Attributes["role_id"], rs.Primary.Attributes["domain_id"]), nil
	}
}

func testAccIdentityGroupRoleAssignmentProjectImportStateFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceName, rs)
		}
		if rs.Primary.Attributes["group_id"] == "" ||
			rs.Primary.Attributes["role_id"] == "" || rs.Primary.Attributes["project_id"] == "" {
			return "", fmt.Errorf("invalid format specified for import ID,"+
				" want '<group_id>/<role_id>/<project_id>', but got '%s/%s/%s'",
				rs.Primary.Attributes["group_id"], rs.Primary.Attributes["role_id"],
				rs.Primary.Attributes["project_id"])
		}
		return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["group_id"],
			rs.Primary.Attributes["role_id"], rs.Primary.Attributes["project_id"]), nil
	}
}

func testAccIdentityGroupRoleAssignmentEpsImportStateFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceName, rs)
		}
		if rs.Primary.Attributes["group_id"] == "" ||
			rs.Primary.Attributes["role_id"] == "" || rs.Primary.Attributes["enterprise_project_id"] == "" {
			return "", fmt.Errorf("invalid format specified for import ID,"+
				" want '<group_id>/<role_id>/<enterprise_project_id>', but got '%s/%s/%s'",
				rs.Primary.Attributes["group_id"], rs.Primary.Attributes["role_id"],
				rs.Primary.Attributes["enterprise_project_id"])
		}
		return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["group_id"],
			rs.Primary.Attributes["role_id"], rs.Primary.Attributes["enterprise_project_id"]), nil
	}
}

func testAccIdentityGroupRoleAssignment_base(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_identity_role_v3" test {
  name        = "%[1]s"
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

resource "flexibleengine_identity_group_v3" "test" {
  name        = "%[1]s"
  description = "A test group"
}`, rName)
}

func testAccIdentityGroupRoleAssignment_domain(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_identity_group_role_assignment" "test" {
  group_id  = flexibleengine_identity_group_v3.test.id
  role_id   = flexibleengine_identity_role_v3.test.id
  domain_id = "%s"
}
`, testAccIdentityGroupRoleAssignment_base(rName), OS_DOMAIN_ID)
}

func testAccIdentityGroupRoleAssignment_project(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_identity_group_role_assignment" "test" {
  group_id   = flexibleengine_identity_group_v3.test.id
  role_id    = flexibleengine_identity_role_v3.test.id
  project_id = "%s"
}
`, testAccIdentityGroupRoleAssignment_base(rName), OS_PROJECT_ID)
}

func testAccIdentityGroupRoleAssignment_allProjects(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_identity_group_role_assignment" "test" {
  group_id   = flexibleengine_identity_group_v3.test.id
  role_id    = flexibleengine_identity_role_v3.test.id
  project_id = "all"
}
`, testAccIdentityGroupRoleAssignment_base(rName))
}

func testAccIdentityGroupRoleAssignment_epsID(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_identity_group_role_assignment" "test" {
  group_id              = flexibleengine_identity_group_v3.test.id
  role_id               = flexibleengine_identity_role_v3.test.id
  enterprise_project_id = "%s"
}
`, testAccIdentityGroupRoleAssignment_base(rName), OS_ENTERPRISE_PROJECT_ID_TEST)
}
