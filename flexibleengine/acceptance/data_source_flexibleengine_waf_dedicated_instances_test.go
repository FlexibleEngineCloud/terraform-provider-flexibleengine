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
