package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVpcV1EipDataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	dataSourceName := "data.flexibleengine_vpc_eip_v1.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpcEipConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcEipDataSourceID(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "status", "UNBOUND"),
					resource.TestCheckResourceAttr(dataSourceName, "type", "5_bgp"),
					resource.TestCheckResourceAttr(dataSourceName, "bandwidth_size", "8"),
					resource.TestCheckResourceAttr(dataSourceName, "bandwidth_share_type", "PER"),
				),
			},
		},
	})
}

func testAccCheckVpcEipDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find eip data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Eip data source ID not set")
		}

		return nil
	}
}

func testAccDataSourceVpcEipConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_vpc_eip_v1" "test" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "%s"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

data "flexibleengine_vpc_eip_v1" "test" {
  public_ip = flexibleengine_vpc_eip_v1.test.address
}
`, rName)
}
