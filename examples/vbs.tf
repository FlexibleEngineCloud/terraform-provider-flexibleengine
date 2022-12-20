resource "flexibleengine_vbs_backup_policy_v2" "vbs" {
  name                = "${var.project}-vbs-policy"
  start_time          = "12:00"
  status              = "ON"
  retain_first_backup = "N"
  rentention_num      = 2
  frequency           = 1
  resources           = flexibleengine_blockstorage_volume_v2.volume.*.id
}

resource "flexibleengine_blockstorage_volume_v2" "volume" {
  count = var.disk_size_gb > 0 ? var.instance_count : 0
  name  = "${var.project}-disk${format("%02d", count.index + 1)}"
  size  = var.disk_size_gb
}

# resource "flexibleengine_vbs_backup_v2" "backups_1" {
#   count     = var.disk_size_gb > 0 ? var.instance_count : 0
#   volume_id = flexibleengine_blockstorage_volume_v2.volume.*.id[count.index]
#   name      = "${var.project}-${flexibleengine_blockstorage_volume_v2.volume.*.name[count.index]}-backup"
# }
