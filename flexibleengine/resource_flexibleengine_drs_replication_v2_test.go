package flexibleengine

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/golangsdk/openstack/drs/v2/replications"
)

// TestAccDRSV2Replication_basic is basic acc test.
func TestAccDRSV2Replication_basic(t *testing.T) {
	var replication replications.Replication

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDRSV2ReplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccDRSV2ReplicationConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDRSV2ReplicationExists("flexibleengine_drs_replication_v2.replication_1", &replication),
				),
			},
		},
	})
}

// testAccCheckDRSV2ReplicationDestroy checks destory.
func testAccCheckDRSV2ReplicationDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.drsV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine drs client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_drs_replication_v2" {
			continue
		}

		_, err := replications.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Replication still exists")
		}
	}

	log.Printf("[DEBUG] testAccCheckDRSV2ReplicationDestroy success!")

	return nil
}

// testAccCheckDRSV2ReplicationExists checks exist.
func testAccCheckDRSV2ReplicationExists(n string, replication *replications.Replication) resource.TestCheckFunc {
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

		found, err := replications.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Replication not found")
		}
		log.Printf("[DEBUG] test found is: %#v", found)
		*replication = *found

		return nil
	}
}

// TestAccDRSV2ReplicationConfig_basic is used to create.
var TestAccDRSV2ReplicationConfig_basic = `
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
`
