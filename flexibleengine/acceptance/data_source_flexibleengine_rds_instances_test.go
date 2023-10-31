package acceptance

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRdsInstanceDataSource_basic(t *testing.T) {
	dataSourceName := "data.flexibleengine_rds_instances.test"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "instances.#", regexp.MustCompile("\\d+")),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.name"),
				),
			},
		},
	})
}

func TestAccRdsInstanceDataSource_ha_basic(t *testing.T) {
	dataSourceName := "data.flexibleengine_rds_instances.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceDataSource_ha_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "instances.#", regexp.MustCompile("\\d+")),
				),
			},
		},
	})
}

func TestAccRdsInstanceDataSource_replica_basic(t *testing.T) {
	dataSourceName := "data.flexibleengine_rds_instances.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceDataSource_replica_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "instances.#", regexp.MustCompile("\\d+")),
				),
			},
		},
	})
}

func TestAccRdsInstanceDataSource_enterprise_basic(t *testing.T) {
	dataSourceName := "data.flexibleengine_rds_instances.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceDataSource_enterprise_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "instances.#", regexp.MustCompile("\\d+")),
				),
			},
		},
	})
}

func testAccRdsInstanceDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_availability_zones" "test" {}

data "flexibleengine_networking_secgroup_v2" "test" {
  name = "default"
}

data "flexibleengine_rds_flavors_v3" "test" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "single"
}

resource "flexibleengine_rds_instance_v3" "test" {
  name              = "%s"
  flavor            = data.flexibleengine_rds_flavors_v3.test.flavors[2].name
  availability_zone = [data.flexibleengine_availability_zones.test.names[0]]
  security_group_id = data.flexibleengine_networking_secgroup_v2.test.id
  subnet_id         = flexibleengine_vpc_subnet_v1.test.id
  vpc_id            = flexibleengine_vpc_v1.test.id
  time_zone         = "UTC+08:00"
  fixed_ip          = "192.168.0.58"

  db {
    password = "FlexibleEngine!120521"
    type     = "MySQL"
    version  = "8.0"
    port     = 8630
  }

  volume {
    type = "COMMON"
    size = 50
  }

  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }

  tags = {
    key = "value"
    foo = "bar"
  }
}

data "flexibleengine_rds_instances" "test" {
  depends_on = [
    flexibleengine_rds_instance_v3.test,
  ]
}
`, testVpc(rName), rName)
}

func testAccRdsInstanceDataSource_ha_basic() string {
	return fmt.Sprintf(`
data "flexibleengine_rds_instances" "test" {
  type           = "Ha"
  datastore_type = "PostgreSQL"
}
`)
}

func testAccRdsInstanceDataSource_replica_basic() string {
	return fmt.Sprintf(`
data "flexibleengine_rds_instances" "test" {
  type           = "Replica"
  datastore_type = "PostgreSQL"
}
`)
}

func testAccRdsInstanceDataSource_enterprise_basic() string {
	return fmt.Sprintf(`
data "flexibleengine_rds_instances" "test" {
  type                  = "Enterprise"
  datastore_type        = "SQLServer"
  enterprise_project_id = "0"
}
`)
}
