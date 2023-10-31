package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceStoragetype_basic(t *testing.T) {
	rName := "data.flexibleengine_rds_storage_types.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceStoragetype_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "db_type", "MySQL"),
					resource.TestCheckResourceAttr(rName, "db_version", "8.0"),
					resource.TestCheckResourceAttrSet(rName, "storage_types.0.name"),
					resource.TestCheckResourceAttrSet(rName, "storage_types.0.az_status.%"),
					resource.TestCheckResourceAttrSet(rName, "storage_types.0.support_compute_group_type.#"),
				),
			},
		},
	})
}

func TestAccDatasourceStoragetype_PostgreSQL_basic(t *testing.T) {
	rName := "data.flexibleengine_rds_storage_types.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceStoragetype_PostgreSQL_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "db_type", "PostgreSQL"),
					resource.TestCheckResourceAttr(rName, "db_version", "14"),
					resource.TestCheckResourceAttrSet(rName, "storage_types.0.name"),
					resource.TestCheckResourceAttrSet(rName, "storage_types.0.az_status.%"),
					resource.TestCheckResourceAttrSet(rName, "storage_types.0.support_compute_group_type.#"),
				),
			},
		},
	})
}

func TestAccDatasourceStoragetype_SQLServer_basic(t *testing.T) {
	rName := "data.flexibleengine_rds_storage_types.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceStoragetype_SQLServer_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "db_type", "SQLServer"),
					resource.TestCheckResourceAttr(rName, "db_version", "2019_SE"),
					resource.TestCheckResourceAttrSet(rName, "storage_types.0.name"),
					resource.TestCheckResourceAttrSet(rName, "storage_types.0.az_status.%"),
					resource.TestCheckResourceAttrSet(rName, "storage_types.0.support_compute_group_type.#"),
				),
			},
		},
	})
}

func testAccDatasourceStoragetype_basic() string {
	return `
data "flexibleengine_rds_storage_types" "test" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "replica"
}`
}

func testAccDatasourceStoragetype_PostgreSQL_basic() string {
	return `
data "flexibleengine_rds_storage_types" "test" {
  db_type       = "PostgreSQL"
  db_version    = "14"
  instance_mode = "ha"
}`
}

func testAccDatasourceStoragetype_SQLServer_basic() string {
	return `
data "flexibleengine_rds_storage_types" "test" {
  db_type       = "SQLServer"
  db_version    = "2019_SE"
  instance_mode = "single"
}`
}
