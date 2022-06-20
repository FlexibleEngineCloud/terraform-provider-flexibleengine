package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/swr/v2/repositories"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getResourceRepository(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	swrClient, err := conf.SwrV2Client(OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("Error creating SWR client: %s", err)
	}

	return repositories.Get(swrClient, state.Primary.Attributes["organization"], state.Primary.ID).Extract()
}

func TestAccSWRRepository_basic(t *testing.T) {
	var repo repositories.ImageRepository
	rName := acceptance.RandomAccResourceName()
	resourceName := "flexibleengine_swr_repository.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&repo,
		getResourceRepository,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSWRRepository_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "category", "linux"),
					resource.TestCheckResourceAttr(resourceName, "is_public", "false"),
					resource.TestCheckResourceAttrPair(resourceName, "organization", "flexibleengine_swr_organization.test", "name"),
				),
			},
			{
				Config: testAccSWRRepository_update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "category", "windows"),
					resource.TestCheckResourceAttr(resourceName, "is_public", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSWRRepositoryImportStateIdFunc(),
			},
		},
	})
}

func testAccSWRRepositoryImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var organization string
		var repositoryID string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "flexibleengine_swr_organization" {
				organization = rs.Primary.Attributes["name"]
			} else if rs.Type == "flexibleengine_swr_repository" {
				repositoryID = rs.Primary.ID
			}
		}
		if organization == "" || repositoryID == "" {
			return "", fmt.Errorf("resource not found: %s/%s", organization, repositoryID)
		}
		return fmt.Sprintf("%s/%s", organization, repositoryID), nil
	}
}

func testAccSWRRepository_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_swr_repository" "test" {
  organization = flexibleengine_swr_organization.test.name
  name         = "%s"
  description  = "Test repository"
  category     = "linux"
  is_public    = false
}
`, testAccSWROrganization_basic(rName), rName)
}

func testAccSWRRepository_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_swr_repository" "test" {
  organization = flexibleengine_swr_organization.test.name
  name         = "%s"
  description  = "Test repository"
  category     = "windows"
  is_public    = true
}
`, testAccSWROrganization_basic(rName), rName)
}
