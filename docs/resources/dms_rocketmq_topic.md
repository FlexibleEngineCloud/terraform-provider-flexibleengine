---
subcategory: "Distributed Message Service (DMS)"
---

# flexibleengine_dms_rocketmq_topic

Manages DMS RocketMQ topic resources within FlexibleEngine.

## Example Usage

```hcl
variable "instance_id" {}

resource "flexibleengine_dms_rocketmq_topic" "test" {
  instance_id = var.instance_id
  name        = "topic_test"
  queue_num   = 3
  permission  = "all"

  brokers {
    name = "broker-0"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the rocketMQ instance.
  Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the topic.
  Changing this parameter will create a new resource.

* `brokers` - (Required, List, ForceNew) Specifies the list of associated brokers of the topic.
  Changing this parameter will create a new resource.
  The [brokers](#dms_arg_brokers) structure is documented below.

* `queue_num` - (Optional, Int, ForceNew) Specifies the number of queues. Default to 3.
  Changing this parameter will create a new resource.

* `permission` - (Optional, String) Specifies the permissions of the topic.
  Value options: **all**, **sub**, **pub**. Default to all.

* `total_read_queue_num` - (Optional, Int) Specifies the total number of read queues.

* `total_write_queue_num` - (Optional, Int) Specifies the total number of write queues.

<a name="dms_arg_brokers"></a>
The `brokers` block supports:

* `name` - (Optional, String) Indicates the name of the broker.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `brokers` - Specifies the list of associated brokers of the topic.
  The [brokers](#dms_attr_brokers) structure is documented below.

<a name="dms_attr_brokers"></a>
  The `brokers` block supports:

* `read_queue_num` - Indicates the read queues number of the broker. It's useless when create a topic.

* `write_queue_num` - Indicates the read queues number of the broker. It's useless when create a topic.

## Import

The rocketmq topic can be imported using the rocketMQ instance ID and topic name separated by a slash, e.g.

```sh
terraform import flexibleengine_dms_rocketmq_topic.test c8057fe5-23a8-46ef-ad83-c0055b4e0c5c/topic_1
```
