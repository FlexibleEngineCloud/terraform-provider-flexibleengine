package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/mrs/v2/jobs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMRSV2Job_basic(t *testing.T) {
	var job jobs.Job
	resourceName := "flexibleengine_mrs_job_v2.test"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	pwd := fmt.Sprintf("TF%s%s%d", acctest.RandString(10), acctest.RandStringFromCharSet(1, "-_"),
		acctest.RandIntRange(0, 99))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMRSV2JobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMRSV2JobConfig_basic(rName, pwd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2JobExists(resourceName, &job),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "SparkSubmit"),
					resource.TestCheckResourceAttr(resourceName, "program_path",
						"s3a://obs-demo-analysis/program/driver_behavior.jar"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccMRSJobImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccCheckMRSV2JobDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.MrsV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine MRS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_mrs_job_v2" {
			continue
		}

		_, err := jobs.Get(client, rs.Primary.Attributes["cluster_id"], rs.Primary.ID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil
			}
			return fmt.Errorf("MRS cluster (%s) is still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckMRSV2JobExists(n string, job *jobs.Job) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource %s not found", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No MRS cluster ID")
		}

		config := testAccProvider.Meta().(*Config)
		client, err := config.MrsV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine MRS client: %s ", err)
		}

		found, err := jobs.Get(client, rs.Primary.Attributes["cluster_id"], rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*job = *found
		return nil
	}
}

func testAccMRSJobImportStateIdFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.ID == "" || rs.Primary.Attributes["cluster_id"] == "" {
			return "", fmt.Errorf("resource not found: %s/%s", rs.Primary.Attributes["cluster_id"], rs.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["cluster_id"], rs.Primary.ID), nil
	}
}

func testAccMRSV2JobConfig_basic(rName, pwd string) string {
	return fmt.Sprintf(`
resource "flexibleengine_mrs_cluster_v1" "cluster1" {
  available_zone_id = "%s"
  cluster_name      = "%s"
  cluster_version   = "MRS 1.8.9"
  cluster_type      = 0
  master_node_num   = 2
  core_node_num     = 3
  master_node_size  = "s3.2xlarge.4.linux.mrs"
  core_node_size    = "s3.xlarge.4.linux.mrs"

  node_public_cert_name = "KeyPair-ci"
  safe_mode             = 0
  cluster_admin_secret  = "%s"

  volume_type = "SATA"
  volume_size = 100
  vpc_id      = "%s"
  subnet_id   = "%s"

  component_list {
      component_name = "Hadoop"
  }
  component_list {
      component_name = "Spark"
  }
  component_list {
      component_name = "Hive"
  }
}

resource "flexibleengine_mrs_job_v2" "test" {
  cluster_id   = flexibleengine_mrs_cluster_v1.cluster1.id
  name         = "%s"
  type         = "SparkSubmit"
  program_path = "s3a://obs-demo-analysis/program/driver_behavior.jar"
  parameters   = "%s %s 1 s3a://obs-demo-analysis/input s3a://obs-demo-analysis/output"

  program_parameters = {
    "--class" = "com.huawei.bigdata.spark.examples.DriverBehavior"
  }
}`, OS_AVAILABILITY_ZONE, rName, pwd, OS_VPC_ID, OS_NETWORK_ID, rName, OS_ACCESS_KEY, OS_SECRET_KEY)
}
