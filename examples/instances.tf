resource "flexibleengine_compute_instance_v2" "webserver" {
  count       = var.instance_count
  name        = "${var.project}-webserver${format("%02d", count.index + 1)}"
  image_name  = var.image_name
  flavor_name = var.flavor_name
  key_pair    = flexibleengine_compute_keypair_v2.keypair.name
  user_data   = "#cloud-config\npackage_update: true\npackages: ['nginx-light']\n"
  security_groups = [
    flexibleengine_networking_secgroup_v2.secgrp_web.name
  ]

  network {
    uuid = flexibleengine_vpc_subnet_v1.subnet.id
  }
}

resource "flexibleengine_compute_volume_attach_v2" "volume_attach" {
  count       = var.disk_size_gb > 0 ? var.instance_count : 0
  instance_id = element(flexibleengine_compute_instance_v2.webserver.*.id, count.index)
  volume_id   = element(flexibleengine_blockstorage_volume_v2.volume.*.id, count.index)
}
