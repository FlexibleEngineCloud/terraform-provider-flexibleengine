resource "orangecloud_lb_loadbalancer_v2" "loadbalancer" {
  count          = "${var.instance_count}"
  name           = "${var.project}-loadbalancer"
  vip_subnet_id  = "${orangecloud_networking_subnet_v2.subnet.id}"
  admin_state_up = "true"
  depends_on     = ["orangecloud_networking_router_interface_v2.interface"]
}

resource "orangecloud_lb_listener_v2" "listener" {
  name             = "${var.project}-listener"
  count          = "${var.instance_count}"
  protocol         = "HTTP"
  protocol_port    = 80
  loadbalancer_id  = "${orangecloud_lb_loadbalancer_v2.loadbalancer.id}"
  admin_state_up   = "true"
  #connection_limit = "-1"
}

resource "orangecloud_lb_pool_v2" "pool" {
  protocol    = "HTTP"
  count          = "${var.instance_count}"
  lb_method   = "ROUND_ROBIN"
  listener_id = "${orangecloud_lb_listener_v2.listener.id}"
}

resource "orangecloud_lb_member_v2" "member" {
  count         = "${var.instance_count}"
  address       = "${element(orangecloud_compute_instance_v2.webserver.*.access_ip_v4, count.index)}"
  pool_id       = "${orangecloud_lb_pool_v2.pool.id}"
  subnet_id     = "${orangecloud_networking_subnet_v2.subnet.id}"
  protocol_port = 80
}

resource "orangecloud_lb_monitor_v2" "monitor" {
  pool_id        = "${orangecloud_lb_pool_v2.pool.id}"
  count          = "${var.instance_count}"
  type           = "HTTP"
  url_path       = "/"
  expected_codes = "200"
  delay          = 20
  timeout        = 10
  max_retries    = 5
}
