package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/sdrs/v1/replications"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getReplicationPairResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := OS_REGION_NAME
	client, err := cfg.SdrsV1Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating SDRS Client: %s", err)
	}
	return replications.Get(client, state.Primary.ID).Extract()
}

func TestAccReplicationPair_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "flexibleengine_sdrs_replication_pair_v1.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getReplicationPairResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testReplicationPair_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "delete_target_volume", "true"),
					resource.TestCheckResourceAttrPair(rName, "group_id", "flexibleengine_sdrs_protectiongroup_v1.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "volume_id", "flexibleengine_blockstorage_volume_v2.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "replication_model"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "target_volume_id"),
				),
			},
			{
				Config: testReplicationPair_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", name)),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"delete_target_volume",
				},
			},
		},
	})
}

func testReplicationPair_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "flexibleengine_availability_zones" "test" {}
data "flexibleengine_sdrs_domain_v1" "test" {}

resource "flexibleengine_sdrs_protectiongroup_v1" "test" {
  name                     = "%[2]s"
  source_availability_zone = data.flexibleengine_availability_zones.test.names[0]
  target_availability_zone = data.flexibleengine_availability_zones.test.names[1]
  domain_id                = data.flexibleengine_sdrs_domain_v1.test.id
  source_vpc_id            = flexibleengine_vpc_v1.test.id
  description              = "test description"
}

resource "flexibleengine_blockstorage_volume_v2" "test" {
  name              = "%[2]s"
  description       = "test volume for sdrs replication pair"
  availability_zone = data.flexibleengine_availability_zones.test.names[0]
  volume_type       = "SSD"
  size              = 100
}
`, testVpc(name), name)
}

func testReplicationPair_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "flexibleengine_sdrs_replication_pair_v1" "test" {
  name                 = "%[2]s"
  group_id             = flexibleengine_sdrs_protectiongroup_v1.test.id
  volume_id            = flexibleengine_blockstorage_volume_v2.test.id
  description          = "test description"
  delete_target_volume = true
}
`, testReplicationPair_base(name), name)
}

func testReplicationPair_basic_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "flexibleengine_sdrs_replication_pair_v1" "test" {
  name                 = "%[2]s_update"
  group_id             = flexibleengine_sdrs_protectiongroup_v1.test.id
  volume_id            = flexibleengine_blockstorage_volume_v2.test.id
  description          = "test description"
  delete_target_volume = true
}
`, testReplicationPair_base(name), name)
}
