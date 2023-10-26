package acceptance

import (
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCssFlavorsDataSource_basic(t *testing.T) {
	var (
		typeFilter        = acceptance.InitDataSourceCheck("data.flexibleengine_css_flavors.type_filter")
		typeFilterCold    = acceptance.InitDataSourceCheck("data.flexibleengine_css_flavors.type_filter_cold")
		typeFilterMaster  = acceptance.InitDataSourceCheck("data.flexibleengine_css_flavors.type_filter_master")
		typeFilterClient  = acceptance.InitDataSourceCheck("data.flexibleengine_css_flavors.type_filter_client")
		versionFilter654  = acceptance.InitDataSourceCheck("data.flexibleengine_css_flavors.version_filter_654")
		versionFilter711  = acceptance.InitDataSourceCheck("data.flexibleengine_css_flavors.version_filter_711")
		versionFilter762  = acceptance.InitDataSourceCheck("data.flexibleengine_css_flavors.version_filter_762")
		versionFilter793  = acceptance.InitDataSourceCheck("data.flexibleengine_css_flavors.version_filter_793")
		versionFilter7102 = acceptance.InitDataSourceCheck("data.flexibleengine_css_flavors.version_filter_7102")
		vcpusFilter       = acceptance.InitDataSourceCheck("data.flexibleengine_css_flavors.vcpus_filter")
		memoryFilter      = acceptance.InitDataSourceCheck("data.flexibleengine_css_flavors.memory_filter")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCssFlavors_basic,
				Check: resource.ComposeTestCheckFunc(
					typeFilter.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					typeFilterCold.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_cold_useful", "true"),
					typeFilterMaster.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_master_useful", "true"),
					typeFilterClient.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_client_useful", "true"),
					versionFilter654.CheckResourceExists(),
					resource.TestCheckOutput("is_version_filter_654_useful", "true"),
					versionFilter711.CheckResourceExists(),
					resource.TestCheckOutput("is_version_filter_711_useful", "true"),
					versionFilter762.CheckResourceExists(),
					resource.TestCheckOutput("is_version_filter_762_useful", "true"),
					versionFilter793.CheckResourceExists(),
					resource.TestCheckOutput("is_version_filter_793_useful", "true"),
					versionFilter7102.CheckResourceExists(),
					resource.TestCheckOutput("is_version_filter_7102_useful", "true"),
					vcpusFilter.CheckResourceExists(),
					resource.TestCheckOutput("is_vcpus_filter_useful", "true"),
					memoryFilter.CheckResourceExists(),
					resource.TestCheckOutput("is_memory_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceCssFlavors_basic = `

data "flexibleengine_css_flavors" "type_filter" {
  type = "ess"
}

output "is_type_filter_useful" {
  value = !contains([for v in data.flexibleengine_css_flavors.type_filter.flavors[*].type : v == "ess"], "false")
}

data "flexibleengine_css_flavors" "type_filter_cold" {
  type = "ess-cold"
}

output "is_type_filter_cold_useful" {
  value = !contains([for v in data.flexibleengine_css_flavors.type_filter_cold.flavors[*].type : v == "ess-cold"], "false")
}

data "flexibleengine_css_flavors" "type_filter_master" {
  type = "ess-master"
}

output "is_type_filter_master_useful" {
  value = !contains([for v in data.flexibleengine_css_flavors.type_filter_master.flavors[*].type : v == "ess-master"], "false")
}

data "flexibleengine_css_flavors" "type_filter_client" {
  type = "ess-client"
}

output "is_type_filter_client_useful" {
  value = !contains([for v in data.flexibleengine_css_flavors.type_filter_client.flavors[*].type : v == "ess-client"], "false")
}

data "flexibleengine_css_flavors" "version_filter_654" {
  version = "6.5.4"
}

output "is_version_filter_654_useful" {
  value = !contains([for v in data.flexibleengine_css_flavors.version_filter_654.flavors[*].version : v == "6.5.4"], "false")
}

data "flexibleengine_css_flavors" "version_filter_711" {
  version = "7.1.1"
}

output "is_version_filter_711_useful" {
  value = !contains([for v in data.flexibleengine_css_flavors.version_filter_711.flavors[*].version : v == "7.1.1"], "false")
}

data "flexibleengine_css_flavors" "version_filter_762" {
  version = "7.6.2"
}

output "is_version_filter_762_useful" {
  value = !contains([for v in data.flexibleengine_css_flavors.version_filter_762.flavors[*].version : v == "7.6.2"], "false")
}

data "flexibleengine_css_flavors" "version_filter_793" {
  version = "7.9.3"
}

output "is_version_filter_793_useful" {
  value = !contains([for v in data.flexibleengine_css_flavors.version_filter_793.flavors[*].version : v == "7.9.3"], "false")
}

data "flexibleengine_css_flavors" "version_filter_7102" {
  version = "7.10.2"
}

output "is_version_filter_7102_useful" {
  value = !contains([for v in data.flexibleengine_css_flavors.version_filter_7102.flavors[*].version : v == "7.10.2"], "false")
}

data "flexibleengine_css_flavors" "vcpus_filter" {
  vcpus = 32
}

output "is_vcpus_filter_useful" {
  value = !contains([for v in data.flexibleengine_css_flavors.vcpus_filter.flavors[*].vcpus : v == 32], "false")
}

data "flexibleengine_css_flavors" "memory_filter" {
  memory = 256
}

output "is_memory_filter_useful" {
  value = !contains([for v in data.flexibleengine_css_flavors.memory_filter.flavors[*].memory : v == 256], "false")
}
`

func TestAccCssFlavorsDataSource_all(t *testing.T) {
	dataSourceName := "data.flexibleengine_css_flavors.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCssFlavors_version_654,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.type", "ess"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.version", "6.5.4"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.id"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.region", "eu-west-0"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.name", "ess.spec-4u8g"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.memory", "8"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.vcpus", "4"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.disk_range"),
				),
			},
			{
				Config: testAccDataSourceCssFlavors_version_711,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.type", "ess-cold"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.version", "7.1.1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.id"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.region", "eu-west-0"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.name", "ess.spec-4u8g"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.memory", "8"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.vcpus", "4"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.disk_range"),
				),
			},
			{
				Config: testAccDataSourceCssFlavors_version_762,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.type", "ess-master"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.version", "7.6.2"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.id"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.region", "eu-west-0"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.name", "ess.spec-4u8g"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.memory", "8"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.vcpus", "4"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.disk_range"),
				),
			},
			{
				Config: testAccDataSourceCssFlavors_version_793,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.type", "ess-client"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.version", "7.9.3"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.id"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.region", "eu-west-0"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.name", "ess.spec-4u8g"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.memory", "8"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.vcpus", "4"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.disk_range"),
				),
			},
			{
				Config: testAccDataSourceCssFlavors_version_7102,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.type", "ess"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.version", "7.10.2"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.id"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.region", "eu-west-0"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.name", "ess.spec-4u8g"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.memory", "8"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.vcpus", "4"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.disk_range"),
				),
			},
		},
	})
}

const testAccDataSourceCssFlavors_version_654 = `
data "flexibleengine_css_flavors" "test" {
  type    = "ess"
  version = "6.5.4"
  vcpus   = 4
  memory  = 8
  region  = "eu-west-0"
  name    = "ess.spec-4u8g"
}
`

const testAccDataSourceCssFlavors_version_711 = `
data "flexibleengine_css_flavors" "test" {
  type    = "ess-cold"
  version = "7.1.1"
  vcpus   = 4
  memory  = 8
  region  = "eu-west-0"
  name    = "ess.spec-4u8g"
}
`

const testAccDataSourceCssFlavors_version_762 = `
data "flexibleengine_css_flavors" "test" {
  type    = "ess-master"
  version = "7.6.2"
  vcpus   = 4
  memory  = 8
  region  = "eu-west-0"
  name    = "ess.spec-4u8g"
}
`

const testAccDataSourceCssFlavors_version_793 = `
data "flexibleengine_css_flavors" "test" {
  type    = "ess-client"
  version = "7.9.3"
  vcpus   = 4
  memory  = 8
  region  = "eu-west-0"
  name    = "ess.spec-4u8g"
}
`

const testAccDataSourceCssFlavors_version_7102 = `
data "flexibleengine_css_flavors" "test" {
  type    = "ess"
  version = "7.10.2"
  vcpus   = 4
  memory  = 8
  region  = "eu-west-0"
  name    = "ess.spec-4u8g"
}
`
