---
subcategory: "Distributed Message Service (DMS)"
description: ""
page_title: "flexibleengine_dms_kafka_topic"
---

# flexibleengine_dms_kafka_topic

Manages a DMS Kafka topic resource within FlexibleEngine.

## Example Usage

```hcl
variable "kafka_instance_id" {}

resource "flexibleengine_dms_kafka_topic" "topic" {
  instance_id = var.kafka_instance_id
  name       = "topic_1"
  partitions = 20
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the DMS Kafka topic resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the DMS Kafka instance to which the topic belongs.
  Changing this creates a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the topic. The name starts with a letter,
  consists of 4 to 64 characters, and supports only letters, digits, hyphens (-) and underscores (_).
  Changing this creates a new resource.

* `partitions` - (Optional, Int, ForceNew) Specifies the partition number.
  The value ranges from 1 to 50 and defaults to 3. Changing this creates a new resource.

* `replicas` - (Optional, Int, ForceNew) Specifies the replica number. The value ranges from 1 to 3 and defaults to 3.
  Changing this creates a new resource.

* `aging_time` - (Optional, Int, ForceNew) Specifies the aging time in hours.
  The value ranges from 1 to 720 and defaults to 72. Changing this creates a new resource.

* `sync_replication` - (Optional, Bool, ForceNew) Whether or not to enable synchronous replication.
  Changing this creates a new resource.

* `sync_flushing` - (Optional, Bool, ForceNew) Whether or not to enable synchronous flushing.
  Changing this creates a new resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which equals to the topic name.

## Import

DMS Kafka topics can be imported using the Kafka instance ID and topic name separated by a slash, e.g.:

```sh
terraform import flexibleengine_dms_kafka_topic.topic c8057fe5-23a8-46ef-ad83-c0055b4e0c5c/topic_1
```
