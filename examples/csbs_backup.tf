# resource "flexibleengine_csbs_backup_v1" "backup_v1" {
#   count         = var.instance_count
#   backup_name   = "${var.project}-${flexibleengine_compute_instance_v2.webserver.*.name[count.index]}-backup"
#   resource_id   = flexibleengine_compute_instance_v2.webserver.*.id[count.index]
#   resource_type = "OS::Nova::Server"
# }
