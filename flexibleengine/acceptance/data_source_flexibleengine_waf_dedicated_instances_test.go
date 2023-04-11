package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceWafDedicatedInstances_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	resourceName1 := "data.flexibleengine_waf_dedicated_instances.instance_1"
	resourceName2 := "data.flexibleengine_waf_dedicated_instances.instance_2"

	dc := acceptance.InitDataSourceCheck(resourceName1)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPrecheckWafInstance(t)
		},
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWafDedicatedInstances_conf(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName1, "name", name),
					resource.TestCheckResourceAttr(resourceName1, "instances.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName1, "instances.0.available_zone"),
					resource.TestCheckResourceAttrSet(resourceName1, "instances.0.cpu_flavor"),
					resource.TestCheckResourceAttrSet(resourceName1, "instances.0.cpu_architecture"),
					resource.TestCheckResourceAttrSet(resourceName1, "instances.0.security_group.#"),
					resource.TestCheckResourceAttrSet(resourceName1, "instances.0.server_id"),
					resource.TestCheckResourceAttrSet(resourceName1, "instances.0.service_ip"),
					resource.TestCheckResourceAttrSet(resourceName1, "instances.0.subnet_id"),
					resource.TestCheckResourceAttrSet(resourceName1, "instances.0.vpc_id"),
					resource.TestCheckResourceAttrSet(resourceName1, "instances.0.run_status"),
					resource.TestCheckResourceAttrSet(resourceName1, "instances.0.access_status"),
					resource.TestCheckResourceAttrSet(resourceName1, "instances.0.upgradable"),

					resource.TestCheckResourceAttr(resourceName2, "instances.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName2, "name"),
					resource.TestCheckResourceAttrSet(resourceName2, "instances.0.available_zone"),
				),
			},
		},
	})
}

func baseDependResource(name string) string {
	return fmt.Sprintf(`
resource "flexibleengine_vpc_v1" "vpc_1" {
  name = "%s_waf"
  cidr = "192.168.0.0/24"
}

resource "flexibleengine_vpc_subnet_v1" "vpc_subnet_1" {
  name       = "%s_waf"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = flexibleengine_vpc_v1.vpc_1.id
}

resource "flexibleengine_networking_secgroup_v2" "secgroup" {
  name        = "%s_waf"
  description = "terraform security group acceptance test"
}

data "flexibleengine_availability_zones" "zones" {}

data "flexibleengine_compute_flavors_v2" "flavors" {
  availability_zone = data.flexibleengine_availability_zones.zones.names[1]
  performance_type  = "normal"
  cpu_core          = 2
}
`, name, name, name)
}

func testAccWafDedicatedInstanceV1_conf(name string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_waf_dedicated_instance" "instance_1" {
  name               = "%s"
  available_zone     = data.flexibleengine_availability_zones.zones.names[1]
  specification_code = "waf.instance.professional"
  ecs_flavor         = data.flexibleengine_compute_flavors_v2.flavors.flavors[0]
  vpc_id             = flexibleengine_vpc_v1.vpc_1.id
  subnet_id          = flexibleengine_vpc_subnet_v1.vpc_subnet_1.id
  
  security_group = [
    flexibleengine_networking_secgroup_v2.secgroup.id
  ]
}
`, baseDependResource(name), name)
}

func testAccWafDedicatedInstances_conf(name string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_waf_dedicated_instances" "instance_1" {
  name = flexibleengine_waf_dedicated_instance.instance_1.name

  depends_on = [
    flexibleengine_waf_dedicated_instance.instance_1
  ]
}

data "flexibleengine_waf_dedicated_instances" "instance_2" {
  id   = flexibleengine_waf_dedicated_instance.instance_1.id
  name = flexibleengine_waf_dedicated_instance.instance_1.name

  depends_on = [
    flexibleengine_waf_dedicated_instance.instance_1
  ]
}
`, testAccWafDedicatedInstanceV1_conf(name))
}
