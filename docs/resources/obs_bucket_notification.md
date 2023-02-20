---
subcategory: "Object Storage Service (OSS)"
description: ""
page_title: "flexibleengine_obs_bucket_notification"
---

# flexibleengine_obs_bucket_notification

Manages an OBS bucket **Notification Configuration** resource within FlexibleEngine.

[Notification Configuration](https://docs.prod-cloud-ocb.orange-business.com/usermanual/obs/en-us_topic_0045853816.html)
OBS leverages SMN to provide the event notification function. In OBS, you can use SMN to send event notifications to
specified subscribers, so that you will be informed of any critical operations (such as upload and deletion) 
that occur on specified buckets in real time.

## Example Usage

### OBS Notification Configuration

```hcl
resource "flexibleengine_obs_bucket" "bucket" {
  bucket = "my-test-bucket"
  acl = "public-read"
}

resource "flexibleengine_obs_bucket_notification" notification {
  bucket = flexibleengine_obs_bucket.bucket.bucket
  topic_configurations {
    name      = "notification_name"
    events    = ["ObjectRemoved:*"]
    prefix    = "tf_update"
    suffix    = ".png"
    topic_urn = "urn:smn:eu-west-0:d8cb0fdcf29b4badb9ed8b2525a3286f:topic"
  }
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required) Specifies the name of the source bucket. Changing this parameter will create a new resource.

* `name` (Required) Specifies the name of OBS Notification. The name must be unique.

* `events` (Required) Type of events that need to be notified. The events include `ObjectCreated:*`,
  `ObjectCreated:Put`, `ObjectCreated:Post`, `ObjectCreated:Copy`, `ObjectCreated:CompleteMultipartUpload`,
  `ObjectRemoved:*`, `ObjectRemoved:Delete`, `ObjectRemoved:DeleteMarkerCreated`.

* `prefix` (Optional) Specifies the prefix filtering rule. The value contains a maximum of 1024 characters.

* `suffix` (Optional) Specifies the suffix filtering rule. The value contains a maximum of 1024 characters.

* `topic_urn` (Required) Specifies the SMN topic that authorizes OBS to publish messages.

## Attributes Reference

The following attributes are exported:

* `id` - The name of the bucket.

## Import

OBS bucket notification configuration can be imported using the `<bucket>` and "name" separated by a slash, e.g.:

```shell
terraform import flexibleengine_obs_bucket_notification.instance <bucket>/name
```
