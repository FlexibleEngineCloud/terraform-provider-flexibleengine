---
subcategory: "Cloud Trace Service (CTS)"
---

# flexibleengine_cts_tracker_v1

CTS Tracker data source allows access of Cloud Tracker.

## Example Usage

```hcl
variable "bucket_name" { }

data "flexibleengine_cts_tracker_v1" "tracker_v1" {
  bucket_name = var.bucket_name
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
