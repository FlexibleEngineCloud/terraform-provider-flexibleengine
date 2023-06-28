package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/cce/v3/nodes"
)

func TestAccCCENodeV3_basic(t *testing.T) {
	var node nodes.Nodes

	cceName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "flexibleengine_cce_node_v3.node_1"
	clusterName := "flexibleengine_cce_cluster_v3.cluster_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccCCEKeyPairPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCENodeV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodeV3_basic(cceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodeV3Exists(resourceName, clusterName, &node),
					resource.TestCheckResourceAttr(resourceName, "name", cceName),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", "s3.large.2"),
					resource.TestCheckResourceAttr(resourceName, "status", "Active"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
				),
			},
			{
				Config: testAccCCENodeV3_update(cceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", cceName+"-update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCCENodeImportStateIdFunc(),
			},
		},
	})
}

func TestAccCCENodeV3_volumes_encryption(t *testing.T) {
	var node nodes.Nodes

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "flexibleengine_cce_node_v3.node_1"
	clusterName := "flexibleengine_cce_cluster_v3.cluster_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccCCEKeyPairPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCENodeV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodeV3_volumes_encryption(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodeV3Exists(resourceName, clusterName, &node),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrPair(resourceName, "root_volume.0.kms_key_id",
						"flexibleengine_kms_key_v1.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "data_volumes.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "data_volumes.0.kms_key_id",
						"flexibleengine_kms_key_v1.test", "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCCENodeImportStateIdFunc(),
			},
		},
	})
}

func testAccCCENodeImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		cluster, ok := s.RootModule().Resources["flexibleengine_cce_cluster_v3.cluster_1"]
		if !ok {
			return "", fmt.Errorf("Cluster not found: %s", cluster)
		}
		node, ok := s.RootModule().Resources["flexibleengine_cce_node_v3.node_1"]
		if !ok {
			return "", fmt.Errorf("Node not found: %s", node)
		}

		if cluster.Primary.ID == "" || node.Primary.ID == "" {
			return "", fmt.Errorf("resource not found: %s/%s", cluster.Primary.ID, node.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", cluster.Primary.ID, node.Primary.ID), nil
	}
}

func testAccCheckCCENodeV3Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	cceClient, err := config.CceV3Client(OS_REGION_NAME)
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
		cceClient, err := config.CceV3Client(OS_REGION_NAME)
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

func testAccCCENodeV3_base(rName string) string {
	return fmt.Sprintf(`
data "flexibleengine_availability_zones" "test" {}

resource "flexibleengine_cce_cluster_v3" "cluster_1" {
  name         = "%s"
  cluster_type = "VirtualMachine"
  flavor_id    = "cce.s1.small"
  vpc_id       = "%s"
  subnet_id    = "%s"
  container_network_type = "overlay_l2"
}`, rName, OS_VPC_ID, OS_NETWORK_ID)
}

func testAccCCENodeV3_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_cce_node_v3" "node_1" {
  cluster_id        = flexibleengine_cce_cluster_v3.cluster_1.id
  name              = "%s"
  flavor_id         = "s3.large.2"
  availability_zone = data.flexibleengine_availability_zones.test.names[0]
  key_pair          = "%s"
  subnet_id         = "%s"

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
}`, testAccCCENodeV3_base(rName), rName, OS_KEYPAIR_NAME, OS_NETWORK_ID)
}

func testAccCCENodeV3_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_cce_node_v3" "node_1" {
  cluster_id        = flexibleengine_cce_cluster_v3.cluster_1.id
  name              = "%s-update"
  flavor_id         = "s3.large.2"
  availability_zone = data.flexibleengine_availability_zones.test.names[0]
  key_pair          = "%s"
  subnet_id         = "%s"

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
}`, testAccCCENodeV3_base(rName), rName, OS_KEYPAIR_NAME, OS_NETWORK_ID)
}

func testAccCCENodeV3_volumes_encryption(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_kms_key_v1" "test" {
  key_alias    = "%s"
  pending_days = "7"
}

resource "flexibleengine_cce_node_v3" "node_1" {
  cluster_id        = flexibleengine_cce_cluster_v3.cluster_1.id
  name              = "%s"
  flavor_id         = "s3.large.2"
  availability_zone = data.flexibleengine_availability_zones.test.names[0]
  key_pair          = "%s"

  root_volume {
    size       = 40
    volumetype = "SSD"
    kms_key_id = flexibleengine_kms_key_v1.test.id
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
    kms_key_id = flexibleengine_kms_key_v1.test.id
  }
  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccCCENodeV3_base(rName), rName, rName, OS_KEYPAIR_NAME)
}
