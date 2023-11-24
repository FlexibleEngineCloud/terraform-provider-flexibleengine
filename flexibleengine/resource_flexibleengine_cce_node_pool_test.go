package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/cce/v3/nodepools"
)

func TestAccCCENodePool_basic(t *testing.T) {
	var nodePool nodepools.NodePool

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	updateName := rName + "update"
	resourceName := "flexibleengine_cce_node_pool_v3.test"
	//clusterName here is used to provide the cluster id to fetch cce node pool.
	clusterName := "flexibleengine_cce_cluster_v3.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCENodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodePool_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodePoolExists(resourceName, clusterName, &nodePool),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "scale_enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "initial_node_count", "1"),
					resource.TestCheckResourceAttr(resourceName, "current_node_count", "1"),
					resource.TestCheckResourceAttr(resourceName, "min_node_count", "0"),
					resource.TestCheckResourceAttr(resourceName, "max_node_count", "0"),
					resource.TestCheckResourceAttr(resourceName, "max_pods", "200"),
					resource.TestCheckResourceAttr(resourceName, "taints.0.key", "bar"),
					resource.TestCheckResourceAttr(resourceName, "taints.0.value", "foo"),
					resource.TestCheckResourceAttr(resourceName, "taints.0.effect", "NoSchedule"),
					resource.TestCheckResourceAttr(resourceName, "labels.pool", "acc-test-pool"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
				),
			},
			{
				Config: testAccCCENodePool_update(rName, updateName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "scale_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "initial_node_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "current_node_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "min_node_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "max_node_count", "9"),
					resource.TestCheckResourceAttr(resourceName, "scale_down_cooldown_time", "100"),
					resource.TestCheckResourceAttr(resourceName, "priority", "1"),
					resource.TestCheckResourceAttr(resourceName, "taints.0.key", "looks"),
					resource.TestCheckResourceAttr(resourceName, "taints.0.value", "bad"),
					resource.TestCheckResourceAttr(resourceName, "taints.0.effect", "NoExecute"),
					resource.TestCheckResourceAttr(resourceName, "labels.pool", "acc-test-pool-update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateIdFunc:       nodePoolImportStateIdFunc(resourceName),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_node_count", "max_pods"},
			},
		},
	})
}

func TestAccCCENodePool_serverGroup(t *testing.T) {
	var nodePool nodepools.NodePool

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "flexibleengine_cce_node_pool_v3.test"
	// clusterName here is used to provide the cluster id to fetch cce node pool.
	clusterName := "flexibleengine_cce_cluster_v3.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCENodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodePool_serverGroup(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodePoolExists(resourceName, clusterName, &nodePool),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrPair(resourceName, "ecs_group_id",
						"flexibleengine_compute_servergroup_v2.test", "id"),
				),
			},
		},
	})
}

func TestAccCCENodePool_securityGroups(t *testing.T) {
	var nodePool nodepools.NodePool

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "flexibleengine_cce_node_pool_v3.test"
	//clusterName here is used to provide the cluster id to fetch cce node pool.
	clusterName := "flexibleengine_cce_cluster_v3.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCENodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNodePool_securityGroups(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodePoolExists(resourceName, clusterName, &nodePool),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "runtime", "docker"),
					resource.TestCheckResourceAttrSet(resourceName, "security_groups.#"),
				),
			},
		},
	})
}

func nodePoolImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		clusterID := rs.Primary.Attributes["cluster_id"]
		return fmt.Sprintf("%s/%s", clusterID, rs.Primary.Attributes["id"]), nil
	}
}

func testAccCheckCCENodePoolDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	cceClient, err := config.CceV3Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating Flexibleengine CCE client: %s", err)
	}

	var clusterId string
	var nodepollId string

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "flexibleengine_cce_cluster_v3" {
			clusterId = rs.Primary.ID
		}

		if rs.Type == "flexibleengine_cce_node_pool_v3" {
			nodepollId = rs.Primary.ID
		}

		if clusterId == "" || nodepollId == "" {
			continue
		}

		_, err := nodepools.Get(cceClient, clusterId, nodepollId).Extract()
		if err == nil {
			return fmt.Errorf("Node still exists")
		}
	}

	return nil
}

func testAccCheckCCENodePoolExists(n string, cluster string, nodePool *nodepools.NodePool) resource.TestCheckFunc {
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
			return fmt.Errorf("Error creating Flexibleengine CCE client: %s", err)
		}

		found, err := nodepools.Get(cceClient, c.Primary.ID, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Metadata.Id != rs.Primary.ID {
			return fmt.Errorf("Node Pool not found")
		}

		*nodePool = *found

		return nil
	}
}

func testAccCCENodePool_Base(rName string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_availability_zones" "test" {}

resource "flexibleengine_compute_keypair_v2" "test" {
  name = "%s"
}

resource "flexibleengine_cce_cluster_v3" "test" {
  name                   = "%s"
  description            = "a description"
  cluster_type           = "VirtualMachine"
  flavor_id              = "cce.s1.small"
  vpc_id                 = flexibleengine_vpc_v1.test.id
  subnet_id              = flexibleengine_vpc_subnet_v1.test.id
  container_network_type = "overlay_l2"
}
`, testAccCCEClusterV3_Base(rName), rName, rName)
}

func testAccCCENodePool_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_cce_node_pool_v3" "test" {
  cluster_id               = flexibleengine_cce_cluster_v3.test.id
  name                     = "%s"
  os                       = "EulerOS 2.5"
  flavor_id                = "s3.large.2"
  availability_zone        = data.flexibleengine_availability_zones.test.names[0]
  key_pair                 = flexibleengine_compute_keypair_v2.test.name
  scale_enable             = false
  initial_node_count       = 1
  min_node_count           = 0
  max_node_count           = 0
  max_pods                 = 200
  scale_down_cooldown_time = 0
  priority                 = 0
  type                     = "vm"

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }

  taints {
    key    = "bar"
    value  = "foo"
    effect = "NoSchedule"
  }
  labels = {
    pool = "acc-test-pool"
  }
  tags = {
    key = "value"
    foo = "bar"
  }
}
`, testAccCCENodePool_Base(rName), rName)
}

func testAccCCENodePool_update(rName, updateName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_cce_node_pool_v3" "test" {
  cluster_id               = flexibleengine_cce_cluster_v3.test.id
  name                     = "%s"
  os                       = "EulerOS 2.5"
  flavor_id                = "s3.large.2"
  availability_zone        = data.flexibleengine_availability_zones.test.names[0]
  key_pair                 = flexibleengine_compute_keypair_v2.test.name
  scale_enable             = true
  initial_node_count       = 2
  min_node_count           = 2
  max_node_count           = 9
  max_pods                 = 200
  scale_down_cooldown_time = 100
  priority                 = 1
  type                     = "vm"

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }

  taints {
    key    = "looks"
    value  = "bad"
    effect = "NoExecute"
  }
  labels = {
    pool = "acc-test-pool-update"
  }
  tags = {
    key   = "value1"
    owner = "terraform"
  }
}
`, testAccCCENodePool_Base(rName), updateName)
}

func testAccCCENodePool_serverGroup(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_compute_servergroup_v2" "test" {
  name     = "%[2]s"
  policies = ["anti-affinity"]
}

resource "flexibleengine_cce_node_pool_v3" "test" {
  cluster_id               = flexibleengine_cce_cluster_v3.test.id
  name                     = "%[2]s"
  os                       = "EulerOS 2.5"
  flavor_id                = "s3.large.2"
  initial_node_count       = 1
  availability_zone        = data.flexibleengine_availability_zones.test.names[0]
  key_pair                 = flexibleengine_compute_keypair_v2.test.name
  scall_enable             = false
  min_node_count           = 0
  max_node_count           = 0
  scale_down_cooldown_time = 0
  priority                 = 0
  type                     = "vm"
  ecs_group_id             = flexibleengine_compute_servergroup_v2.test.id

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }
}
`, testAccCCENodePool_Base(rName), rName)
}

func testAccNodePool_securityGroups(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_networking_secgroup_v2" "test1" {
  name        = "%[2]s-test1"
  description = "terraform security group"
}

resource "flexibleengine_cce_node_pool_v3" "test" {
  cluster_id               = flexibleengine_cce_cluster_v3.test.id
  name                     = "%[2]s"
  os                       = "EulerOS 2.5"
  runtime                  = "docker"
  flavor_id                = "s3.large.2"
  availability_zone        = data.flexibleengine_availability_zones.test.names[0]
  key_pair                 = flexibleengine_compute_keypair_v2.test.name
  scale_enable             = false
  initial_node_count       = 1
  min_node_count           = 0
  max_node_count           = 0
  max_pods                 = 200
  scale_down_cooldown_time = 0
  priority                 = 0
  type                     = "vm"

  security_groups = [
    flexibleengine_networking_secgroup_v2.test1.id
  ]

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }
}
`, testAccCCENodePool_Base(rName), rName)
}
