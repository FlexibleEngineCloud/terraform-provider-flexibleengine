---
subcategory: "Cloud Trace Service (CTS)"
description: ""
page_title: "flexibleengine_cts_tracker_v1"
---

# flexibleengine_cts_tracker_v1

Allows you to collect, store, and query cloud resource operation records.

## Example Usage

 ```hcl
 variable "bucket_name" { }
 
 resource "flexibleengine_cts_tracker_v1" "tracker_v1" {
  bucket_name      = var.bucket_name
  file_prefix_name = "tracker"
 }

 ```

## Argument Reference

The following arguments are supported:

* `bucket_name` - (Required) The OBS bucket name for a tracker.

* `file_prefix_name` - (Optional) The prefix of a log that needs to be stored in an OBS bucket.

* `status` - The status of a tracker. The value should be **enabled** when creating a tracker,
  and can be enabled or disabled when updating it.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `tracker_name` - The tracker name. Currently, only tracker **system** is available.

## Import

CTS tracker can be imported using  `tracker_name`, e.g.

```
$ terraform import flexibleengine_cts_tracker_v1.tracker system
```
