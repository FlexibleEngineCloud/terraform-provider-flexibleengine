package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCCENodesV3DataSource_basic(t *testing.T) {
	var cceName = fmt.Sprintf("terra-test-%s", acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccCCEKeyPairPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodeV3DataSource_basic(cceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodeV3DataSourceID("data.flexibleengine_cce_node_v3.nodes"),
					resource.TestCheckResourceAttr("data.flexibleengine_cce_node_v3.nodes", "name", cceName),
					resource.TestCheckResourceAttr("data.flexibleengine_cce_node_v3.nodes", "flavor_id", "s1.medium"),
				),
			},
		},
	})
}

func testAccCheckCCENodeV3DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find nodes data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Node data source ID not set ")
		}

		return nil
	}
}

func testAccCCENodeV3DataSource_basic(cceName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_cce_cluster_v3" "cluster_1" {
  name = "%s"
  cluster_type="VirtualMachine"
  flavor_id="cce.s1.small"
  cluster_version = "v1.9.7-r1"
  vpc_id="%s"
  subnet_id="%s"
  container_network_type="overlay_l2"
}

resource "flexibleengine_cce_node_v3" "node_1" {
cluster_id = "${flexibleengine_cce_cluster_v3.cluster_1.id}"
  name = "%s"
  flavor_id="s1.medium"
  availability_zone= "%s"
  key_pair="%s"
  root_volume = {
    size= 40,
    volumetype= "SATA"
  }
  data_volumes = [
    {
      size= 100,
      volumetype= "SATA"
    },
  ]
}
data "flexibleengine_cce_node_v3" "nodes" {
	cluster_id = "${flexibleengine_cce_cluster_v3.cluster_1.id}"
	name = "${flexibleengine_cce_node_v3.node_1.name}"
}
`, cceName, OS_VPC_ID, OS_NETWORK_ID, cceName, OS_AVAILABILITY_ZONE, OS_KEYPAIR_NAME)
}
