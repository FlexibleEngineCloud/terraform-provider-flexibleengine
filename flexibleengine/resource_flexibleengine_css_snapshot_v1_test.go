package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/golangsdk/openstack/css/v1/snapshots"
)

func TestAccCssSnapshotV1_basic(t *testing.T) {
	rand := acctest.RandString(5)
	resourceKey := "flexibleengine_css_snapshot_v1.snapshot"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCssSnapshotV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCssSnapshotV1_basic(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCssSnapshotV1Exists(),
					resource.TestCheckResourceAttr(
						resourceKey, "name", fmt.Sprintf("snapshot-%s", rand)),
					resource.TestCheckResourceAttr(
						resourceKey, "status", "COMPLETED"),
					resource.TestCheckResourceAttr(
						resourceKey, "backup_type", "manual"),
				),
			},
		},
	})
}

func testAccCheckCssSnapshotV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.sdkClient(OS_REGION_NAME, "css")
	if err != nil {
		return fmt.Errorf("Error creating css client, err=%s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_css_snapshot_v1" {
			continue
		}

		clusterId := rs.Primary.Attributes["cluster_id"]
		snapList, err := snapshots.List(client, clusterId).Extract()
		if err != nil {
			return err
		}

		for _, v := range snapList {
			if v.ID == rs.Primary.ID {
				return fmt.Errorf("flexibleengine_css_snapshot_v1 %s still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckCssSnapshotV1Exists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		client, err := config.sdkClient(OS_REGION_NAME, "css")
		if err != nil {
			return fmt.Errorf("Error creating css client, err=%s", err)
		}

		rs, ok := s.RootModule().Resources["flexibleengine_css_snapshot_v1.snapshot"]
		if !ok {
			return fmt.Errorf("Error checking flexibleengine_css_snapshot_v1.snapshot exist, err=not found this resource")
		}

		clusterId := rs.Primary.Attributes["cluster_id"]
		snapList, err := snapshots.List(client, clusterId).Extract()
		if err != nil {
			return err
		}

		for _, v := range snapList {
			if v.ID == rs.Primary.ID {
				return nil
			}
		}

		return fmt.Errorf("flexibleengine_css_snapshot_v1 %s is not exist", rs.Primary.ID)
	}
}

func testAccCssSnapshotV1_basic(val string) string {
	return fmt.Sprintf(`
resource "flexibleengine_networking_secgroup_v2" "secgroup" {
  name = "terraform_test_sg-%s"
  description = "terraform security group acceptance test"
}

resource "flexibleengine_css_cluster_v1" "cluster" {
  name = "tf-css-cluster-%s"
  engine_version = "7.1.1"
  node_number    = 1

  node_config {
    flavor = "ess.spec-4u16g"
    network_info {
      security_group_id = flexibleengine_networking_secgroup_v2.secgroup.id
      subnet_id = "%s"
      vpc_id = "%s"
    }
    volume {
      volume_type = "COMMON"
      size = 40
    }
    availability_zone = "%s"
  }

  backup_strategy {
    start_time = "00:00 GMT+03:00"
    prefix     = "snapshot"
    keep_days  = 14
  }
}

resource "flexibleengine_css_snapshot_v1" "snapshot" {
  name        = "snapshot-%s"
  description = "a snapshot created by terraform acctest"
  cluster_id  = flexibleengine_css_cluster_v1.cluster.id
}
	`, val, val, OS_NETWORK_ID, OS_VPC_ID, OS_AVAILABILITY_ZONE, val)
}
