package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceListeners_basic(t *testing.T) {
	var (
		rName            = acceptance.RandomAccResourceNameWithDash()
		dcByName         = acceptance.InitDataSourceCheck("data.flexibleengine_lb_listeners.by_name")
		dcByProtocol     = acceptance.InitDataSourceCheck("data.flexibleengine_lb_listeners.by_protocol")
		dcByProtocolPort = acceptance.InitDataSourceCheck("data.flexibleengine_lb_listeners.by_protocol_port")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceListeners_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("name_query_result_validation", "true"),
					resource.TestCheckResourceAttrSet("data.flexibleengine_lb_listeners.by_name",
						"listeners.0.name"),
					resource.TestCheckResourceAttrSet("data.flexibleengine_lb_listeners.by_name",
						"listeners.0.protocol"),
					resource.TestCheckResourceAttrSet("data.flexibleengine_lb_listeners.by_name",
						"listeners.0.protocol_port"),
					resource.TestCheckResourceAttrSet("data.flexibleengine_lb_listeners.by_name",
						"listeners.0.connection_limit"),
					resource.TestCheckResourceAttrSet("data.flexibleengine_lb_listeners.by_name",
						"listeners.0.http2_enable"),
					resource.TestCheckResourceAttr("data.flexibleengine_lb_listeners.by_name",
						"listeners.0.loadbalancers.#", "1"),
					dcByProtocol.CheckResourceExists(),
					resource.TestCheckOutput("protocol_query_result_validation", "true"),
					dcByProtocolPort.CheckResourceExists(),
					resource.TestCheckOutput("protocol_port_query_result_validation", "true"),
				),
			},
		},
	})
}

func testAccDatasourceListeners_base(rName string) string {
	rCidr := acceptance.RandomCidr()
	return fmt.Sprintf(`
variable "listener_configuration" {
  type = list(object({
    protocol_port = number
    protocol      = string
  }))
  default = [
    {protocol_port = 306, protocol = "TCP"},
    {protocol_port = 406, protocol = "UDP"},
    {protocol_port = 506, protocol = "HTTP"},
  ]
}

resource "flexibleengine_vpc_v1" "test" {
  name = "%[1]s"
  cidr = "%[2]s"
}

resource "flexibleengine_vpc_subnet_v1" "test" {
  vpc_id = flexibleengine_vpc_v1.test.id

  name       = "%[1]s"
  cidr       = cidrsubnet(flexibleengine_vpc_v1.test.cidr, 4, 1)
  gateway_ip = cidrhost(cidrsubnet(flexibleengine_vpc_v1.test.cidr, 4, 1), 1)
}

resource "flexibleengine_lb_loadbalancer" "test" {
  name          = "%[1]s"
  vip_subnet_id = flexibleengine_vpc_subnet_v1.test.ipv4_subnet_id
}

resource "flexibleengine_lb_listener" "test" {
  count = length(var.listener_configuration)

  loadbalancer_id = flexibleengine_lb_loadbalancer.test.id

  name          = "%[1]s-${count.index}"
  protocol      = var.listener_configuration[count.index]["protocol"]
  protocol_port = var.listener_configuration[count.index]["protocol_port"]
}
`, rName, rCidr)
}

func testAccDatasourceListeners_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_lb_listeners" "by_name" {
  depends_on = [flexibleengine_lb_listener.test]

  name = flexibleengine_lb_listener.test[0].name
}

data "flexibleengine_lb_listeners" "by_protocol" {
  depends_on = [flexibleengine_lb_listener.test]

  protocol = flexibleengine_lb_listener.test[1].protocol
}

data "flexibleengine_lb_listeners" "by_protocol_port" {
  depends_on = [flexibleengine_lb_listener.test]

  protocol_port = flexibleengine_lb_listener.test[2].protocol_port
}

output "name_query_result_validation" {
  value = contains(data.flexibleengine_lb_listeners.by_name.listeners[*].id,
  flexibleengine_lb_listener.test[0].id) && !contains(data.flexibleengine_lb_listeners.by_name.listeners[*].id,
  flexibleengine_lb_listener.test[1].id) && !contains(data.flexibleengine_lb_listeners.by_name.listeners[*].id,
  flexibleengine_lb_listener.test[2].id)
}

output "protocol_query_result_validation" {
  value = contains(data.flexibleengine_lb_listeners.by_protocol.listeners[*].id,
  flexibleengine_lb_listener.test[1].id) && !contains(data.flexibleengine_lb_listeners.by_protocol.listeners[*].id,
  flexibleengine_lb_listener.test[0].id) && !contains(data.flexibleengine_lb_listeners.by_protocol.listeners[*].id,
  flexibleengine_lb_listener.test[2].id)
}

output "protocol_port_query_result_validation" {
  value = contains(data.flexibleengine_lb_listeners.by_protocol_port.listeners[*].id,
  flexibleengine_lb_listener.test[2].id) && !contains(data.flexibleengine_lb_listeners.by_protocol_port.listeners[*].id,
  flexibleengine_lb_listener.test[0].id) && !contains(data.flexibleengine_lb_listeners.by_protocol_port.listeners[*].id,
  flexibleengine_lb_listener.test[1].id)
}
`, testAccDatasourceListeners_base(rName))
}
