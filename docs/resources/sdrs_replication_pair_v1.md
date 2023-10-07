---
subcategory: "Storage Disaster Recovery Service (SDRS)"
description: ""
page_title: "flexibleengine_sdrs_replication_pair_v1"
---

# flexibleengine_sdrs_replication_pair_v1

Manages a SDRS replication pair resource within FlexibleEngine.

## Example Usage

```hcl
resource "flexibleengine_vpc_v1" "example_vpc" {
  name = "example-vpc"
  cidr = "192.168.0.0/16"
}

data "flexibleengine_sdrs_domain_v1" "domain_1" {
  name = "SDRS_HypeDomain01"
}

resource "flexibleengine_sdrs_protectiongroup_v1" "group_1" {
  name        = "group_1"
  description = "test description"
  source_availability_zone = "eu-west-0a"
  target_availability_zone = "eu-west-0b"
  domain_id     = data.flexibleengine_sdrs_domain_v1.domain_1.id
  source_vpc_id = flexibleengine_vpc_v1.example_vpc.id
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

* `name` - (Required, String) The name of a replication pair. The name can contain a maximum of 64 bytes.
  The value can contain only letters (a to z and A to Z), digits (0 to 9), decimal points (.),
  underscores (_), and hyphens (-).

* `description` - (Optional, String, ForceNew) The description of a replication pair. Changing this creates a new pair.

* `group_id` - (Required, String, ForceNew) Specifies the ID of a protection group. Changing this creates a new pair.

* `volume_id` - (Required, String, ForceNew) Specifies the ID of a source disk. Changing this creates a new pair.

* `delete_target_volume` - (Optional, Bool) Specifies whether to delete the target disk.
  The default value is `false`.

## Attribute Reference

The following attributes are exported:

* `id` -  ID of the replication pair.

* `fault_level` - Specifies the fault level of a replication pair.

* `replication_model` - Specifies the replication mode of a replication pair. The default value is `hypermetro`.

* `status` - Specifies the status of a replication pair.

* `target_volume_id` - Specifies the ID of the disk in the protection availability zone.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

Replication pairs can be imported using the `id`, e.g.

```shell
terraform import flexibleengine_sdrs_replication_pair_v1.replication_1 43b28b66-770b-4e9e-b5c6-cfc43f0593d9
```
