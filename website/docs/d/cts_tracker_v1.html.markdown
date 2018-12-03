---
layout: "flexibleengine"
page_title: "FlexibleEngine: flexibleengine_cts_tracker_v1"
sidebar_current: "docs-flexibleengine-datasource-cts-tracker-v1"
description: |-
  CTS tracker allows you to collect, store, and query cloud resource operation records and use these records for security analysis, compliance auditing, resource tracking, and fault locating.
---

# Data Source: flexibleengine_cts_tracker_v1

CTS Tracker data source allows access of Cloud Tracker.

## Example Usage


```hcl
variable "bucket_name" { }

data "flexibleengine_cts_tracker_v1" "tracker_v1" {
  bucket_name = "${var.bucket_name}"
}

```

## Argument Reference
The following arguments are supported:

* `tracker_name` - (Optional) The tracker name. 

* `bucket_name` - (Optional) The OBS bucket name for a tracker.

* `file_prefix_name` - (Optional) The prefix of a log that needs to be stored in an OBS bucket. 

* `status` - (Optional) Status of a tracker. 


## Attributes Reference

All above argument parameters can be exported as attribute parameters.
    