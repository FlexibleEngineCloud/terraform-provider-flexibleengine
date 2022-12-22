---
subcategory: "Storage Disaster Recovery Service (SDRS)"
description: ""
page_title: "flexibleengine_sdrs_drill_v1"
---

# flexibleengine_sdrs_drill_v1

Manages a Disaster Recovery Drill resource within FlexibleEngine.

## Example Usage

```hcl
data "flexibleengine_sdrs_domain_v1" "domain_1" {
  name = "SDRS_HypeDomain01"
}

resource "flexibleengine_vpc_v1" "example_vpc" {
  name = "example-vpc"
  cidr = "192.168.0.0/24"
}

resource "flexibleengine_vpc_v1" "example_vpc_drill" {
  name = "example-vpc"
  cidr = "192.168.1.0/24"
}

resource "flexibleengine_sdrs_protectiongroup_v1" "group_1" {
  name = "group_1"
  description = "test description"
  source_availability_zone = "eu-west-0a"
  target_availability_zone = "eu-west-0b"
  domain_id = data.flexibleengine_sdrs_domain_v1.domain_1.id
  source_vpc_id = flexibleengine_vpc_v1.example_vpc.id
  dr_type = "migration"
}

resource "flexibleengine_sdrs_drill_v1" "drill_1" {
  name         = "drill_1"
  group_id     = flexibleengine_sdrs_protectiongroup_v1.group_1.id
  drill_vpc_id = flexibleengine_vpc_v1.example_vpc_drill.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of a DR drill. The name can contain a maximum of 64 bytes.
  The value can contain only letters (a to z and A to Z), digits (0 to 9), decimal points (.),
  underscores (_), and hyphens (-).

* `group_id` - (Required) Specifies the ID of a protection group. Changing this creates a new drill.

* `drill_vpc_id` - (Required) Specifies the ID used for a DR drill. Changing this creates a new drill.

## Attributes Reference

The following attributes are exported:

* `id` -  ID of a DR drill.
* `name` - See Argument Reference above.
* `group_id` - See Argument Reference above.
* `drill_vpc_id` - See Argument Reference above.
* `status` - The status of a DR drill.
  For details, see [DR Drill Status](https://docs.prod-cloud-ocb.orange-business.com/en-us/api/sdrs/en-us_topic_0126152933.html).

## Import

DR drill can be imported using the `id`, e.g.

```shell
terraform import flexibleengine_sdrs_drill_v1.drill_1 22fce838-4bfb-4a92-b9aa-fc80a583eb59
```
