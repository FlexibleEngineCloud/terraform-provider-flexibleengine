---
layout: "flexibleengine"
page_title: "FlexibleEngine: flexibleengine_smn_topic_v2"
sidebar_current: "docs-flexibleengine-resource-smn-topic-v2"
description: |-
  Manages a V2 topic resource within FlexibleEngine.
---

# flexibleengine\_smn\_topic\_v2

Manages a V2 topic resource within FlexibleEngine.

## Example Usage

```hcl
resource "flexibleengine_smn_topic_v2" "topic_1" {
  name            = "topic_1"
  display_name    = "The display name of topic_1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the topic to be created.

* `display_name` - (Optional) Topic display name, which is presented as the
    name of the email sender in an email message.

## Attributes Reference

The following attributes are exported:

* `name` - See Argument Reference above.
* `display_name` - See Argument Reference above.
* `topic_urn` - Resource identifier of a topic, which is unique.
* `push_policy` - Message pushing policy. 0 indicates that the message
    sending fails and the message is cached in the queue. 1 indicates that the
    failed message is discarded.
* `create_time` - Time when the topic was created.
* `update_time` - Time when the topic was updated.
