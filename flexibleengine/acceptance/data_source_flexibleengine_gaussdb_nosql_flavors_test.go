package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGaussDBNoSQLFlavors_basic(t *testing.T) {
	dataSourceName := "data.flexibleengine_gaussdb_nosql_flavors.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDBNoSQLFlavors_default(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.engine", "cassandra"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.#"),
				),
			},
			{
				Config: testAccGaussDBNoSQLFlavors_cassandra(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "engine", "cassandra"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.#"),
				),
			},
		},
	})
}

func TestAccGaussDBNoSQLFlavors_influxdb(t *testing.T) {
	dataSourceName := "data.flexibleengine_gaussdb_nosql_flavors.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDBNoSQLFlavors_influxdb(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "engine", "influxdb"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.#"),
				),
			},
		},
	})
}

func TestAccGaussDBNoSQLFlavors_vcpus(t *testing.T) {
	dataSourceName := "data.flexibleengine_gaussdb_nosql_flavors.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDBNoSQLFlavors_vcpus(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.vcpus", "4"),
				),
			},
		},
	})
}

func TestAccGaussDBNoSQLFlavors_memory(t *testing.T) {
	dataSourceName := "data.flexibleengine_gaussdb_nosql_flavors.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDBNoSQLFlavors_memory(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.memory", "16"),
				),
			},
		},
	})
}

func TestAccGaussDBNoSQLFlavors_az(t *testing.T) {
	dataSourceName := "data.flexibleengine_gaussdb_nosql_flavors.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDBNoSQLFlavors_az(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(dataSourceName, "flavors.0.availability_zones.0",
						"data.flexibleengine_availability_zones.test", "names.0"),
				),
			},
		},
	})
}

func testAccGaussDBNoSQLFlavors_default() string {
	return fmt.Sprintf(`
data "flexibleengine_gaussdb_nosql_flavors" "test" {}
`)
}

func testAccGaussDBNoSQLFlavors_cassandra() string {
	return fmt.Sprintf(`
data "flexibleengine_gaussdb_nosql_flavors" "test" {
  engine = "cassandra"
}
`)
}

func testAccGaussDBNoSQLFlavors_influxdb() string {
	return fmt.Sprintf(`
data "flexibleengine_gaussdb_nosql_flavors" "test" {
  engine = "influxdb"
}
`)
}

func testAccGaussDBNoSQLFlavors_vcpus() string {
	return fmt.Sprintf(`
data "flexibleengine_gaussdb_nosql_flavors" "test" {
  vcpus = 4
}
`)
}

func testAccGaussDBNoSQLFlavors_memory() string {
	return fmt.Sprintf(`
data "flexibleengine_gaussdb_nosql_flavors" "test" {
  memory = 16
}
`)
}

func testAccGaussDBNoSQLFlavors_az() string {
	return fmt.Sprintf(`
data "flexibleengine_availability_zones" "test" {}

data "flexibleengine_gaussdb_nosql_flavors" "test" {
  availability_zone = data.flexibleengine_availability_zones.test.names[0]
}
`)
}
