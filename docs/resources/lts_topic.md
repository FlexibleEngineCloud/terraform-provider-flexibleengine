---
subcategory: "Log Tank Service (LTS)"
description: ""
page_title: "flexibleengine_lts_topic"
---

# flexibleengine_lts_topic

Manage a log topic resource within FlexibleEngine.

## Example Usage

```hcl
resource "flexibleengine_lts_group" "test_group" {
  group_name = "test_group"
}

resource "flexibleengine_lts_topic" "test_topic" {
  group_id   = flexibleengine_lts_group.test_group.id
  topic_name = "test1"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the log topic resource.
  If omitted, the provider-level region will be used. Changing this will create a new log topic resource.

* `group_id` - (Required, String, ForceNew) Specifies the ID of a created log group.
  Changing this parameter will create a new resource.

* `topic_name` - (Required, String, ForceNew) Specifies the log topic name.
  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The log topic ID in UUID format.

* `filter_count` - The Number of metric filter.

## Import

Log topic can be imported using the group ID and topic ID separated by a slash, e.g.

```sh
terraform import flexibleengine_lts_topic.topic_1 393f2bfd-2244-11ea-adb7-286ed488c87f/137159d3-e3b7-11eb-b952-286ed488cb76
```
