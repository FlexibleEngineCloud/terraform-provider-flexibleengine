package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGeminiDBInstancesDataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	password := acceptance.RandomPassword()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiDBInstancesDataSource_basic(rName, password),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGeminiDBInstancesDataSourceID("data.flexibleengine_gaussdb_cassandra_instances.test"),
					resource.TestCheckResourceAttr("data.flexibleengine_gaussdb_cassandra_instances.test", "instances.#", "1"),
				),
			},
		},
	})
}

func testAccCheckGeminiDBInstancesDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find GaussDB cassandra instance data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("GaussDB cassandra instances data source ID not set ")
		}

		return nil
	}
}

func testAccGeminiDBInstancesDataSource_basic(rName, password string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_gaussdb_cassandra_instances" "test" {
  name = "%s"
  depends_on = [
    flexibleengine_gaussdb_cassandra_instance.test,
  ]
}
`, testAccGeminiDBInstanceConfig_basic(rName, password), rName)
}
