package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceLoadBalancers_basic(t *testing.T) {
	rName := "data.flexibleengine_elb_loadbalancers.test"
	dc := acceptance.InitDataSourceCheck(rName)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceLoadBalancers_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.#"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.name"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.id"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.ipv4_address"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.ipv4_port_id"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.l4_flavor_id"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.l7_flavor_id"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.vpc_id"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.enterprise_project_id"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("vpc_id_filter_is_useful", "true"),
					resource.TestCheckOutput("ipv4_subnet_id_filter_is_useful", "true"),
					resource.TestCheckOutput("description_filter_is_useful", "true"),
					resource.TestCheckOutput("l4_flavor_id_filter_is_useful", "true"),
					resource.TestCheckOutput("l7_flavor_id_filter_is_useful", "true"),
					resource.TestCheckOutput("type_is_useful", "true"),
					resource.TestCheckOutput("enterprise_project_id_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceLoadBalancers_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "flexibleengine_availability_zones" "test" {}

data "flexibleengine_elb_flavors" "l4flavors" {
  type            = "L4"
  max_connections = 1000000
  cps             = 20000
  bandwidth       = 100
}

data "flexibleengine_elb_flavors" "l7flavors" {
  type            = "L7"
  max_connections = 400000
  cps             = 4000
  bandwidth       = 100
}

resource "flexibleengine_lb_loadbalancer_v3" "test" {
  name           = "%[2]s"
  vpc_id         = flexibleengine_vpc_v1.test.id
  ipv4_subnet_id = flexibleengine_vpc_subnet_v1.test.ipv4_subnet_id
  description    = "update flavors"
  l4_flavor_id   = data.flexibleengine_elb_flavors.l4flavors.ids[0]
  l7_flavor_id   = data.flexibleengine_elb_flavors.l7flavors.ids[0]
  
  availability_zone = [
    data.flexibleengine_availability_zones.test.names[0]
  ]

  backend_subnets = [
    flexibleengine_vpc_subnet_v1.test.id
  ]

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, testVpc(rName), rName)
}

func testAccDatasourceLoadBalancers_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "flexibleengine_elb_loadbalancers" "test" {
  depends_on = [flexibleengine_lb_loadbalancer_v3.test]
}

data "flexibleengine_elb_loadbalancers" "name_filter" {
  name       = "%[2]s"
  depends_on = [flexibleengine_lb_loadbalancer_v3.test]
}

output "name_filter_is_useful" {
  value = length(data.flexibleengine_elb_loadbalancers.name_filter.loadbalancers) > 0 && alltrue(
  [for v in data.flexibleengine_elb_loadbalancers.name_filter.loadbalancers[*].name :v == "%[2]s"]
  )  
}

data "flexibleengine_elb_loadbalancers" "vpc_id_filter" {
  vpc_id     = flexibleengine_lb_loadbalancer_v3.test.vpc_id
  depends_on = [flexibleengine_lb_loadbalancer_v3.test]
}
locals {
  vpc_id = flexibleengine_lb_loadbalancer_v3.test.vpc_id
}
output "vpc_id_filter_is_useful" {
  value = length(data.flexibleengine_elb_loadbalancers.vpc_id_filter.loadbalancers) > 0 && alltrue(
  [for v in data.flexibleengine_elb_loadbalancers.vpc_id_filter.loadbalancers[*].vpc_id : v == local.vpc_id]
  )  
}

data "flexibleengine_elb_loadbalancers" "ipv4_subnet_id_filter" {
  ipv4_subnet_id = flexibleengine_lb_loadbalancer_v3.test.ipv4_subnet_id
  depends_on     = [flexibleengine_lb_loadbalancer_v3.test]
}
locals {
  ipv4_subnet_id = flexibleengine_lb_loadbalancer_v3.test.ipv4_subnet_id
}
output "ipv4_subnet_id_filter_is_useful" {
  value = length(data.flexibleengine_elb_loadbalancers.vpc_id_filter.loadbalancers) > 0 && alltrue(
  [for v in data.flexibleengine_elb_loadbalancers.ipv4_subnet_id_filter.loadbalancers[*].ipv4_subnet_id : 
  v == local.ipv4_subnet_id]
  )  
}

data "flexibleengine_elb_loadbalancers" "description_filter" {
  description = flexibleengine_lb_loadbalancer_v3.test.description
  depends_on  = [flexibleengine_lb_loadbalancer_v3.test]
}
locals {
  description = flexibleengine_lb_loadbalancer_v3.test.description
}
output "description_filter_is_useful" {
  value = length(data.flexibleengine_elb_loadbalancers.description_filter.loadbalancers) > 0 && alltrue(
  [for v in data.flexibleengine_elb_loadbalancers.description_filter.loadbalancers[*].description : 
  v == local.description]
  )  
}


data "flexibleengine_elb_loadbalancers" "l4_flavor_id_filter" {
  l4_flavor_id = flexibleengine_lb_loadbalancer_v3.test.l4_flavor_id
  depends_on   = [flexibleengine_lb_loadbalancer_v3.test]
}
locals {
  l4_flavor_id = flexibleengine_lb_loadbalancer_v3.test.l4_flavor_id
}
output "l4_flavor_id_filter_is_useful" {
  value = length(data.flexibleengine_elb_loadbalancers.l4_flavor_id_filter.loadbalancers) > 0 && alltrue(
  [for v in data.flexibleengine_elb_loadbalancers.l4_flavor_id_filter.loadbalancers[*].l4_flavor_id : 
  v == local.l4_flavor_id]
  )  
}

data "flexibleengine_elb_loadbalancers" "l7_flavor_id_filter" {
  l7_flavor_id = flexibleengine_lb_loadbalancer_v3.test.l7_flavor_id
  depends_on   = [flexibleengine_lb_loadbalancer_v3.test]
}
locals {
  l7_flavor_id = flexibleengine_lb_loadbalancer_v3.test.l7_flavor_id
}
output "l7_flavor_id_filter_is_useful" {
  value = length(data.flexibleengine_elb_loadbalancers.l7_flavor_id_filter.loadbalancers) > 0 && alltrue(
  [for v in data.flexibleengine_elb_loadbalancers.l7_flavor_id_filter.loadbalancers[*].l7_flavor_id : 
  v == local.l7_flavor_id]
  ) 
}


data "flexibleengine_elb_loadbalancers" "type_filter" {
  type       = "dedicated"
  depends_on = [flexibleengine_lb_loadbalancer_v3.test]
}
locals {
  type       = "dedicated"
}
output "type_is_useful" {
  value = length(data.flexibleengine_elb_loadbalancers.type_filter.loadbalancers) > 0 && alltrue(
  [for v in data.flexibleengine_elb_loadbalancers.type_filter.loadbalancers[*].type : 
  v == local.type]
  ) 
}

data "flexibleengine_elb_loadbalancers" "enterprise_project_id_filter" {
  depends_on = [flexibleengine_lb_loadbalancer_v3.test]
}
locals {
  enterprise_project_id = flexibleengine_lb_loadbalancer_v3.test.enterprise_project_id
}
output "enterprise_project_id_is_useful" {
  value = length(data.flexibleengine_elb_loadbalancers.enterprise_project_id_filter.loadbalancers) > 0 && alltrue(
  [for v in data.flexibleengine_elb_loadbalancers.enterprise_project_id_filter.loadbalancers[*].enterprise_project_id : 
  v == local.enterprise_project_id]
  ) 
}

`, testAccDatasourceLoadBalancers_base(name), name)
}
