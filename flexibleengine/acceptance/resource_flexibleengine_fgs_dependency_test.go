package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/fgs/v2/dependencies"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getDependencyResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.FgsV2Client(OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating FunctionGraph v2 client: %s", err)
	}
	return dependencies.Get(c, state.Primary.ID)
}

func TestAccFunctionGraphResourceDependency_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "flexibleengine_fgs_dependency.test"
	pkgLocation := fmt.Sprintf("https://%s.oss.%s.prod-cloud-ocb.orange-business.com/dependencies/sdkcore.zip",
		OS_FGS_BUCKET, OS_REGION_NAME)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckOBS(t)
		},
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionGraphResourceDependency_basic(rName, pkgLocation),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(resourceName, "runtime", "Python2.7"),
					resource.TestCheckResourceAttr(resourceName, "link", pkgLocation),
				),
			},
			{
				Config: testAccFunctionGraphResourceDependency_update(rName, pkgLocation),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName+"_update"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated by terraform script"),
					resource.TestCheckResourceAttr(resourceName, "runtime", "Python3.6"),
					resource.TestCheckResourceAttr(resourceName, "link", pkgLocation),
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

func testAccFunctionGraphResourceDependency_basic(rName, pkgLocation string) string {
	return fmt.Sprintf(`
resource "flexibleengine_fgs_dependency" "test" {
  name        = "%s"
  description = "Created by terraform script"
  runtime     = "Python2.7"
  link        = "%s"
}
`, rName, pkgLocation)
}

func testAccFunctionGraphResourceDependency_update(rName, pkgLocation string) string {
	return fmt.Sprintf(`
resource "flexibleengine_fgs_dependency" "test" {
  name        = "%s_update"
  description = "Updated by terraform script"
  runtime     = "Python3.6"
  link        = "%s"
}
`, rName, pkgLocation)
}
