---
layout: "flexibleengine"
page_title: "FlexibleEngine: flexibleengine_sdrs_replication_pair_v1"
sidebar_current: "docs-flexibleengine-resource-sdrs-replication-pair-v1"
description: |-
  Manages a V1 SDRS replication pair resource within FlexibleEngine.
---

# flexibleengine_sdrs_replication_pair_v1

Manages a SDRS replication pair resource within FlexibleEngine.

## Example Usage

```hcl
data "flexibleengine_sdrs_domain_v1" "domain_1" {
  name = "SDRS_HypeDomain01"
}

resource "flexibleengine_sdrs_protectiongroup_v1" "group_1" {
  name        = "group_1"
  description = "test description"
  source_availability_zone = "eu-west-0a"
  target_availability_zone = "eu-west-0b"
  domain_id     = data.flexibleengine_sdrs_domain_v1.domain_1.id
  source_vpc_id = "{{ vpc_id }}"
  dr_type       = "migration"
}
resource "flexibleengine_sdrs_replication_pair_v1" "replication_1" {
  name        = "replication_1"
  description = "test description"
  group_id    = flexibleengine_sdrs_protectiongroup_v1.group_1.id
  volume_id   = "{{ volume_id }}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of a replication pair. The name can contain a maximum of 64 bytes.
  The value can contain only letters (a to z and A to Z), digits (0 to 9), decimal points (.),
  underscores (_), and hyphens (-).

* `description` - (Optional) The description of a replication pair. Changing this creates a new pair.

* `group_id` - (Required) Specifies the ID of a protection group. Changing this creates a new pair.

* `volume_id` - (Required) Specifies the ID of a source disk. Changing this creates a new pair.

* `delete_target_volume` - (Optional) Specifies whether to delete the target disk.
  The default value is `false`.


## Attributes Reference

The following attributes are exported:

* `id` -  ID of the replication pair.

* `fault_level` - Specifies the fault level of a replication pair.

* `replication_model` - Specifies the replication mode of a replication pair. The default value is `hypermetro`.

* `status` - Specifies the status of a replication pair.

* `target_volume_id` - Specifies the ID of the disk in the protection availability zone.

## Import

Replication pairs can be imported using the `id`, e.g.

```
$ terraform import flexibleengine_sdrs_replication_pair_v1.replication_1 43b28b66-770b-4e9e-b5c6-cfc43f0593d9
```
