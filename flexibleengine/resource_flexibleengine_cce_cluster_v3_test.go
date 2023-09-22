package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/cce/v3/clusters"
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
					resource.TestCheckResourceAttr(resourceName, "cluster_version", "v1.17.9-r0"),
					resource.TestCheckResourceAttr(resourceName, "container_network_type", "overlay_l2"),
					resource.TestCheckResourceAttr(resourceName, "authentication_mode", "rbac"),
					resource.TestCheckResourceAttrSet(resourceName, "internal_endpoint"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"eip",
				},
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
	cceClient, err := config.CceV3Client(OS_REGION_NAME)
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
		cceClient, err := config.CceV3Client(OS_REGION_NAME)
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

func testAccCCEClusterV3_Base(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_vpc_v1" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "test" {
  name       = "%s"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"

  //dns is required for cce node installing
  primary_dns   = "100.125.0.41"
  secondary_dns = "100.126.0.41"
  vpc_id        = flexibleengine_vpc_v1.test.id
}
`, rName, rName)
}

func testAccCCEClusterV3_basic(cceName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name       = "test"
    size       = 10
    share_type = "PER"
  }
}

resource "flexibleengine_cce_cluster_v3" "cluster_1" {
  name                   = "%s"
  description            = "a description"
  cluster_type           = "VirtualMachine"
  cluster_version        = "v1.17.9-r0"
  flavor_id              = "cce.s1.small"
  vpc_id                 = flexibleengine_vpc_v1.test.id
  subnet_id              = flexibleengine_vpc_subnet_v1.test.id
  container_network_type = "overlay_l2"
  eip                    = flexibleengine_vpc_eip.test.address
}`, testAccCCEClusterV3_Base(cceName), cceName)
}

func testAccCCEClusterV3_update(cceName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name       = "test"
    size       = 10
    share_type = "PER"
  }
}

resource "flexibleengine_vpc_eip" "update" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name       = "test"
    size       = 10
    share_type = "PER"
  }
}

resource "flexibleengine_cce_cluster_v3" "cluster_1" {
  name                   = "%s"
  description            = "a updated description"
  cluster_type           = "VirtualMachine"
  cluster_version        = "v1.17.9-r0"
  flavor_id              = "cce.s1.small"
  vpc_id                 = flexibleengine_vpc_v1.test.id
  subnet_id              = flexibleengine_vpc_subnet_v1.test.id
  container_network_type = "overlay_l2"
  eip                    = flexibleengine_vpc_eip.update.address
}`, testAccCCEClusterV3_Base(cceName), cceName)
}

func TestAccCluster_hibernate(t *testing.T) {
	var cluster clusters.Clusters

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "flexibleengine_cce_cluster_v3.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCEClusterV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCluster_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEClusterV3Exists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "Available"),
				),
			},
			{
				Config: testAccCluster_hibernate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEClusterV3Exists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "Hibernation"),
				),
			},
			{
				Config: testAccCluster_awake(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEClusterV3Exists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "Available"),
				),
			},
		},
	})
}

func testAccCluster_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_cce_cluster_v3" "test" {
  name                   = "%s"
  flavor_id              = "cce.s1.small"
  cluster_type           = "VirtualMachine"
  vpc_id                 = flexibleengine_vpc_v1.test.id
  subnet_id              = flexibleengine_vpc_subnet_v1.test.id
  container_network_type = "overlay_l2"
  service_network_cidr   = "10.248.0.0/16"
}
`, testAccCCEClusterV3_Base(rName), rName)
}

func testAccCluster_hibernate(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_cce_cluster_v3" "test" {
  name                   = "%s"
  flavor_id              = "cce.s1.small"
  cluster_type           = "VirtualMachine"
  vpc_id                 = flexibleengine_vpc_v1.test.id
  subnet_id              = flexibleengine_vpc_subnet_v1.test.id
  container_network_type = "overlay_l2"
  service_network_cidr   = "10.248.0.0/16"
  hibernate              = true
}
`, testAccCCEClusterV3_Base(rName), rName)
}

func testAccCluster_awake(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_cce_cluster_v3" "test" {
  name                   = "%s"
  flavor_id              = "cce.s1.small"
  cluster_type           = "VirtualMachine"
  vpc_id                 = flexibleengine_vpc_v1.test.id
  subnet_id              = flexibleengine_vpc_subnet_v1.test.id
  container_network_type = "overlay_l2"
  service_network_cidr   = "10.248.0.0/16"
  hibernate              = false
}
`, testAccCCEClusterV3_Base(rName), rName)
}
