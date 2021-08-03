package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/golangsdk/openstack/cce/v3/nodes"
)

func TestAccCCENodeV3_basic(t *testing.T) {
	var node nodes.Nodes
	var cceName = fmt.Sprintf("terra-test-%s", acctest.RandString(5))
	resourceName := "flexibleengine_cce_node_v3.node_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccCCEKeyPairPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCENodeV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodeV3_basic(cceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodeV3Exists("flexibleengine_cce_node_v3.node_1", "flexibleengine_cce_cluster_v3.cluster_1", &node),
					resource.TestCheckResourceAttr(resourceName, "name", "test-node"),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", "s1.medium"),
					resource.TestCheckResourceAttr(resourceName, "status", "Active"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
				),
			},
			{
				Config: testAccCCENodeV3_update(cceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-node2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
		},
	})
}

func testAccCheckCCENodeV3Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	cceClient, err := config.cceV3Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating flexibleengine CCE client: %s", err)
	}

	var clusterID string
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "flexibleengine_cce_cluster_v3" {
			clusterID = rs.Primary.ID
			break
		}
	}
	if clusterID == "" {
		return fmt.Errorf("cce cluster not found")
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_cce_node_v3" {
			continue
		}

		_, err := nodes.Get(cceClient, clusterID, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Node still exists")
		}
	}

	return nil
}

func testAccCheckCCENodeV3Exists(n string, cluster string, node *nodes.Nodes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		c, ok := s.RootModule().Resources[cluster]
		if !ok {
			return fmt.Errorf("Cluster not found: %s", c)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		if c.Primary.ID == "" {
			return fmt.Errorf("Cluster id is not set")
		}

		config := testAccProvider.Meta().(*Config)
		cceClient, err := config.cceV3Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine CCE client: %s", err)
		}

		found, err := nodes.Get(cceClient, c.Primary.ID, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Metadata.Id != rs.Primary.ID {
			return fmt.Errorf("Node not found")
		}

		*node = *found

		return nil
	}
}

func testAccCCENodeV3_basic(cceName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_cce_cluster_v3" "cluster_1" {
  name         = "%s"
  cluster_type = "VirtualMachine"
  flavor_id    = "cce.s1.small"
  vpc_id       = "%s"
  subnet_id    = "%s"
  container_network_type = "overlay_l2"
}

resource "flexibleengine_cce_node_v3" "node_1" {
  cluster_id        = flexibleengine_cce_cluster_v3.cluster_1.id
  name              = "test-node"
  flavor_id         = "s1.medium"
  availability_zone = "%s"
  key_pair          = "%s"

  root_volume {
    size       = 40
    volumetype = "SATA"
  }
  data_volumes {
    size       = 100
    volumetype = "SATA"
  }
  tags = {
    key = "value"
    foo = "bar"
  }
}`, cceName, OS_VPC_ID, OS_NETWORK_ID, OS_AVAILABILITY_ZONE, OS_KEYPAIR_NAME)
}

func testAccCCENodeV3_update(cceName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_cce_cluster_v3" "cluster_1" {
  name         = "%s"
  cluster_type = "VirtualMachine"
  flavor_id    = "cce.s1.small"
  vpc_id       = "%s"
  subnet_id    = "%s"
  container_network_type = "overlay_l2"
}

resource "flexibleengine_cce_node_v3" "node_1" {
  cluster_id        = flexibleengine_cce_cluster_v3.cluster_1.id
  name              = "test-node2"
  flavor_id         = "s1.medium"
  availability_zone = "%s"
  key_pair          = "%s"

  root_volume {
    size       = 40
    volumetype = "SATA"
  }
  data_volumes {
    size       = 100
    volumetype = "SATA"
  }
  tags = {
    key   = "value1"
    owner = "terraform"
  }
}`, cceName, OS_VPC_ID, OS_NETWORK_ID, OS_AVAILABILITY_ZONE, OS_KEYPAIR_NAME)
}
