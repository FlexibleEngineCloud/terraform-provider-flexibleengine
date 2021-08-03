package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/golangsdk/openstack/sdrs/v1/protectedinstances"
)

func TestAccSdrsReplicationAttachV1_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckSdrs(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSdrsAttachV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSdrsAttachV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSdrsAttachV1Exists("flexibleengine_sdrs_replication_attach_v1.attach_1"),
					resource.TestCheckResourceAttr(
						"flexibleengine_sdrs_replication_attach_v1.attach_1", "device", "/dev/vdb"),
				),
			},
		},
	})
}

func testAccCheckSdrsAttachV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	sdrsClient, err := config.sdrsV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine SDRS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_sdrs_replication_attach_v1" {
			continue
		}

		instId, replicaId := ExtractAttachId(rs.Primary.ID)
		n, err := protectedinstances.Get(sdrsClient, instId).Extract()
		if err != nil {
			return nil
		}

		for _, attach := range n.Attachment {
			if attach.Replication == replicaId {
				return fmt.Errorf("SDRS DR drill still exists")
			}
		}
	}

	return nil
}

func testAccCheckSdrsAttachV1Exists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		sdrsClient, err := config.sdrsV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine SDRS client: %s", err)
		}

		instId, replicaId := ExtractAttachId(rs.Primary.ID)
		n, err := protectedinstances.Get(sdrsClient, instId).Extract()
		if err != nil {
			return err
		}

		find := false
		for _, attach := range n.Attachment {
			if attach.Replication == replicaId {
				find = true
				break
			}
		}

		if find == false {
			return fmt.Errorf("SDRS Replication attach not found")
		}
		return nil
	}
}

var testAccSdrsAttachV1_basic = fmt.Sprintf(`
data "flexibleengine_sdrs_domain_v1" "domain_1" {
  name = "SDRS_HypeDomain01"
}

# create ecs server for protection 
resource "flexibleengine_compute_instance_v2" "server_1" {
  name = "server_1"
  security_groups = ["default"]
  availability_zone = "eu-west-0a"
  network {
    uuid = "%s"
  }
}

resource "flexibleengine_sdrs_protectiongroup_v1" "group_1" {
  name = "group_1"
  description = "test description"
  source_availability_zone = "eu-west-0a"
  target_availability_zone = "eu-west-0b"
  domain_id = data.flexibleengine_sdrs_domain_v1.domain_1.id
  source_vpc_id = "%s"
  dr_type = "migration"
}

resource "flexibleengine_sdrs_protectedinstance_v1" "instance_1" {
  name        = "instance_1"
  description = "test description"
  group_id    = flexibleengine_sdrs_protectiongroup_v1.group_1.id
  server_id   = flexibleengine_compute_instance_v2.server_1.id
  delete_target_server = true
  delete_target_eip = true
}

resource "flexibleengine_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  description = "volume for replication pair"
  availability_zone = "eu-west-0a"
  size = 1
}

resource "flexibleengine_sdrs_replication_pair_v1" "replication_1" {
  name = "replication_1"
  description = "test replication pair"
  group_id = flexibleengine_sdrs_protectiongroup_v1.group_1.id
  volume_id = flexibleengine_blockstorage_volume_v2.volume_1.id
  delete_target_volume = true
}

resource "flexibleengine_sdrs_replication_attach_v1" "attach_1" {
  instance_id = flexibleengine_sdrs_protectedinstance_v1.instance_1.id
  replication_id = flexibleengine_sdrs_replication_pair_v1.replication_1.id
  device = "/dev/vdb"
}`, OS_NETWORK_ID, OS_VPC_ID)
