---
subcategory: "Simple Message Notification (SMN)"
---

# flexibleengine_smn_topics

Use this data source to get an array of SMN topics.

## Example Usage

```hcl
variable "topic_name" {}

data "flexibleengine_smn_topics" "tpoic_1" {
  name = var.topic_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to obtain the SMN topics. If omitted, the
  provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the topic.

* `topic_urn` - (Optional, String) Specifies the topic URN.

* `display_name` - (Optional, String) Specifies the topic display name.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project id of the topic.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `topics` - An array of SMN topics found. The [topics](#smn_topics) object structure is documented below.

<a name="smn_topics"></a>
The `topics` block supports:

* `name` - The name of the topic.

* `id` - The topic ID. The value is the topic URN.

* `topic_urn` - The topic URN.

* `display_name` - The topic display name.

* `enterprise_project_id` - The enterprise project id of the topic.

* `push_policy` - Message pushing policy.
  + **0**: indicates that the message sending fails and the message is cached in the queue.
  + **1**: indicates that the failed message is discarded.
