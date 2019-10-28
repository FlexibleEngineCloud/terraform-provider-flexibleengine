package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/cce/v3/clusters"
)

func TestAccCCEClusterV3_basic(t *testing.T) {
	var cluster clusters.Clusters
	var cceName = fmt.Sprintf("terra-test-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCEClusterV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEClusterV3_basic(cceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEClusterV3Exists("flexibleengine_cce_cluster_v3.cluster_1", &cluster),
					resource.TestCheckResourceAttr(
						"flexibleengine_cce_cluster_v3.cluster_1", "name", cceName),
					resource.TestCheckResourceAttr(
						"flexibleengine_cce_cluster_v3.cluster_1", "status", "Available"),
					resource.TestCheckResourceAttr(
						"flexibleengine_cce_cluster_v3.cluster_1", "cluster_type", "VirtualMachine"),
					resource.TestCheckResourceAttr(
						"flexibleengine_cce_cluster_v3.cluster_1", "flavor_id", "cce.s1.small"),
					resource.TestCheckResourceAttr(
						"flexibleengine_cce_cluster_v3.cluster_1", "cluster_version", "v1.11.7"),
					resource.TestCheckResourceAttr(
						"flexibleengine_cce_cluster_v3.cluster_1", "container_network_type", "overlay_l2"),
					resource.TestCheckResourceAttr(
						"flexibleengine_cce_cluster_v3.cluster_1", "authentication_mode", "x509"),
				),
			},
			{
				Config: testAccCCEClusterV3_update(cceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"flexibleengine_cce_cluster_v3.cluster_1", "description", "new description"),
				),
			},
		},
	})
}

func TestAccCCEClusterV3_timeout(t *testing.T) {
	var cluster clusters.Clusters
	var cceName = fmt.Sprintf("terra-test-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCEClusterV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEClusterV3_timeout(cceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEClusterV3Exists("flexibleengine_cce_cluster_v3.cluster_1", &cluster),
					resource.TestCheckResourceAttr(
						"flexibleengine_cce_cluster_v3.cluster_1", "authentication_mode", "rbac"),
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
  name = "%s"
  cluster_type="VirtualMachine"
  flavor_id="cce.s1.small"
  cluster_version = "v1.11.7"
  vpc_id="%s"
  subnet_id="%s"
  container_network_type="overlay_l2"
}`, cceName, OS_VPC_ID, OS_NETWORK_ID)
}

func testAccCCEClusterV3_update(cceName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_cce_cluster_v3" "cluster_1" {
  name = "%s"
  cluster_type="VirtualMachine"
  flavor_id="cce.s1.small"
  cluster_version = "v1.11.7"
  vpc_id="%s"
  subnet_id="%s"
  container_network_type="overlay_l2"
  description="new description"
}`, cceName, OS_VPC_ID, OS_NETWORK_ID)
}

func testAccCCEClusterV3_timeout(cceName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_networking_floatingip_v2" "fip_1" {
}

resource "flexibleengine_cce_cluster_v3" "cluster_1" {
  name = "%s"
  cluster_type="VirtualMachine"
  flavor_id="cce.s1.small"
  vpc_id="%s"
  subnet_id="%s"
  eip="${flexibleengine_networking_floatingip_v2.fip_1.address}"
  authentication_mode = "rbac"
  container_network_type="overlay_l2"
    timeouts {
    create = "10m"
    delete = "10m"
  }
}
`, cceName, OS_VPC_ID, OS_NETWORK_ID)
}
