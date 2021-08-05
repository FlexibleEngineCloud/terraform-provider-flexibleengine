package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCCEClusterV3DataSource_basic(t *testing.T) {
	var cceName = fmt.Sprintf("terra-test-%s", acctest.RandString(5))
	resourceName := "data.flexibleengine_cce_cluster_v3.clusters"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEClusterV3DataSource_basic(cceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEClusterV3DataSourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", cceName),
					resource.TestCheckResourceAttr(resourceName, "status", "Available"),
					resource.TestCheckResourceAttr(resourceName, "cluster_type", "VirtualMachine"),
				),
			},
		},
	})
}

func testAccCheckCCEClusterV3DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find cluster data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("cluster data source ID not set ")
		}

		return nil
	}
}

func testAccCCEClusterV3DataSource_basic(cceName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_cce_cluster_v3" "cluster_1" {
  name                   = "%s"
  cluster_type           = "VirtualMachine"
  flavor_id              = "cce.s1.small"
  vpc_id                 = "%s"
  subnet_id              = "%s"
  container_network_type = "overlay_l2"
}

data "flexibleengine_cce_cluster_v3" "clusters" {
  name = flexibleengine_cce_cluster_v3.cluster_1.name
}`, cceName, OS_VPC_ID, OS_NETWORK_ID)
}
