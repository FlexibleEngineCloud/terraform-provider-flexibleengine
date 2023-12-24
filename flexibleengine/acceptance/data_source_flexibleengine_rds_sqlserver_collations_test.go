package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceSQLServerCollations_basic(t *testing.T) {
	rName := "data.flexibleengine_rds_sqlserver_collations.collations"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceSQLServerCollations_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "char_sets.#"),
				),
			},
		},
	})
}

func testAccDatasourceSQLServerCollations_basic() string {
	return `
data "flexibleengine_rds_sqlserver_collations" "collations" {}
`
}
