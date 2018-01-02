resource "flexibleengine_compute_instance_v2" "webserver" {
  count           = "${var.instance_count}"
  name            = "${var.project}-webserver${format("%02d", count.index+1)}"
  image_name      = "${var.image_name}"
  flavor_name     = "${var.flavor_name}"
  key_pair        = "${flexibleengine_compute_keypair_v2.keypair.name}"
  security_groups = [
    "${flexibleengine_compute_secgroup_v2.secgrp_web.name}"
  ]

  network {
    uuid           = "${flexibleengine_networking_network_v2.network.id}"
  }
}

resource "flexibleengine_compute_volume_attach_v2" "volume_attach" {
  count       = "${var.disk_size_gb > 0 ? var.instance_count : 0}"
  instance_id = "${element(flexibleengine_compute_instance_v2.webserver.*.id, count.index)}"
  volume_id   = "${element(flexibleengine_blockstorage_volume_v2.volume.*.id, count.index)}"
}
