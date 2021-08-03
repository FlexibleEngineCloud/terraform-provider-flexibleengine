package flexibleengine

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/golangsdk/openstack/drs/v2/replicationconsistencygroups"
)

// TestAccDRSV2ReplicationConsistencyGroup_basic is basic acc test.
func TestAccDRSV2ReplicationConsistencyGroup_basic(t *testing.T) {
	var rcg replicationconsistencygroups.ReplicationConsistencyGroup

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDRSV2ReplicationConsistencyGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccDRSV2ReplicationConsistencyGroupConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDRSV2ReplicationConsistencyGroupExists(
						"flexibleengine_drs_replicationconsistencygroup_v2.replicationconsistencygroup_1",
						&rcg),
				),
			},
			{
				Config: TestAccDRSV2ReplicationConsistencyGroupConfig_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"flexibleengine_drs_replicationconsistencygroup_v2.replicationconsistencygroup_1",
						"name",
						"replicationconsistencygroup_1_updated"),
					resource.TestCheckResourceAttr(
						"flexibleengine_drs_replicationconsistencygroup_v2.replicationconsistencygroup_1",
						"description",
						"The description of replicationconsistencygroup_1_updated"),
				),
			},
		},
	})
}

// testAccCheckDRSV2ReplicationConsistencyGroupDestroy checks destory.
func testAccCheckDRSV2ReplicationConsistencyGroupDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.drsV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine drs client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_drs_replicationconsistencygroup_v2" {
			continue
		}

		_, err := replicationconsistencygroups.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("replicationconsistencygroup still exists")
		}
	}

	log.Printf("[DEBUG] testAccCheckDRSV2ReplicationConsistencyGroupDestroy success!")

	return nil
}

// testAccCheckDRSV2ReplicationConsistencyGroupExists checks exist.
func testAccCheckDRSV2ReplicationConsistencyGroupExists(n string, rcg *replicationconsistencygroups.ReplicationConsistencyGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		client, err := config.drsV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine drs client: %s", err)
		}

		found, err := replicationconsistencygroups.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("replicationconsistencygroup not found")
		}
		log.Printf("[DEBUG] test found is: %#v", found)
		*rcg = *found

		return nil
	}
}

// TestAccDRSV2ReplicationConsistencyGroupConfig_basic is used to create.
var TestAccDRSV2ReplicationConsistencyGroupConfig_basic = `
resource "flexibleengine_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  size = 1
  availability_zone = "eu-west-0a"
}

resource "flexibleengine_blockstorage_volume_v2" "volume_2" {
  name = "volume_2"
  size = 1
  availability_zone = "eu-west-0b"
}

resource "flexibleengine_drs_replication_v2" "replication_1" {
  name = "replication_1"
  description = "The description of replication_1"
  volume_ids = ["${flexibleengine_blockstorage_volume_v2.volume_1.id}", "${flexibleengine_blockstorage_volume_v2.volume_2.id}"]
  priority_station = "eu-west-0a"
}

resource "flexibleengine_drs_replicationconsistencygroup_v2" "replicationconsistencygroup_1" {
  name = "replicationconsistencygroup_1"
  description = "The description of replicationconsistencygroup_1"
  replication_ids = ["${flexibleengine_drs_replication_v2.replication_1.id}"]
  priority_station = "eu-west-0a"
}
`

// TestAccDRSV2ReplicationConsistencyGroupConfig_update is used to update.
var TestAccDRSV2ReplicationConsistencyGroupConfig_update = `
resource "flexibleengine_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  size = 1
  availability_zone = "eu-west-0a"
}

resource "flexibleengine_blockstorage_volume_v2" "volume_2" {
  name = "volume_2"
  size = 1
  availability_zone = "eu-west-0b"
}

resource "flexibleengine_drs_replication_v2" "replication_1" {
  name = "replication_1"
  description = "The description of replication_1"
  volume_ids = ["${flexibleengine_blockstorage_volume_v2.volume_1.id}", "${flexibleengine_blockstorage_volume_v2.volume_2.id}"]
  priority_station = "eu-west-0a"
}

resource "flexibleengine_drs_replicationconsistencygroup_v2" "replicationconsistencygroup_1" {
  name = "replicationconsistencygroup_1_updated"
  description = "The description of replicationconsistencygroup_1_updated"
  replication_ids = ["${flexibleengine_drs_replication_v2.replication_1.id}"]
  priority_station = "eu-west-0a"
}
`
