resource "random_string" "bucket" {
  length           = 16
  upper            = false
  override_special = "-"
}

resource "flexibleengine_s3_bucket" "bucket" {
  bucket = "${var.project}-${random_string.bucket.result}-bucket"
  acl    = "public-read"
}

resource "flexibleengine_smn_topic_v2" "topic_1" {
  name         = "topic_check"
  display_name = "The display name of topic_1"
}

resource "flexibleengine_cts_tracker_v1" "tracker_v1" {
  bucket_name      = flexibleengine_s3_bucket.bucket.bucket
  file_prefix_name = "yO8Q"
}
