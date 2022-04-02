package flexibleengine

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFunctionGraphDependencies_basic(t *testing.T) {
	dataSourceName := "data.flexibleengine_fgs_dependencies.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionGraphDependencies_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceName, "packages.#", regexp.MustCompile(`[1-9][0-9]*`)),
				),
			},
		},
	})
}

func TestAccFunctionGraphDependencies_name(t *testing.T) {
	dataSourceName := "data.flexibleengine_fgs_dependencies.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionGraphDependencies_name(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "type", "public"),
					resource.TestCheckResourceAttr(dataSourceName, "name", "obssdk"),
					resource.TestCheckResourceAttr(dataSourceName, "packages.#", "1"),
				),
			},
		},
	})
}

func TestAccFunctionGraphDependencies_runtime(t *testing.T) {
	dataSourceName := "data.flexibleengine_fgs_dependencies.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionGraphDependencies_runtime(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "type", "public"),
					resource.TestCheckResourceAttr(dataSourceName, "runtime", "Python2.7"),
					resource.TestMatchResourceAttr(dataSourceName, "packages.#", regexp.MustCompile(`[1-9][0-9]*`)),
				),
			},
		},
	})
}

func testAccFunctionGraphDependencies_basic() string {
	return fmt.Sprintf(`
data "flexibleengine_fgs_dependencies" "test" {}
`)
}

func testAccFunctionGraphDependencies_name() string {
	return fmt.Sprintf(`
data "flexibleengine_fgs_dependencies" "test" {
  type = "public"
  name = "obssdk"
}
`)
}

func testAccFunctionGraphDependencies_runtime() string {
	return fmt.Sprintf(`
data "flexibleengine_fgs_dependencies" "test" {
  type    = "public"
  runtime = "Python2.7"
}
`)
}
