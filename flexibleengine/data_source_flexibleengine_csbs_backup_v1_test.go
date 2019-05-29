package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCSBSBackupV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCSBSBackupV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCSBSBackupV1DataSourceID("data.flexibleengine_csbs_backup_v1.csbs"),
					resource.TestCheckResourceAttr("data.flexibleengine_csbs_backup_v1.csbs", "backup_name", "csbs-test"),
					resource.TestCheckResourceAttr("data.flexibleengine_csbs_backup_v1.csbs", "resource_name", "instance_1"),
				),
			},
		},
	})
}

func testAccCheckCSBSBackupV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find backup data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("backup data source ID not set ")
		}

		return nil
	}
}

var testAccCSBSBackupV1DataSource_basic = fmt.Sprintf(`
resource "flexibleengine_compute_instance_v2" "instance_1" {
  name = "instance_1"
  image_id = "%s"
  security_groups = ["default"]
  availability_zone = "%s"
  flavor_id = "%s"
  metadata = {
    foo = "bar"
  }
  network {
    uuid = "%s"
  }
}
resource "flexibleengine_csbs_backup_v1" "csbs" {
  backup_name      = "csbs-test"
  description      = "test-code"
  resource_id = "${flexibleengine_compute_instance_v2.instance_1.id}"
  resource_type = "OS::Nova::Server"
}
data "flexibleengine_csbs_backup_v1" "csbs" {
  id = "${flexibleengine_csbs_backup_v1.csbs.id}"
}
`, OS_IMAGE_ID, OS_AVAILABILITY_ZONE, OS_FLAVOR_ID, OS_NETWORK_ID)
