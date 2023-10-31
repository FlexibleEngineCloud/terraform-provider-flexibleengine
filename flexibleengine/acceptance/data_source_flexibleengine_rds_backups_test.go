package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceBackup_basic(t *testing.T) {
	rName := "data.flexibleengine_rds_backups.test"
	dc := acceptance.InitDataSourceCheck(rName)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceBackup_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "backups.0.id", "flexibleengine_rds_backup.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "backups.0.name", "flexibleengine_rds_backup.test", "name"),
					resource.TestCheckResourceAttrPair(rName, "backups.0.instance_id",
						"flexibleengine_rds_instance_v3.test", "id"),
					resource.TestCheckResourceAttr(rName, "backups.0.type", "manual"),
					resource.TestCheckResourceAttrSet(rName, "backups.0.size"),
					resource.TestCheckResourceAttrSet(rName, "backups.0.status"),
					resource.TestCheckResourceAttrSet(rName, "backups.0.begin_time"),
					resource.TestCheckResourceAttrSet(rName, "backups.0.end_time"),
					resource.TestCheckResourceAttrSet(rName, "backups.0.associated_with_ddm"),
					resource.TestCheckResourceAttr(rName, "backups.0.datastore.#", "1"),
					resource.TestCheckResourceAttr(rName, "backups.0.databases.#", "0"),
				),
			},
		},
	})
}

func TestAccDatasourceBackup_auto_basic(t *testing.T) {
	rName := "data.flexibleengine_rds_backups.test"
	dc := acceptance.InitDataSourceCheck(rName)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceBackup_auto_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "backups.0.id", "flexibleengine_rds_backup.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "backups.0.name", "flexibleengine_rds_backup.test", "name"),
					resource.TestCheckResourceAttrPair(rName, "backups.0.instance_id",
						"flexibleengine_rds_instance_v3.test", "id"),
					resource.TestCheckResourceAttr(rName, "backups.0.type", "auto"),
					resource.TestCheckResourceAttrSet(rName, "backups.0.size"),
					resource.TestCheckResourceAttrSet(rName, "backups.0.status"),
					resource.TestCheckResourceAttrSet(rName, "backups.0.begin_time"),
					resource.TestCheckResourceAttrSet(rName, "backups.0.end_time"),
					resource.TestCheckResourceAttrSet(rName, "backups.0.associated_with_ddm"),
					resource.TestCheckResourceAttr(rName, "backups.0.datastore.#", "1"),
					resource.TestCheckResourceAttr(rName, "backups.0.databases.#", "0"),
				),
			},
		},
	})
}

func TestAccDatasourceBackup_incremental_basic(t *testing.T) {
	rName := "data.flexibleengine_rds_backups.test"
	dc := acceptance.InitDataSourceCheck(rName)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceBackup_incremental_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "backups.0.id", "flexibleengine_rds_backup.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "backups.0.name", "flexibleengine_rds_backup.test", "name"),
					resource.TestCheckResourceAttrPair(rName, "backups.0.instance_id",
						"flexibleengine_rds_instance_v3.test", "id"),
					resource.TestCheckResourceAttr(rName, "backups.0.type", "incremental"),
					resource.TestCheckResourceAttrSet(rName, "backups.0.size"),
					resource.TestCheckResourceAttrSet(rName, "backups.0.status"),
					resource.TestCheckResourceAttrSet(rName, "backups.0.begin_time"),
					resource.TestCheckResourceAttrSet(rName, "backups.0.end_time"),
					resource.TestCheckResourceAttrSet(rName, "backups.0.associated_with_ddm"),
					resource.TestCheckResourceAttr(rName, "backups.0.datastore.#", "1"),
					resource.TestCheckResourceAttr(rName, "backups.0.databases.#", "0"),
				),
			},
		},
	})
}

func testAccDatasourceBackup_base(name string) string {
	return fmt.Sprintf(`
%[1]s

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
  name              = "%[2]s"
  flavor            = data.flexibleengine_rds_flavors_v3.test.flavors[2].name
  availability_zone = [data.flexibleengine_availability_zones.test.names[0]]
  security_group_id = data.flexibleengine_networking_secgroup_v2.test.id
  subnet_id         = flexibleengine_vpc_subnet_v1.test.id
  vpc_id            = flexibleengine_vpc_v1.test.id
  time_zone         = "UTC+08:00"

  db {
    password = "FlexibleEngine!120521"
    type     = "MySQL"
    version  = "8.0"
    port     = 8630
  }

  volume {
    type = "COMMON"
    size = 60
  }

  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }

  lifecycle {
    ignore_changes = [
      backup_strategy,
    ]
  }
}

resource "flexibleengine_rds_backup" "test" {
  name        = "%[2]s"
  instance_id = flexibleengine_rds_instance_v3.test.id
}

`, testVpc(name), name)
}

func testAccDatasourceBackup_basic(name string) string {
	return fmt.Sprintf(`
%s 

data "flexibleengine_rds_backups" "test" {
  instance_id = flexibleengine_rds_instance_v3.test.id
  backup_type = "manual"

  depends_on = [
    flexibleengine_rds_backup.test
  ]
}
`, testAccDatasourceBackup_base(name))
}

func testAccDatasourceBackup_auto_basic(name string) string {
	return fmt.Sprintf(`
%s 

data "flexibleengine_rds_backups" "test" {
  instance_id = flexibleengine_rds_instance_v3.test.id
  backup_type = "auto"

  depends_on = [
    flexibleengine_rds_backup.test
  ]
}
`, testAccDatasourceBackup_base(name))
}

func testAccDatasourceBackup_incremental_basic(name string) string {
	return fmt.Sprintf(`
%s 

data "flexibleengine_rds_backups" "test" {
  instance_id = flexibleengine_rds_instance_v3.test.id
  backup_type = "incremental"

  depends_on = [
    flexibleengine_rds_backup.test
  ]
}
`, testAccDatasourceBackup_base(name))
}
