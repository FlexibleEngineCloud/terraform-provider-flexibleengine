package acceptance

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceModelartsDatasets_basic(t *testing.T) {
	dataSourceName := "data.flexibleengine_modelarts_datasets.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	name := acceptance.RandomAccResourceName()
	obsName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDatasets_basic(name, obsName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(dataSourceName, "datasets.0.id",
						"flexibleengine_modelarts_dataset.test", "id"),
					resource.TestCheckResourceAttr(dataSourceName, "datasets.0.name", name),
					resource.TestCheckResourceAttr(dataSourceName, "datasets.0.type", "1"),
					resource.TestCheckResourceAttrPair(dataSourceName, "datasets.0.description",
						"flexibleengine_modelarts_dataset.test", "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "datasets.0.output_path",
						"flexibleengine_modelarts_dataset.test", "output_path"),
					resource.TestCheckResourceAttr(dataSourceName, "datasets.0.data_source.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "datasets.0.labels.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "datasets.0.schemas.#", "0"),
				),
			},
		},
	})
}

func testAccDataSourceDatasets_basic(rName, obsName string) string {
	datasets := testAccDateset_basic(rName, obsName)
	return fmt.Sprintf(`
%s

data "flexibleengine_modelarts_datasets" "test" {
  name = "%s"
  type = 1

  depends_on = [
    flexibleengine_modelarts_dataset.test
  ]
}
`, datasets, rName)
}
