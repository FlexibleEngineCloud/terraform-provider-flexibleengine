resource "flexibleengine_vpc_eip" "eip_1" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "test"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "flexibleengine_vpc_eip_associate" "associated" {
  public_ip  = flexibleengine_vpc_eip.eip_1.address
  fixed_ip   = flexibleengine_lb_loadbalancer_v2.loadbalancer.vip_address
  network_id = flexibleengine_vpc_subnet_v1.subnet.id
}

resource "flexibleengine_antiddos_v1" "antiddos_1" {
  floating_ip_id         = flexibleengine_vpc_eip.eip_1.id
  enable_l7              = true
  traffic_pos_id         = 1
  http_request_pos_id    = 2
  cleaning_access_pos_id = 1
  app_type_id            = 0
}
