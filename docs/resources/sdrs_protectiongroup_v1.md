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

resource "flexibleengine_sdrs_protectiongroup_v1" "group_1" {
  name = "group_1"
  description = "test description"
  source_availability_zone = "eu-west-0a"
  target_availability_zone = "eu-west-0b"
  domain_id = data.flexibleengine_sdrs_domain_v1.domain_1.id
  source_vpc_id = "{{ vpc_id }}"
  dr_type = "migration"
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of a protection group.

* `description` - (Optional) The description of a protection group. Changing this creates a new group.

* `source_availability_zone` - (Required) Specifies the source AZ of a protection group. Changing this creates a new group.

* `target_availability_zone` - (Required) Specifies the target AZ of a protection group. Changing this creates a new group.

* `domain_id` - (Required) Specifies the ID of an active-active domain. Changing this creates a new group.

* `source_vpc_id` - (Required) Specifies the ID of the source VPC. Changing this creates a new group.

* `dr_type` - (Optional) Specifies the deployment model. The default value is migration indicating migration within a VPC.
  Changing this creates a new group.

* `enable` - (Optional) Enable protection or not. It can only be set to true when there's replication pairs within
  the protection group.

## Attributes Reference

The following attributes are exported:

* `id` -  ID of the protection group.

## Import

Protection groups can be imported using the `id`, e.g.

```
$ terraform import flexibleengine_sdrs_protectiongroup_v1.group_1 7117d38e-4c8f-4624-a505-bd96b97d024c
```
