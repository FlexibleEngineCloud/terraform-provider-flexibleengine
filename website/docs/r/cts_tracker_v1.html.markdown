---
layout: "flexibleengine"
page_title: "FlexibleEngine: resource_flexibleengine_cts_tracker_v1"
sidebar_current: "docs-flexibleengine-resource-cts-tracker-v1"
description: |-
   CTS tracker allows you to collect, store, and query cloud resource operation records and use these records for security analysis, compliance auditing, resource tracking, and fault locating.
---

# flexibleengine_cts_tracker_v1

Allows you to collect, store, and query cloud resource operation records.

## Example Usage

 ```hcl
 variable "bucket_name" { }
 
 resource "flexibleengine_cts_tracker_v1" "tracker_v1" {
   bucket_name      = "${var.bucket_name}"
   file_prefix_name      = "yO8Q"
 }

 ```
## Argument Reference
The following arguments are supported:

* `bucket_name` - (Required) The OBS bucket name for a tracker.

* `file_prefix_name` - (Optional) The prefix of a log that needs to be stored in an OBS bucket. 


## Attributes Reference
In addition to all arguments above, the following attributes are exported:

* `status` - The status of a tracker. The value is **enabled**.

* `tracker_name` - The tracker name. Currently, only tracker **system** is available.


## Import

CTS tracker can be imported using  `tracker_name`, e.g.

```
$ terraform import flexibleengine_cts_tracker_v1.tracker system
```




