package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccBMSNicV2DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBMSNicV2DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBMSNicV2DataSourceID("data.flexibleengine_compute_bms_nic_v2.nic_1"),
					resource.TestCheckResourceAttr("data.flexibleengine_compute_bms_nic_v2.nic_1", "status", "ACTIVE"),
				),
			},
		},
	})
}

func testAccCheckBMSNicV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find nic data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("nic data source ID not set ")
		}

		return nil
	}
}

var testAccBMSNicV2DataSource_basic = fmt.Sprintf(`
resource "flexibleengine_compute_instance_v2" "instance_1" {
  name = "BMSinstance_1"
  image_id = "%s"
  security_groups = ["default"]
  availability_zone = "%s"
  flavor_id = "physical.o2.medium"
  flavor_name = "physical.o2.medium"
  metadata = {
    foo = "bar"
  }
  network {
    uuid = "%s"
  }
}
data "flexibleengine_compute_bms_nic_v2" "nic_1" {
  server_id = "${flexibleengine_compute_instance_v2.instance_1.id}"
}
`, OS_IMAGE_ID, OS_AVAILABILITY_ZONE, OS_NETWORK_ID)
