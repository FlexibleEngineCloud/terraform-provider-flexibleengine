resource "flexibleengine_networking_secgroup_v2" "secgrp_web" {
  name        = "${var.project}-secgrp-web-elb"
  description = "Webserver Security Group"
}

resource "flexibleengine_networking_secgroup_rule_v2" "secgrp_web_rule_1" {
  security_group_id = flexibleengine_networking_secgroup_v2.secgrp_web.id
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  port_range_min    = 22
  port_range_max    = 22
}

resource "flexibleengine_networking_secgroup_rule_v2" "secgrp_web_rule_2" {
  security_group_id = flexibleengine_networking_secgroup_v2.secgrp_web.id
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  port_range_min    = 80
  port_range_max    = 80
}

resource "flexibleengine_networking_secgroup_rule_v2" "secgrp_web_rule_3" {
  security_group_id = flexibleengine_networking_secgroup_v2.secgrp_web.id
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "icmp"
  remote_ip_prefix  = "0.0.0.0/0"
}
