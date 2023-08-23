package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/nat/v3/gateways"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/nat"
)

func getPrivateGatewayResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NatV3Client(OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating NAT v3 client: %s", err)
	}

	return gateways.Get(client, state.Primary.ID)
}

func TestAccPrivateGateway_basic(t *testing.T) {
	var (
		obj gateways.Gateway

		rName      = "flexibleengine_nat_private_gateway.test"
		name       = acceptance.RandomAccResourceNameWithDash()
		updateName = acceptance.RandomAccResourceNameWithDash()
		baseConfig = testBaseNetwork(name)
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPrivateGatewayResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateGateway_basic_step_1(name, baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "subnet_id", "flexibleengine_vpc_subnet_v1.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "spec", nat.PrivateSpecTypeSmall),
					resource.TestCheckResourceAttr(rName, "description", "Created by acc test"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
				),
			},
			{
				Config: testAccPrivateGateway_basic_step_2(updateName, baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "subnet_id", "flexibleengine_vpc_subnet_v1.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "baaar"),
					resource.TestCheckResourceAttr(rName, "tags.newKey", "value"),
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

func testAccPrivateGateway_basic_step_1(name, relatedConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "flexibleengine_nat_private_gateway" "test" {
  subnet_id             = flexibleengine_vpc_subnet_v1.test.id
  name                  = "%[2]s"
  description           = "Created by acc test"
  enterprise_project_id = "0"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, relatedConfig, name)
}

func testAccPrivateGateway_basic_step_2(name, relatedConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "flexibleengine_nat_private_gateway" "test" {
  subnet_id             = flexibleengine_vpc_subnet_v1.test.id
  name                  = "%[2]s"
  enterprise_project_id = "0"

  tags = {
    foo    = "baaar"
    newKey = "value"
  }
}
`, relatedConfig, name)
}
