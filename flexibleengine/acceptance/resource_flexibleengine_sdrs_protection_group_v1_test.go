package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/sdrs/v1/protectiongroups"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getProtectionGroupResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := OS_REGION_NAME
	client, err := cfg.SdrsV1Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating SDRS Client: %s", err)
	}
	return protectiongroups.Get(client, state.Primary.ID).Extract()
}

// Lack of testing for `enable`, will test it in resource replication pair
func TestAccProtectionGroup_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "flexibleengine_sdrs_protectiongroup_v1.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getProtectionGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testProtectionGroup_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttrPair(rName, "source_availability_zone", "data.flexibleengine_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrPair(rName, "target_availability_zone", "data.flexibleengine_availability_zones.test", "names.1"),
					resource.TestCheckResourceAttrPair(rName, "domain_id", "data.flexibleengine_sdrs_domain_v1.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "source_vpc_id", "flexibleengine_vpc_v1.test", "id"),
				),
			},
			{
				Config: testProtectionGroup_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", name)),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testProtectionGroup_base(name string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_sdrs_domain_v1" "test" {}
data "flexibleengine_availability_zones" "test" {}
`, testVpc(name))
}

func testProtectionGroup_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "flexibleengine_sdrs_protectiongroup_v1" "test" {
  name                     = "%[2]s"
  source_availability_zone = data.flexibleengine_availability_zones.test.names[0]
  target_availability_zone = data.flexibleengine_availability_zones.test.names[1]
  domain_id                = data.flexibleengine_sdrs_domain_v1.test.id
  source_vpc_id            = flexibleengine_vpc_v1.test.id
  description              = "test description"
}
`, testProtectionGroup_base(name), name)
}

func testProtectionGroup_basic_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "flexibleengine_sdrs_protectiongroup_v1" "test" {
  name                     = "%[2]s_update"
  source_availability_zone = data.flexibleengine_availability_zones.test.names[0]
  target_availability_zone = data.flexibleengine_availability_zones.test.names[1]
  domain_id                = data.flexibleengine_sdrs_domain_v1.test.id
  source_vpc_id            = flexibleengine_vpc_v1.test.id
  description              = "test description"
}
`, testProtectionGroup_base(name), name)
}
