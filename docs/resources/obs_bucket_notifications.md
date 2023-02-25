---
subcategory: "Object Storage Service (OBS)"
description: ""
page_title: "flexibleengine_obs_bucket_notifications"
---

# flexibleengine_obs_bucket_notifications

Manages an OBS bucket **Notification Configuration** resource within FlexibleEngine.

**The resource overwrites an existing configuration**.

[Notification Configuration](https://docs.prod-cloud-ocb.orange-business.com/usermanual/obs/en-us_topic_0045853816.html)
OBS leverages SMN to provide the event notification function. In OBS, you can use SMN to send event notifications to
specified subscribers, so that you will be informed of any critical operations (such as upload and deletion)
that occur on specified buckets in real time.

## Example Usage

### OBS Notification Configuration

```hcl
resource "flexibleengine_obs_bucket" "bucket" {
  bucket = "my-test-bucket"
  acl    = "public-read"
}

resource "flexibleengine_obs_bucket_notifications" notification {
  bucket = flexibleengine_obs_bucket.bucket.bucket

  notifications {
    name      = "notification_name"
    events    = ["ObjectCreated:*"]
    prefix    = "tf"
    suffix    = ".jpg"
    topic_urn = "urn:smn:eu-west-0:d8cb0fdcf29b4badb9ed8b2525a3286f:topic"
  }
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, String, ForceNew) Specifies the name of the source bucket.
  Changing this parameter will create a new resource.

* `notifications` - (Optional, List) Specifies the list of OBS bucket Notification Configurations.

The `notifications` block supports:

* `topic_urn` (Required, String) Specifies the SMN topic that authorizes OBS to publish messages.

* `events` (Required, List) Type of events that need to be notified. The events include `ObjectCreated:*`,
  `ObjectCreated:Put`, `ObjectCreated:Post`, `ObjectCreated:Copy`, `ObjectCreated:CompleteMultipartUpload`,
  `ObjectRemoved:*`, `ObjectRemoved:Delete`, `ObjectRemoved:DeleteMarkerCreated`.

* `name` (Optional, String) Specifies the name of OBS Notification. If not specified, the system assigns an ID
  automatically.

* `prefix` (Optional, String) Specifies the prefix filtering rule. The value contains a maximum of 1024 characters.

* `suffix` (Optional, String) Specifies the suffix filtering rule. The value contains a maximum of 1024 characters.

## Attributes Reference

The following attributes are exported:

* `id` - The name of the bucket.

## Import

OBS bucket notification configuration can be imported using the `bucket`, e.g.

```shell
terraform import flexibleengine_obs_bucket_notifications.instance <bucket-name>
```
