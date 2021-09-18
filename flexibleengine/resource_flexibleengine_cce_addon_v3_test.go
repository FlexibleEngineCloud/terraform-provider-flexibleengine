package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/cce/v3/addons"
)

func TestAccCCEAddon_basic(t *testing.T) {
	var addon addons.Addon

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "flexibleengine_cce_addon_v3.test"
	clusterName := "flexibleengine_cce_cluster_v3.cluster_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCEAddonDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEAddon_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEAddonExists(resourceName, clusterName, &addon),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCCEAddonImportStateIdFunc(),
			},
		},
	})
}

func testAccCheckCCEAddonDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	cceClient, err := config.CceAddonV3Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine CCE Addon client: %s", err)
	}

	var clusterID string
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "flexibleengine_cce_cluster_v3" {
			clusterID = rs.Primary.ID
			break
		}
	}

	if clusterID == "" {
		return nil
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_cce_addon_v3" {
			continue
		}

		_, err := addons.Get(cceClient, rs.Primary.ID, clusterID).Extract()
		if err == nil {
			return fmt.Errorf("addon still exists")
		}
	}
	return nil
}

func testAccCheckCCEAddonExists(n string, cluster string, addon *addons.Addon) resource.TestCheckFunc {
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
		cceClient, err := config.CceAddonV3Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine CCE Addon client: %s", err)
		}

		found, err := addons.Get(cceClient, rs.Primary.ID, c.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Metadata.Id != rs.Primary.ID {
			return fmt.Errorf("Addon not found")
		}

		*addon = *found

		return nil
	}
}

func testAccCCEAddonImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var clusterID string
		var addonID string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "flexibleengine_cce_cluster_v3" {
				clusterID = rs.Primary.ID
			} else if rs.Type == "flexibleengine_cce_addon_v3" {
				addonID = rs.Primary.ID
			}
		}
		if clusterID == "" || addonID == "" {
			return "", fmt.Errorf("resource not found: %s/%s", clusterID, addonID)
		}
		return fmt.Sprintf("%s/%s", clusterID, addonID), nil
	}
}

func testAccCCEAddon_basic(rName string) string {
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
}

resource "flexibleengine_cce_addon_v3" "test" {
  cluster_id    = flexibleengine_cce_cluster_v3.cluster_1.id
  version       = "1.0.6"
  template_name = "metrics-server"
  depends_on    = [flexibleengine_cce_node_v3.node_1]
}
`, rName, OS_VPC_ID, OS_NETWORK_ID, OS_AVAILABILITY_ZONE, OS_KEYPAIR_NAME)
}
