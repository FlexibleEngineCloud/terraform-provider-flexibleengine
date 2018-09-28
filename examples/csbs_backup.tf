resource "flexibleengine_csbs_backup_v1" "backup_v1" {
  backup_name      = "${var.project}-backup"
  resource_id      = "${flexibleengine_compute_instance_v2.webserver.id}"
  resource_type    = "OS::Nova::Server"
}