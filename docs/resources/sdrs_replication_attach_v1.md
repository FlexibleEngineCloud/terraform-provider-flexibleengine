---
subcategory: "Storage Disaster Recovery Service (SDRS)"
---

# flexibleengine_sdrs_replication_attach_v1

Manages a SDRS replication attch resource within FlexibleEngine.

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

resource "flexibleengine_sdrs_protectedinstance_v1" "instance_1" {
  group_id = flexibleengine_sdrs_protectiongroup_v1.group_1.id
  server_id = "{{ server_id }}"
  name = "instance_1"
  description = "test description"
}

resource "flexibleengine_sdrs_replication_pair_v1" "replication_1" {
  name        = "replication_1"
  description = "test description"
  group_id    = flexibleengine_sdrs_protectiongroup_v1.group_1.id
  volume_id   = "{{ volume_id }}"
}

resource "flexibleengine_sdrs_replication_attach_v1" "attach_1" {
  instance_id = flexibleengine_sdrs_protectedinstance_v1.instance_1.id
  replication_id = flexibleengine_sdrs_replication_pair_v1.replication_1.id
  device = "/dev/vdb"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) Specifies the ID of a protected instance. Changing this creates a new replication attach.

* `replication_id` - (Required) Specifies the ID of a replication pair. Changing this creates a new replication attach.

* `device` - (Required) Specifies the device name, eg. /dev/vdb. Changing this creates a new replication attach.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format of <instance_id>:<replication_id>.

* `status` - The status of the SDRS replication attch resource.
