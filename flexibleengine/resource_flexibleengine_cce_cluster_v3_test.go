package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/golangsdk/openstack/cce/v3/clusters"
)

func TestAccCCEClusterV3_basic(t *testing.T) {
	var cluster clusters.Clusters
	var cceName = fmt.Sprintf("terra-test-%s", acctest.RandString(5))
	resourceName := "flexibleengine_cce_cluster_v3.cluster_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCEClusterV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEClusterV3_basic(cceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEClusterV3Exists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", cceName),
					resource.TestCheckResourceAttr(resourceName, "description", "a description"),
					resource.TestCheckResourceAttr(resourceName, "status", "Available"),
					resource.TestCheckResourceAttr(resourceName, "cluster_type", "VirtualMachine"),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", "cce.s1.small"),
					resource.TestCheckResourceAttr(resourceName, "cluster_version", "v1.11.7-r2"),
					resource.TestCheckResourceAttr(resourceName, "container_network_type", "overlay_l2"),
					resource.TestCheckResourceAttr(resourceName, "authentication_mode", "rbac"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCCEClusterV3_update(cceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "a updated description"),
				),
			},
		},
	})
}

func testAccCheckCCEClusterV3Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	cceClient, err := config.cceV3Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating flexibleengine CCE client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_cce_cluster_v3" {
			continue
		}

		_, err := clusters.Get(cceClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Cluster still exists")
		}
	}

	return nil
}

func testAccCheckCCEClusterV3Exists(n string, cluster *clusters.Clusters) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		cceClient, err := config.cceV3Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating flexibleengine CCE client: %s", err)
		}

		found, err := clusters.Get(cceClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Metadata.Id != rs.Primary.ID {
			return fmt.Errorf("Cluster not found")
		}

		*cluster = *found

		return nil
	}
}

func testAccCCEClusterV3_basic(cceName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_cce_cluster_v3" "cluster_1" {
  name            = "%s"
  description     = "a description"
  cluster_type    = "VirtualMachine"
  cluster_version = "v1.11.7-r2"
  flavor_id       = "cce.s1.small"
  vpc_id          = "%s"
  subnet_id       = "%s"
  container_network_type = "overlay_l2"
}`, cceName, OS_VPC_ID, OS_NETWORK_ID)
}

func testAccCCEClusterV3_update(cceName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_cce_cluster_v3" "cluster_1" {
  name            = "%s"
  description     = "a updated description"
  cluster_type    = "VirtualMachine"
  cluster_version = "v1.11.7-r2"
  flavor_id       = "cce.s1.small"
  vpc_id          = "%s"
  subnet_id       = "%s"
  container_network_type = "overlay_l2"
}`, cceName, OS_VPC_ID, OS_NETWORK_ID)
}
