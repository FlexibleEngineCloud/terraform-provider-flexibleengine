package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/golangsdk/openstack/mrs/v1/cluster"
)

func TestAccMRSV1HybridCluster_basic(t *testing.T) {
	var mrsCluster cluster.Cluster

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckMrs(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMRSV1ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMRSV1HybridClusterConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV1ClusterExists("flexibleengine_mrs_hybrid_cluster_v1.cluster1", &mrsCluster),
					resource.TestCheckResourceAttr(
						"flexibleengine_mrs_hybrid_cluster_v1.cluster1", "cluster_name", "mrs-hybrid-cluster-acc"),
					resource.TestCheckResourceAttr(
						"flexibleengine_mrs_hybrid_cluster_v1.cluster1", "state", "running"),
					resource.TestCheckResourceAttr(
						"flexibleengine_mrs_hybrid_cluster_v1.cluster1", "log_collection", "1"),
					resource.TestCheckResourceAttr(
						"flexibleengine_mrs_hybrid_cluster_v1.cluster1", "total_node_number", "3"),
				),
			},
		},
	})
}

var testAccMRSV1HybridClusterConfig_basic = fmt.Sprintf(`
resource "flexibleengine_mrs_hybrid_cluster_v1" "cluster1" {
  available_zone  = "%s"
  cluster_name    = "mrs-hybrid-cluster-acc"
  cluster_version = "MRS 2.0.1"
  cluster_admin_secret  = "Cluster@123"
  master_node_key_pair = "KeyPair-ci"
  vpc_id = "%s"
  subnet_id = "%s"
  component_list = ["Hadoop", "Storm", "Spark", "Hive"]

  master_nodes {
    node_number = 1
    flavor = "s3.2xlarge.4.linux.mrs"
    data_volume_type = "SATA"
    data_volume_size = 100
    data_volume_count = 1
  }
  analysis_core_nodes {
    node_number = 1
    flavor = "s3.xlarge.4.linux.mrs"
    data_volume_type = "SATA"
    data_volume_size = 100
    data_volume_count = 1
  }

  streaming_core_nodes {
    node_number = 1
    flavor = "s3.xlarge.4.linux.mrs"
    data_volume_type = "SATA"
    data_volume_size = 100
    data_volume_count = 1
  }
}`, OS_AVAILABILITY_ZONE, OS_VPC_ID, OS_NETWORK_ID)
