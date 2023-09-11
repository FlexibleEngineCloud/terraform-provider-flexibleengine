package acceptance

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceDdmInstances_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	name = strings.ReplaceAll(name, "_", "-")
	rName := "data.flexibleengine_ddm_instances.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDdmInstances_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "instances.0.status", "RUNNING"),
					resource.TestCheckResourceAttr(rName, "instances.0.name", name),
					acceptance.TestCheckResourceAttrWithVariable(rName, "instances.0.vpc_id",
						"${flexibleengine_vpc_v1.test.id}"),
					acceptance.TestCheckResourceAttrWithVariable(rName, "instances.0.subnet_id",
						"${flexibleengine_vpc_subnet_v1.test.id}"),
					acceptance.TestCheckResourceAttrWithVariable(rName, "instances.0.security_group_id",
						"${flexibleengine_networking_secgroup_v2.test.id}"),
					resource.TestCheckResourceAttr(rName, "instances.0.node_num", "2"),
					resource.TestCheckResourceAttr(rName, "instances.0.engine_version", "3.0.8.5"),
				),
			},
		},
	})
}

func testAccDatasourceDdmInstances_basic(name string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_ddm_instances" "test" {
  name = flexibleengine_ddm_instance.test.name
}
`, testDdmInstance_basic(name))
}
