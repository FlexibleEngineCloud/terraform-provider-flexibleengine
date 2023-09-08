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

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the CTS tracker resource.
  If omitted, the provider-level region will be used. Changing this will create a new CTS tracker resource.

* `bucket_name` - (Required, String) The OBS bucket name for a tracker.

* `file_prefix_name` - (Optional, String) The prefix of a log that needs to be stored in an OBS bucket.
  The value can contain letters, digits, and special characters `.-_`, but cannot contain spaces.
  The length is 1 to 64 characters.

* `status` - (Optional, String) The status of a tracker. The value should be **enabled** when creating a tracker,
  and can be enabled or disabled when updating it.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `tracker_name` - The tracker name. Currently, only tracker **system** is available.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default 5 minutes
* `delete` - Default 5 minutes

## Import

CTS tracker can be imported using  `tracker_name`, e.g.

```shell
terraform import flexibleengine_cts_tracker_v1.tracker system
```
