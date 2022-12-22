resource "flexibleengine_csbs_backup_policy_v1" "backup_policy_v1" {
  name = "${var.project}-csbs-policy"
  dynamic "resource" {
    for_each = flexibleengine_compute_instance_v2.webserver.*.id
    content {
      id   = resource.value
      type = "OS::Nova::Server"
      name = flexibleengine_compute_instance_v2.webserver.*.name[resource.key]
    }
  }

  scheduled_operation {
    name            = "mybackup"
    enabled         = true
    description     = "My backup policy"
    operation_type  = "backup"
    max_backups     = 2
    trigger_pattern = "BEGIN:VCALENDAR\r\nBEGIN:VEVENT\r\nRRULE:FREQ=WEEKLY;BYDAY=TH;BYHOUR=12;BYMINUTE=27\r\nEND:VEVENT\r\nEND:VCALENDAR\r\n"
  }
}
