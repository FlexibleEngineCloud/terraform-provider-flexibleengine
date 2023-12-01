---
subcategory: "Virtual Private Cloud (VPC)"
description: ""
page_title: "flexibleengine_vpc_flow_log_v1"
---

# flexibleengine_vpc_flow_log_v1

Manages a VPC flow log resource within FlexibleEngine.

## Example Usage

```hcl
resource "flexibleengine_lts_group" "log_group1" {
  group_name = var.log_group_name
}

resource "flexibleengine_lts_topic" "log_topic1" {
  group_id   = flexibleengine_lts_group.log_group1.id
  topic_name = var.log_topic_name
}

resource "flexibleengine_vpc_flow_log_v1" "flowlog1" {
  name          = var.flowlog_name
  description   = var.flowlog_desc
  resource_id   = var.port_id
  traffic_type  = "all"
  log_group_id  = flexibleengine_lts_group.log_group1.id
  log_topic_id  = flexibleengine_lts_topic.log_topic1.id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the VPC flow log resource.
  If omitted, the provider-level region will be used. Changing this creates a new VPC flow log.

* `name` - (Required, String) Specifies the VPC flow log name. The value is a string of 1 to 64 characters
  that can contain letters, digits, underscores (_), hyphens (-) and periods (.).

* `resource_id` - (Required, String, ForceNew) Specifies the network port ID.
  Changing this creates a new VPC flow log.

* `log_group_id` - (Required, String, ForceNew) Specifies the LTS log group ID.
  Changing this creates a new VPC flow log.

* `log_topic_id` - (Required, String, ForceNew) Specifies the LTS log topic ID.
  Changing this creates a new VPC flow log.

* `traffic_type` - (Optinal, String, ForceNew) Specifies the type of traffic to log. The value can be:
  - *all*: specifies that both accepted and rejected traffic of the specified resource will be logged.
  - *accept*: specifies that only accepted inbound and outbound traffic of the specified resource will be logged.
  - *reject*: specifies that only rejected inbound and outbound traffic of the specified resource will be logged.

  Defauts to *all*. Changing this creates a new VPC flow log.

* `description` - (Optinal, String) Specifies supplementary information about the VPC flow log.
  The value is a string of no more than 255 characters and cannot contain angle brackets (< or >).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The VPC flow log ID in UUID format.

* `resource_type` - The type of resource on which to create the VPC flow log. The value is fixed to *port*.

* `status` - The status of the flow log. The value can be `ACTIVE`, `DOWN` or `ERROR`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minute.
* `delete` - Default is 10 minute.

## Import

VPC flow logs can be imported using the `id`, e.g.

```shell
terraform import flexibleengine_vpc_flow_log_v1.flowlog1 41b9d73f-eb1c-4795-a100-59a99b062513
```
