package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataBackup_basic(t *testing.T) {
	randName := acceptance.RandomAccResourceNameWithDash()
	dataSourceName := "data.flexibleengine_cbr_backup.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataBackup_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
				),
			},
		},
	})
}

func testAccDataBackup_base(name string) string {
	return fmt.Sprintf(`
resource "flexibleengine_vpc_v1" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "test" {
  name       = "%[1]s"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = flexibleengine_vpc_v1.test.id

  tags = {
    foo = "bar"
    key = "value"
  }
}

resource "flexibleengine_networking_secgroup_v2" "test" {
  name        = "%[1]s"
  description = "terraform security group acceptance test"
}

data "flexibleengine_availability_zones" "test" {}

data "flexibleengine_compute_flavors_v2" "test" {
  availability_zone = data.flexibleengine_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core          = 2
  memory_size       = 4
}

data "flexibleengine_images_image" "test" {
  name = "OBS Ubuntu 18.04"
}
`, name)
}

func testAccDataBackup_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "flexibleengine_compute_instance_v2" "test" {
  name              = "%[2]s"
  image_id          = data.flexibleengine_images_image.test.id
  flavor_id         = data.flexibleengine_compute_flavors_v2.test.flavors[0]
  security_groups   = [flexibleengine_networking_secgroup_v2.test.name]
  availability_zone = data.flexibleengine_availability_zones.test.names[0]

  network {
    uuid = flexibleengine_vpc_subnet_v1.test.id
  }
}

resource "flexibleengine_cbr_vault" "test" {
  name             = "%[2]s"
  type             = "server"
  consistent_level = "app_consistent"
  protection_type  = "backup"
  size             = 200
}

resource "flexibleengine_images_image" "test" {
  name        = "%[2]s"
  instance_id = flexibleengine_compute_instance_v2.test.id
  vault_id    = flexibleengine_cbr_vault.test.id
}

data "flexibleengine_cbr_backup" "test" {
  id = flexibleengine_images_image.test.backup_id
}
`, testAccDataBackup_base(name), name)
}
