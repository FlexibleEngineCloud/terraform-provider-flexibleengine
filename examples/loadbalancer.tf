resource "flexibleengine_lb_loadbalancer" "loadbalancer" {
  name           = "${var.project}-loadbalancer"
  vip_subnet_id  = flexibleengine_vpc_subnet_v1.subnet.ipv4_subnet_id
  admin_state_up = "true"
  depends_on     = [flexibleengine_vpc_subnet_v1.subnet]
}

resource "flexibleengine_lb_listener" "listener" {
  name            = "${var.project}-listener"
  protocol        = "HTTP"
  protocol_port   = 80
  loadbalancer_id = flexibleengine_lb_loadbalancer.loadbalancer.id
}

resource "flexibleengine_lb_pool" "pool" {
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = flexibleengine_lb_listener.listener.id
}

resource "flexibleengine_lb_member" "member" {
  count         = var.instance_count
  address       = element(flexibleengine_compute_instance_v2.webserver.*.access_ip_v4, count.index)
  pool_id       = flexibleengine_lb_pool.pool.id
  subnet_id     = flexibleengine_vpc_subnet_v1.subnet.ipv4_subnet_id
  protocol_port = 80
}

resource "flexibleengine_lb_monitor" "monitor" {
  pool_id        = flexibleengine_lb_pool.pool.id
  type           = "HTTP"
  url_path       = "/"
  expected_codes = "200"
  delay          = 20
  timeout        = 10
  max_retries    = 5
}
