package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/nat/v2/snats"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getPublicSnatRuleResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NatGatewayClient(OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating NAT v2 client: %s", err)
	}

	return snats.Get(client, state.Primary.ID)
}

func TestAccPublicSnatRule_basic(t *testing.T) {
	var (
		obj snats.Rule

		rName = "flexibleengine_nat_snat_rule_v2.test"
		name  = acceptance.RandomAccResourceNameWithDash()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPublicSnatRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPublicSnatRule_basic_step_1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "nat_gateway_id", "flexibleengine_nat_gateway_v2.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "subnet_id", "flexibleengine_vpc_subnet_v1.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "floating_ip_id", "flexibleengine_vpc_eip.test.0", "id"),
					resource.TestCheckResourceAttr(rName, "description", "Created by acc test"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
				),
			},
			{
				Config: testAccPublicSnatRule_basic_step_2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "nat_gateway_id", "flexibleengine_nat_gateway_v2.test", "id"),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
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

func testAccPublicSnatRule_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "flexibleengine_vpc_eip" "test" {
  count = 2

  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = format("%[2]s-%%d", count.index)
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "flexibleengine_compute_instance_v2" "test" {
  name              = "instance_1"
  security_groups   = ["default"]
  image_id          = data.flexibleengine_images_image.test.id
  flavor_id         = data.flexibleengine_compute_flavors_v2.test.flavors[0]
  availability_zone = data.flexibleengine_availability_zones.test.names[0]
  metadata = {
    foo = "bar"
  }
  network {
    uuid = flexibleengine_vpc_subnet_v1.test.id
  }
  tags = {
    key1 = "value1"
    key2 = "value.key"
  }
}

resource "flexibleengine_nat_gateway_v2" "test" {
  name        = "%[2]s"
  description = "test for terraform"
  spec        = "2"
  vpc_id      = flexibleengine_vpc_v1.test.id
  subnet_id   = flexibleengine_vpc_subnet_v1.test.id
}
`, testBaseComputeResources(name), name)
}

func testAccPublicSnatRule_basic_step_1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "flexibleengine_nat_snat_rule_v2" "test" {
  nat_gateway_id = flexibleengine_nat_gateway_v2.test.id
  subnet_id      = flexibleengine_vpc_subnet_v1.test.id
  floating_ip_id = flexibleengine_vpc_eip.test[0].id
  description    = "Created by acc test"
}
`, testAccPublicSnatRule_base(name))
}

func testAccPublicSnatRule_basic_step_2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "flexibleengine_nat_snat_rule_v2" "test" {
  nat_gateway_id = flexibleengine_nat_gateway_v2.test.id
  subnet_id      = flexibleengine_vpc_subnet_v1.test.id
  floating_ip_id = join(",", flexibleengine_vpc_eip.test[*].id)
}
`, testAccPublicSnatRule_base(name))
}
