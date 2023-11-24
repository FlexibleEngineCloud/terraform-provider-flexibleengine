package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceDcsInstances_basic(t *testing.T) {
	rName := "data.flexibleengine_dcs_instances.test"
	name := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDcsInstances_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "instances.0.name", name),
					resource.TestCheckResourceAttr(rName, "instances.0.port", "6379"),
					resource.TestCheckResourceAttr(rName, "instances.0.flavor", "redis.cluster.xu1.large.r2.s1.8"),
				),
			},
		},
	})
}

func testAccDatasourceDcsInstances_basic(name string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_dcs_instances" "test" {
  name   = flexibleengine_dcs_instance_v1.instance_1.name
  status = "RUNNING"
}
`, testAccDcsV1Instance_basic(name))
}
