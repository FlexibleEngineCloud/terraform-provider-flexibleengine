package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/cse/dedicated/v2/engines"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getEngineFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.CseV2Client(OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CSE V2 client: %s", err)
	}
	return engines.Get(c, state.Primary.ID, state.Primary.Attributes["enterprise_project_id"])
}

func TestAccMicroserviceEngine_basic(t *testing.T) {
	var (
		engine       engines.Engine
		randName     = acceptance.RandomAccResourceNameWithDash()
		resourceName = "flexibleengine_cse_microservice_engine.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&engine,
		getEngineFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccMicroserviceEngine_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform test"),
					resource.TestCheckResourceAttr(resourceName, "flavor", "cse.s1.small"),
					resource.TestCheckResourceAttrPair(resourceName, "network_id", "data.flexibleengine_vpc_subnet_v1.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "auth_type", "RBAC"),
					resource.TestCheckResourceAttrSet(resourceName, "admin_pass"),
					resource.TestCheckResourceAttr(resourceName, "availability_zones.#", "3"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zones.0",
						"data.flexibleengine_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zones.1",
						"data.flexibleengine_availability_zones.test", "names.1"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zones.2",
						"data.flexibleengine_availability_zones.test", "names.2"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "service_limit", "200"),
					resource.TestCheckResourceAttr(resourceName, "instance_limit", "100"),
					resource.TestCheckResourceAttrSet(resourceName, "service_registry_addresses.0.private"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"admin_pass",
					"extend_params",
				},
			},
		},
	})
}

func testAccMicroserviceEngine_base(rName string) string {
	return fmt.Sprintf(`
data "flexibleengine_availability_zones" "test" {}

data "flexibleengine_vpc_subnet_v1" "test" {
  name = "subnet-default"
}

resource "flexibleengine_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "%s"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}
`, rName)
}

func testAccMicroserviceEngine_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_cse_microservice_engine" "test" {
  name                  = "%s"
  description           = "Created by terraform test"
  flavor                = "cse.s1.small"
  network_id            = data.flexibleengine_vpc_subnet_v1.test.id
  enterprise_project_id = "0"
  eip_id                = flexibleengine_vpc_eip.test.id

  auth_type  = "RBAC"
  admin_pass = "AccTest!123"

  availability_zones = slice(data.flexibleengine_availability_zones.test.names, 0, 3)

}`, testAccMicroserviceEngine_base(rName), rName)
}
