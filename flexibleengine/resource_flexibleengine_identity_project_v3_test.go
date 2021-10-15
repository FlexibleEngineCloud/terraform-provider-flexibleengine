package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/identity/v3/projects"
)

func TestAccIdentityProjectV3_basic(t *testing.T) {
	var project projects.Project
	var projectName = fmt.Sprintf("eu-west-0-ACCPTTEST-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIdentityProjectV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProjectV3_basic(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityProjectV3Exists("flexibleengine_identity_project_v3.project_1", &project),
					resource.TestCheckResourceAttrPtr(
						"flexibleengine_identity_project_v3.project_1", "name", &project.Name),
					resource.TestCheckResourceAttrPtr(
						"flexibleengine_identity_project_v3.project_1", "description", &project.Description),
					resource.TestCheckResourceAttrPtr(
						"flexibleengine_identity_project_v3.project_1", "domain_id", &project.DomainID),
					resource.TestCheckResourceAttrPtr(
						"flexibleengine_identity_project_v3.project_1", "parent_id", &project.ParentID),
				),
			},
			{
				Config: testAccIdentityProjectV3_update(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityProjectV3Exists("flexibleengine_identity_project_v3.project_1", &project),
					resource.TestCheckResourceAttrPtr(
						"flexibleengine_identity_project_v3.project_1", "name", &project.Name),
					resource.TestCheckResourceAttrPtr(
						"flexibleengine_identity_project_v3.project_1", "description", &project.Description),
					resource.TestCheckResourceAttrPtr(
						"flexibleengine_identity_project_v3.project_1", "domain_id", &project.DomainID),
					resource.TestCheckResourceAttrPtr(
						"flexibleengine_identity_project_v3.project_1", "parent_id", &project.ParentID),
				),
			},
		},
	})
}

func testAccCheckIdentityProjectV3Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	identityClient, err := config.identityV3Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine identity client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_identity_project_v3" {
			continue
		}

		_, err := projects.Get(identityClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Project still exists")
		}
	}

	return nil
}

func testAccCheckIdentityProjectV3Exists(n string, project *projects.Project) resource.TestCheckFunc {
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

		found, err := projects.Get(identityClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Project not found")
		}

		*project = *found

		return nil
	}
}

func testAccIdentityProjectV3_basic(projectName string) string {
	return fmt.Sprintf(`
    resource "flexibleengine_identity_project_v3" "project_1" {
      name = "%s"
      description = "A ACC test project"
    }
  `, projectName)
}

func testAccIdentityProjectV3_update(projectName string) string {
	return fmt.Sprintf(`
    resource "flexibleengine_identity_project_v3" "project_1" {
      name = "%s"
      description = "Some Project"
    }
  `, projectName)
}
