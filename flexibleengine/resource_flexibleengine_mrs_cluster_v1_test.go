package flexibleengine

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/mrs/v1/cluster"
)

func TestAccMRSV1Cluster_basic(t *testing.T) {
	var clusterGet cluster.Cluster

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckMrs(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMRSV1ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccMRSV1ClusterConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV1ClusterExists("flexibleengine_mrs_cluster_v1.cluster1", &clusterGet),
					resource.TestCheckResourceAttr(
						"flexibleengine_mrs_cluster_v1.cluster1", "cluster_state", "running"),
				),
			},
		},
	})
}

func testAccCheckMRSV1ClusterDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	mrsClient, err := config.MrsV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating flexibleengine mrs: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_mrs_cluster_v1" {
			continue
		}

		clusterGet, err := cluster.Get(mrsClient, rs.Primary.ID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil
			}
			return fmt.Errorf("cluster still exists. err : %s", err)
		}
		if clusterGet.Clusterstate == "terminated" {
			return nil
		}
	}

	return nil
}

func testAccCheckMRSV1ClusterExists(n string, clusterGet *cluster.Cluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s. ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set. ")
		}

		config := testAccProvider.Meta().(*Config)
		mrsClient, err := config.MrsV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating flexibleengine mrs client: %s ", err)
		}

		found, err := cluster.Get(mrsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Clusterid != rs.Primary.ID {
			return fmt.Errorf("Cluster not found. ")
		}

		*clusterGet = *found
		time.Sleep(5 * time.Second)

		return nil
	}
}

var TestAccMRSV1ClusterConfig_basic = fmt.Sprintf(`
resource "flexibleengine_mrs_cluster_v1" "cluster1" {
  cluster_name = "mrs-cluster-acc"
  region = "%s"
  billing_type = 12
  master_node_num = 2
  core_node_num = 3
  master_node_size = "s1.4xlarge.linux.mrs"
  core_node_size = "s1.xlarge.linux.mrs"
  available_zone_id = "%s"
  vpc_id = "%s"
  subnet_id = "%s"
  cluster_version = "MRS 1.5.0"
  volume_type = "SATA"
  volume_size = 100
  safe_mode = 0
  cluster_type = 0
  node_public_cert_name = "KeyPair-ci"
  cluster_admin_secret = ""
  component_list {
      component_name = "Hadoop"
  }
  component_list {
      component_name = "Spark"
  }
  component_list {
      component_name = "Hive"
  }
}`, OS_REGION_NAME, OS_AVAILABILITY_ZONE, OS_VPC_ID, OS_NETWORK_ID)
