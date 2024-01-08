package acceptance

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/elb/v2/monitors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccLBV2Monitor_basic(t *testing.T) {
	var monitor monitors.Monitor
	rName := acceptance.RandomAccResourceNameWithDash()
	rNameUpdate := rName + "-update"
	resourceName := "flexibleengine_lb_monitor.monitor_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&monitor,
		getMonitorResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccLBV2MonitorConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "TCP"),
					resource.TestCheckResourceAttr(resourceName, "delay", "20"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "10"),
					resource.TestCheckResourceAttr(resourceName, "max_retries", "5"),
				),
			},
			{
				Config: testAccLBV2MonitorConfig_update(rName, rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "delay", "30"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "15"),
					resource.TestCheckResourceAttr(resourceName, "max_retries", "3"),
					resource.TestCheckResourceAttr(resourceName, "port", "8888"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLBV2Monitor_udp(t *testing.T) {
	var monitor monitors.Monitor
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "flexibleengine_lb_monitor.monitor_udp"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&monitor,
		getMonitorResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccLBV2MonitorConfig_udp(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "UDP_CONNECT"),
					resource.TestCheckResourceAttr(resourceName, "delay", "20"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "10"),
					resource.TestCheckResourceAttr(resourceName, "max_retries", "5"),
				),
			},
		},
	})
}

func TestAccLBV2Monitor_http(t *testing.T) {
	var monitor monitors.Monitor
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "flexibleengine_lb_monitor.monitor_http"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&monitor,
		getMonitorResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccLBV2MonitorConfig_http(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "delay", "20"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "10"),
					resource.TestCheckResourceAttr(resourceName, "max_retries", "5"),
					resource.TestCheckResourceAttr(resourceName, "url_path", "/api"),
					resource.TestCheckResourceAttr(resourceName, "http_method", "GET"),
					resource.TestCheckResourceAttr(resourceName, "expected_codes", "200-202"),
				),
			},
		},
	})
}

func testAccLBV2MonitorConfig_base(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_vpc_v1" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "test" {
  name       = "%[1]s"
  cidr       = "192.168.0.0/24"
  vpc_id     = flexibleengine_vpc_v1.test.id
  gateway_ip = "192.168.0.1"
}

resource "flexibleengine_lb_loadbalancer" "loadbalancer_1" {
  name          = "%[1]s"
  vip_subnet_id = flexibleengine_vpc_subnet_v1.test.ipv4_subnet_id
}

resource "flexibleengine_lb_listener" "listener_1" {
  name            = "%[1]s"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = flexibleengine_lb_loadbalancer.loadbalancer_1.id
}

resource "flexibleengine_lb_pool" "pool_1" {
  name        = "%[1]s"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = flexibleengine_lb_listener.listener_1.id
}
`, rName)
}

func testAccLBV2MonitorConfig_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_lb_monitor" "monitor_1" {
  pool_id     = flexibleengine_lb_pool.pool_1.id
  name        = "%s"
  type        = "TCP"
  delay       = 20
  timeout     = 10
  max_retries = 5
}
`, testAccLBV2MonitorConfig_base(rName), rName)
}

func testAccLBV2MonitorConfig_update(rName, rNameUpdate string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_lb_monitor" "monitor_1" {
  pool_id     = flexibleengine_lb_pool.pool_1.id
  name        = "%s"
  type        = "TCP"
  delay       = 30
  timeout     = 15
  max_retries = 3
  port        = 8888
}
`, testAccLBV2MonitorConfig_base(rName), rNameUpdate)
}

func testAccLBV2MonitorConfig_http(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_lb_monitor" "monitor_http" {
  pool_id        = flexibleengine_lb_pool.pool_1.id
  name           = "%s"
  type           = "HTTP"
  delay          = 20
  timeout        = 10
  max_retries    = 5
  url_path       = "/api"
  expected_codes = "200-202"
}
`, testAccLBV2MonitorConfig_base(rName), rName)
}

func testAccLBV2MonitorConfig_udp(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_vpc_v1" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "test" {
  name       = "%[1]s"
  cidr       = "192.168.0.0/24"
  vpc_id     = flexibleengine_vpc_v1.test.id
  gateway_ip = "192.168.0.1"
}

resource "flexibleengine_lb_loadbalancer" "loadbalancer_1" {
  name          = "%[1]s"
  vip_subnet_id = flexibleengine_vpc_subnet_v1.test.ipv4_subnet_id
}

resource "flexibleengine_lb_listener" "listener_1" {
  name            = "%[1]s"
  protocol        = "UDP"
  protocol_port   = 8080
  loadbalancer_id = flexibleengine_lb_loadbalancer.loadbalancer_1.id
}

resource "flexibleengine_lb_pool" "pool_1" {
  name        = "%[1]s"
  protocol    = "UDP"
  lb_method   = "ROUND_ROBIN"
  listener_id = flexibleengine_lb_listener.listener_1.id
}

resource "flexibleengine_lb_monitor" "monitor_udp" {
  pool_id     = flexibleengine_lb_pool.pool_1.id
  name        = "%[1]s"
  type        = "UDP_CONNECT"
  delay       = 20
  timeout     = 10
  max_retries = 5
}
`, rName)
}
