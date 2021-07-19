package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccBlockStorageVolumeV2DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBlockStorageVolumeV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceBlockStorageVolumeV2Check(
						"data.flexibleengine_blockstorage_volume_v2.volume_ds"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_blockstorage_volume_v2.volume_ds", "name", "volume_ds"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_blockstorage_volume_v2.volume_ds", "status", "available"),
				),
			},
		},
	})
}

func testAccDataSourceBlockStorageVolumeV2Check(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find volume data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Volume data source ID not set")
		}

		return nil
	}
}

const testAccDataSourceBlockStorageVolumeV2Config = `
resource "flexibleengine_blockstorage_volume_v2" "volume_1" {
  name = "volume_ds"
  description = "first test volume"
  size = 1
}

data "flexibleengine_blockstorage_volume_v2" "volume_ds" {
  name = "${flexibleengine_blockstorage_volume_v2.volume_1.name}"
}
`
