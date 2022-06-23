package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/mrs/v1/cluster"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

type GroupNodeNum struct {
	AnalysisCoreNum int
	StreamCoreNum   int
	AnalysisTaskNum int
	StreamTaskNum   int
}

func TestAccMrsMapReduceCluster_basic(t *testing.T) {
	var clusterGet cluster.Cluster
	resourceName := "flexibleengine_mrs_cluster_v2.test"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	password := fmt.Sprintf("TF%s%s%d", acctest.RandString(10), acctest.RandStringFromCharSet(1, "-_"),
		acctest.RandIntRange(0, 99))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckMRSV2ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMrsMapReduceClusterConfig_basic(rName, password),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2ClusterExists(resourceName, &clusterGet),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "STREAMING"),
					resource.TestCheckResourceAttr(resourceName, "safe_mode", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				Config: testAccMrsMapReduceClusterConfig_update(rName, password),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo1", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "update_value"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"manager_admin_pwd",
				},
			},
		},
	})
}

func TestAccMrsMapReduceCluster_analysis(t *testing.T) {
	var clusterGet cluster.Cluster
	resourceName := "flexibleengine_mrs_cluster_v2.test"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	password := fmt.Sprintf("TF%s%s%d", acctest.RandString(10), acctest.RandStringFromCharSet(1, "-_"),
		acctest.RandIntRange(0, 99))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckMRSV2ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMrsMapReduceClusterConfig_analysis(rName, password, buildGroupNodeNumbers(3, 0, 1, 0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2ClusterExists(resourceName, &clusterGet),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "ANALYSIS"),
					resource.TestCheckResourceAttr(resourceName, "safe_mode", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
					resource.TestCheckResourceAttr(resourceName, "analysis_core_nodes.0.node_number", "3"),
					resource.TestCheckResourceAttr(resourceName, "analysis_task_nodes.0.node_number", "1"),
				),
			},
			{
				Config: testAccMrsMapReduceClusterConfig_analysis(rName, password, buildGroupNodeNumbers(4, 0, 2, 0)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
					resource.TestCheckResourceAttr(resourceName, "analysis_core_nodes.0.node_number", "4"),
					resource.TestCheckResourceAttr(resourceName, "analysis_task_nodes.0.node_number", "2"),
				),
			},
			{
				Config: testAccMrsMapReduceClusterConfig_analysis(rName, password, buildGroupNodeNumbers(3, 0, 1, 0)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
					resource.TestCheckResourceAttr(resourceName, "analysis_core_nodes.0.node_number", "3"),
					resource.TestCheckResourceAttr(resourceName, "analysis_task_nodes.0.node_number", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"manager_admin_pwd",
				},
			},
		},
	})
}

func TestAccMrsMapReduceCluster_stream(t *testing.T) {
	var clusterGet cluster.Cluster
	resourceName := "flexibleengine_mrs_cluster_v2.test"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	password := fmt.Sprintf("TF%s%s%d", acctest.RandString(10), acctest.RandStringFromCharSet(1, "-_"),
		acctest.RandIntRange(0, 99))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckMRSV2ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMrsMapReduceClusterConfig_stream(rName, password, buildGroupNodeNumbers(0, 3, 0, 0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2ClusterExists(resourceName, &clusterGet),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "STREAMING"),
					resource.TestCheckResourceAttr(resourceName, "safe_mode", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
					resource.TestCheckResourceAttr(resourceName, "streaming_core_nodes.0.node_number", "3"),
				),
			},
			{
				Config: testAccMrsMapReduceClusterConfig_stream(rName, password, buildGroupNodeNumbers(0, 4, 0, 0)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
					resource.TestCheckResourceAttr(resourceName, "streaming_core_nodes.0.node_number", "4"),
				),
			},
			{
				Config: testAccMrsMapReduceClusterConfig_stream(rName, password, buildGroupNodeNumbers(0, 3, 0, 0)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
					resource.TestCheckResourceAttr(resourceName, "streaming_core_nodes.0.node_number", "3"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"manager_admin_pwd",
				},
			},
		},
	})
}

func TestAccMrsMapReduceCluster_hybrid(t *testing.T) {
	var clusterGet cluster.Cluster
	resourceName := "flexibleengine_mrs_cluster_v2.test"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	password := fmt.Sprintf("TF%s%s%d", acctest.RandString(10), acctest.RandStringFromCharSet(1, "-_"),
		acctest.RandIntRange(0, 99))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckMRSV2ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMrsMapReduceClusterConfig_hybrid(rName, password, buildGroupNodeNumbers(3, 3, 1, 1)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2ClusterExists(resourceName, &clusterGet),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "MIXED"),
					resource.TestCheckResourceAttr(resourceName, "safe_mode", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
					resource.TestCheckResourceAttr(resourceName, "analysis_core_nodes.0.node_number", "3"),
					resource.TestCheckResourceAttr(resourceName, "streaming_core_nodes.0.node_number", "3"),
					resource.TestCheckResourceAttr(resourceName, "analysis_task_nodes.0.node_number", "1"),
					resource.TestCheckResourceAttr(resourceName, "streaming_task_nodes.0.node_number", "1"),
					resource.TestCheckResourceAttr(resourceName, "analysis_core_nodes.0.host_ips.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "streaming_core_nodes.0.host_ips.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "analysis_task_nodes.0.host_ips.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "streaming_task_nodes.0.host_ips.#", "1"),
				),
			},
			{
				Config: testAccMrsMapReduceClusterConfig_hybrid(rName, password, buildGroupNodeNumbers(4, 4, 2, 2)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
					resource.TestCheckResourceAttr(resourceName, "analysis_core_nodes.0.node_number", "4"),
					resource.TestCheckResourceAttr(resourceName, "streaming_core_nodes.0.node_number", "4"),
					resource.TestCheckResourceAttr(resourceName, "analysis_task_nodes.0.node_number", "2"),
					resource.TestCheckResourceAttr(resourceName, "streaming_task_nodes.0.node_number", "2"),
					resource.TestCheckResourceAttr(resourceName, "analysis_core_nodes.0.host_ips.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "streaming_core_nodes.0.host_ips.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "analysis_task_nodes.0.host_ips.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "streaming_task_nodes.0.host_ips.#", "2"),
				),
			},
			{
				Config: testAccMrsMapReduceClusterConfig_hybrid(rName, password, buildGroupNodeNumbers(3, 3, 1, 1)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
					resource.TestCheckResourceAttr(resourceName, "analysis_core_nodes.0.node_number", "3"),
					resource.TestCheckResourceAttr(resourceName, "streaming_core_nodes.0.node_number", "3"),
					resource.TestCheckResourceAttr(resourceName, "analysis_task_nodes.0.node_number", "1"),
					resource.TestCheckResourceAttr(resourceName, "streaming_task_nodes.0.node_number", "1"),
					resource.TestCheckResourceAttr(resourceName, "analysis_core_nodes.0.host_ips.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "streaming_core_nodes.0.host_ips.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "analysis_task_nodes.0.host_ips.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "streaming_task_nodes.0.host_ips.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"manager_admin_pwd",
				},
			},
		},
	})
}

func TestAccMrsMapReduceCluster_custom(t *testing.T) {
	var clusterGet cluster.Cluster
	resourceName := "flexibleengine_mrs_cluster_v2.test"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	password := fmt.Sprintf("TF%s%s%d", acctest.RandString(10), acctest.RandStringFromCharSet(1, "-_"),
		acctest.RandIntRange(0, 99))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckMRSV2ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMrsMapReduceClusterConfig_custom(rName, password, 3),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2ClusterExists(resourceName, &clusterGet),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "CUSTOM"),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
					resource.TestCheckResourceAttr(resourceName, "custom_nodes.0.node_number", "3"),
					resource.TestCheckResourceAttr(resourceName, "custom_nodes.0.host_ips.#", "3"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"manager_admin_pwd",
					"template_id",
				},
			},
		},
	})
}

func TestAccMrsMapReduceCluster_publicIp(t *testing.T) {
	var clusterGet cluster.Cluster
	resourceName := "flexibleengine_mrs_cluster_v2.test"
	eipResourceName := "flexibleengine_vpc_eip_v1.test"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	bName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	password := fmt.Sprintf("TF%s%s%d", acctest.RandString(10), acctest.RandStringFromCharSet(1, "-_"),
		acctest.RandIntRange(0, 99))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckMRSV2ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMrsMapReduceClusterConfig_publicIp(rName, password, bName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2ClusterExists(resourceName, &clusterGet),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "STREAMING"),
					resource.TestCheckResourceAttr(resourceName, "safe_mode", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttrPair(resourceName, "eip_id", eipResourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "public_ip", eipResourceName, "address"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"manager_admin_pwd",
				},
			},
		},
	})
}

func buildGroupNodeNumbers(analysisCoreNum, streamCoreNum, analysisTaskNum, streamTaskNum int) GroupNodeNum {
	return GroupNodeNum{
		AnalysisCoreNum: analysisCoreNum,
		StreamCoreNum:   streamCoreNum,
		AnalysisTaskNum: analysisTaskNum,
		StreamTaskNum:   streamTaskNum,
	}
}

func testAccCheckMRSV2ClusterDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.MrsV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine MRS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_mrs_cluster_v2" {
			continue
		}

		clusterGet, err := cluster.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil
			}
			return fmt.Errorf("MRS cluster (%s) is still exists", rs.Primary.ID)
		}
		if clusterGet.Clusterstate == "terminated" {
			return nil
		}
	}

	return nil
}

func testAccCheckMRSV2ClusterExists(n string, clusterGet *cluster.Cluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource %s not found", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No MRS cluster ID")
		}

		config := testAccProvider.Meta().(*Config)
		mrsClient, err := config.MrsV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine MRS client: %s ", err)
		}

		found, err := cluster.Get(mrsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*clusterGet = *found
		return nil
	}
}

func testAccMrsMapReduceClusterConfig_base(rName string) string {
	return fmt.Sprintf(`
data "flexibleengine_availability_zones" "test" {}

resource "flexibleengine_vpc_v1" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "test" {
  name       = "%s"
  cidr       = "192.168.0.0/24"
  vpc_id     = flexibleengine_vpc_v1.test.id
  gateway_ip = "192.168.0.1"
}

resource "flexibleengine_compute_keypair_v2" "test" {
  name = "%s"
}`, rName, rName, rName)
}

// The task node has not contain data disks.
func testAccMrsMapReduceClusterConfig_basic(rName, pwd string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_mrs_cluster_v2" "test" {
  availability_zone  = data.flexibleengine_availability_zones.test.names[0]
  name               = "%s"
  type               = "STREAMING"
  version            = "MRS 2.0.1"
  manager_admin_pwd  = "%s"
  node_key_pair      = flexibleengine_compute_keypair_v2.test.name
  subnet_id          = flexibleengine_vpc_subnet_v1.test.id
  vpc_id             = flexibleengine_vpc_v1.test.id
  component_list     = ["Storm", "Kafka"]

  master_nodes {
    flavor            = "c6.4xlarge.4.linux.mrs"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
  streaming_core_nodes {
    flavor            = "c6.4xlarge.4.linux.mrs"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
  streaming_task_nodes {
    flavor            = "c6.4xlarge.4.linux.mrs"
    node_number       = 1
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_count = 0
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}`, testAccMrsMapReduceClusterConfig_base(rName), rName, pwd)
}

func testAccMrsMapReduceClusterConfig_update(rName, pwd string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_mrs_cluster_v2" "test" {
  availability_zone  = data.flexibleengine_availability_zones.test.names[0]
  name               = "%s"
  type               = "STREAMING"
  version            = "MRS 2.0.1"
  manager_admin_pwd  = "%s"
  node_key_pair      = flexibleengine_compute_keypair_v2.test.name
  subnet_id          = flexibleengine_vpc_subnet_v1.test.id
  vpc_id             = flexibleengine_vpc_v1.test.id
  component_list     = ["Storm", "Kafka"]

  master_nodes {
    flavor            = "c6.4xlarge.4.linux.mrs"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
  streaming_core_nodes {
    flavor            = "c6.4xlarge.4.linux.mrs"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
  streaming_task_nodes {
    flavor            = "c6.4xlarge.4.linux.mrs"
    node_number       = 1
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_count = 0
  }

  tags = {
    foo1 = "bar"
    key  = "update_value"
  }
}`, testAccMrsMapReduceClusterConfig_base(rName), rName, pwd)
}

func testAccMrsMapReduceClusterConfig_analysis(rName, pwd string, nodeNums GroupNodeNum) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_mrs_cluster_v2" "test" {
  availability_zone  = data.flexibleengine_availability_zones.test.names[0]
  name               = "%s"
  type               = "ANALYSIS"
  version            = "MRS 3.1.0-LTS.1"
  manager_admin_pwd  = "%s"
  node_key_pair      = flexibleengine_compute_keypair_v2.test.name
  subnet_id          = flexibleengine_vpc_subnet_v1.test.id
  vpc_id             = flexibleengine_vpc_v1.test.id
  component_list     = ["Hadoop", "Hive", "Tez"]

  master_nodes {
    flavor            = "s3.4xlarge.4.linux.mrs"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 480
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1
  }
  analysis_core_nodes {
    flavor            = "s3.4xlarge.4.linux.mrs"
    node_number       = %d
    root_volume_type  = "SAS"
    root_volume_size  = 480
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1
  }
  analysis_task_nodes {
    flavor            = "s3.4xlarge.4.linux.mrs"
    node_number       = %d
    root_volume_type  = "SAS"
    root_volume_size  = 480
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1
  }
}`, testAccMrsMapReduceClusterConfig_base(rName), rName, pwd,
		nodeNums.AnalysisCoreNum, nodeNums.AnalysisTaskNum)
}

func testAccMrsMapReduceClusterConfig_stream(rName, pwd string, nodeNums GroupNodeNum) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_mrs_cluster_v2" "test" {
  availability_zone  = data.flexibleengine_availability_zones.test.names[0]
  name               = "%s"
  type               = "STREAMING"
  version            = "MRS 3.1.0-LTS.1"
  manager_admin_pwd  = "%s"
  node_key_pair      = flexibleengine_compute_keypair_v2.test.name
  subnet_id          = flexibleengine_vpc_subnet_v1.test.id
  vpc_id             = flexibleengine_vpc_v1.test.id
  component_list     = ["Ranger", "Kafka", "ZooKeeper"]

  master_nodes {
    flavor            = "c6.4xlarge.4.linux.mrs"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 480
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1
  }
  streaming_core_nodes {
    flavor            = "c6.4xlarge.4.linux.mrs"
    node_number       = %d
    root_volume_type  = "SAS"
    root_volume_size  = 480
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1
  }
}`, testAccMrsMapReduceClusterConfig_base(rName), rName, pwd, nodeNums.StreamCoreNum)
}

func testAccMrsMapReduceClusterConfig_hybrid(rName, pwd string, nodeNums GroupNodeNum) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_mrs_cluster_v2" "test" {
  availability_zone  = data.flexibleengine_availability_zones.test.names[0]
  name               = "%s"
  type               = "MIXED"
  version            = "MRS 2.0.1"
  safe_mode          = true
  manager_admin_pwd  = "%s"
  node_key_pair      = flexibleengine_compute_keypair_v2.test.name
  subnet_id          = flexibleengine_vpc_subnet_v1.test.id
  vpc_id             = flexibleengine_vpc_v1.test.id
  component_list     = ["Hadoop", "Spark", "Hive", "Tez", "Storm"]

  master_nodes {
    flavor            = "c6.4xlarge.4.linux.mrs"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
  analysis_core_nodes {
    flavor            = "c6.4xlarge.4.linux.mrs"
    node_number       = %d
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
  streaming_core_nodes {
    flavor            = "c6.4xlarge.4.linux.mrs"
    node_number       = %d
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
  analysis_task_nodes {
    flavor            = "c6.4xlarge.4.linux.mrs"
    node_number       = %d
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
  streaming_task_nodes {
    flavor            = "c6.4xlarge.4.linux.mrs"
    node_number       = %d
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
}`, testAccMrsMapReduceClusterConfig_base(rName), rName, pwd,
		nodeNums.AnalysisCoreNum, nodeNums.StreamCoreNum, nodeNums.AnalysisTaskNum, nodeNums.StreamTaskNum)
}

func testAccMrsMapReduceClusterConfig_custom(rName, pwd string, nodeNum1 int) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_mrs_cluster_v2" "test" {
  availability_zone  = data.flexibleengine_availability_zones.test.names[0]
  name               = "%s"
  type               = "CUSTOM"
  version            = "MRS 3.1.0-LTS.1"
  safe_mode          = true
  manager_admin_pwd  = "%s"
  node_key_pair      = flexibleengine_compute_keypair_v2.test.name
  subnet_id          = flexibleengine_vpc_subnet_v1.test.id
  vpc_id             = flexibleengine_vpc_v1.test.id
  template_id        = "mgmt_control_combined_v4"
  component_list     = ["DBService", "Hadoop", "ZooKeeper", "Ranger", "ClickHouse"]

master_nodes {
    flavor            = "c6.4xlarge.4.linux.mrs"
    node_number       = 3
    root_volume_type  = "SAS"
    root_volume_size  = 480
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1
    assigned_roles = [
      "OMSServer:1,2",
      "SlapdServer:1,2",
      "KerberosServer:1,2",
      "KerberosAdmin:1,2",
      "quorumpeer:1,2,3",
      "NameNode:2,3",
      "Zkfc:2,3",
      "JournalNode:1,2,3",
      "ResourceManager:2,3",
      "JobHistoryServer:3",
      "DBServer:1,3",
      "HttpFS:1,3",
      "TimelineServer:3",
      "RangerAdmin:1,2",
      "UserSync:2",
      "TagSync:2",
      "KerberosClient",
      "SlapdClient",
      "meta",
      "ClickHouseBalancer:1,2"
    ]
  }

  custom_nodes {
    group_name        = "node_group_1"
    flavor            = "c6.4xlarge.4.linux.mrs"
    node_number       = %d
    root_volume_type  = "SAS"
    root_volume_size  = 480
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1
    assigned_roles = [
      "DataNode",
      "NodeManager",
      "KerberosClient",
      "SlapdClient",
      "meta"
    ]
  }

  custom_nodes {
    group_name        = "ClickHouse"
    flavor            = "c6.4xlarge.4.linux.mrs"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 480
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1
    assigned_roles = [
      "ClickHouseServer",
      "meta",
      "KerberosClient",
      "SlapdClient"
    ]
  }
  
}`, testAccMrsMapReduceClusterConfig_base(rName), rName, pwd, nodeNum1)
}

func testAccMrsMapReduceClusterConfig_publicIp(rName, pwd, bName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_vpc_eip_v1" "test" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
	name        = "%s"
    share_type  = "PER"
    size        = 5
    charge_mode = "traffic"
  }
}

resource "flexibleengine_mrs_cluster_v2" "test" {
  availability_zone  = data.flexibleengine_availability_zones.test.names[0]
  name               = "%s"
  type               = "STREAMING"
  version            = "MRS 2.0.1"
  manager_admin_pwd  = "%s"
  node_key_pair      = flexibleengine_compute_keypair_v2.test.name
  subnet_id          = flexibleengine_vpc_subnet_v1.test.id
  vpc_id             = flexibleengine_vpc_v1.test.id
  public_ip          = flexibleengine_vpc_eip_v1.test.address
  component_list     = ["Storm"]

  master_nodes {
    flavor            = "c6.4xlarge.4.linux.mrs"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
  streaming_core_nodes {
    flavor            = "c6.4xlarge.4.linux.mrs"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
  streaming_task_nodes {
    flavor            = "c6.4xlarge.4.linux.mrs"
    node_number       = 1
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_count = 0
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}`, testAccMrsMapReduceClusterConfig_base(rName), bName, rName, pwd)
}
