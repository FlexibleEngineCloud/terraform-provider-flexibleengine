package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/eps/v1/enterpriseprojects"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getResourceEnterpriseProject(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	epsClient, err := config.EnterpriseProjectClient(OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("Unable to create EPS client : %s", err)
	}

	return enterpriseprojects.Get(epsClient, state.Primary.ID).Extract()

}

func TestAccEnterpriseProject_basic(t *testing.T) {
	var project enterpriseprojects.Project

	rName := acceptance.RandomAccResourceName()
	updateName := rName + "update"
	resourceName := "flexibleengine_enterprise_project.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&project,
		getResourceEnterpriseProject,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckEnterpriseProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEnterpriseProject_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform test"),
					resource.TestCheckResourceAttr(resourceName, "status", "1"),
				),
			},
			{
				Config: testAccEnterpriseProject_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform test update"),
					resource.TestCheckResourceAttr(resourceName, "status", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckEnterpriseProjectDestroy(s *terraform.State) error {
	conf := testAccProvider.Meta().(*config.Config)
	epsClient, err := conf.EnterpriseProjectClient(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Unable to create EPS client : %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_enterprise_project" {
			continue
		}

		project, err := enterpriseprojects.Get(epsClient, rs.Primary.ID).Extract()
		if err == nil {
			if project.Status != 2 {
				return fmt.Errorf("Project still active")
			}
		}
	}

	return nil
}

func testAccEnterpriseProject_basic(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_enterprise_project" "test" {
  name        = "%s"
  description = "terraform test"
}`, rName)
}

func testAccEnterpriseProject_update(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_enterprise_project" "test" {
  name        = "%s"
  description = "terraform test update"
}`, rName)
}
