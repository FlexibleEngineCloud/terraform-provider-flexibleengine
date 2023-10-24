---
subcategory: "Storage Disaster Recovery Service (SDRS)"
description: ""
page_title: "flexibleengine_sdrs_protectiongroup_v1"
---

# flexibleengine_sdrs_protectiongroup_v1

Manages a SDRS protection group resource within FlexibleEngine.

## Example Usage

```hcl
data "flexibleengine_sdrs_domain_v1" "domain_1" {
  name = "SDRS_HypeDomain01"
}

resource "flexibleengine_vpc_v1" "example_vpc" {
  name = "example-vpc"
  cidr = "192.168.0.0/16"
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

```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) The name of a protection group.

* `description` - (Optional, String, ForceNew) The description of a protection group. Changing this creates a new group.

* `source_availability_zone` - (Required, String, ForceNew) Specifies the source AZ of a protection group.
  Changing this creates a new group.

* `target_availability_zone` - (Required, String, ForceNew) Specifies the target AZ of a protection group.
  Changing this creates a new group.

* `domain_id` - (Required, String, ForceNew) Specifies the ID of an active-active domain.
  Changing this creates a new group.

* `source_vpc_id` - (Required, String, ForceNew) Specifies the ID of the source VPC.
  Changing this creates a new group.

* `dr_type` - (Optional, String, ForceNew) Specifies the deployment model. The default value is migration indicating
  migration within a VPC. Changing this creates a new group.

* `enable` - (Optional, Bool) Enable protection or not. It can only be set to true when there's replication pairs within
  the protection group.

## Attribute Reference

The following attributes are exported:

* `id` -  ID of the protection group.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

Protection groups can be imported using the `id`, e.g.

```shell
terraform import flexibleengine_sdrs_protectiongroup_v1.group_1 7117d38e-4c8f-4624-a505-bd96b97d024c
```
