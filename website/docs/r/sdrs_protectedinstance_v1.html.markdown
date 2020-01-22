---
layout: "flexibleengine"
page_title: "FlexibleEngine: flexibleengine_sdrs_protectedinstance_v1"
sidebar_current: "docs-flexibleengine-resource-sdrs-protectedinstance-v1"
description: |-
  Manages a V1 SDRS protected instance resource within FlexibleEngine.
---

# flexibleengine_sdrs_protectedinstance_v1

Manages a SDRS protected instance resource within FlexibleEngine.

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

resource "flexibleengine_sdrs_protectedinstance_v1" "instance_1" {
  group_id = flexibleengine_sdrs_protectiongroup_v1.group_1.id
  server_id = "{{ server_id }}"
  name = "instance_1"
  description = "test description"
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of a protected instance.

* `description` - (Optional) The description of a protected instance. Changing this creates a new instance.

* `group_id` - (Required) Specifies the ID of the protection group where a protected instance is added. Changing this creates a new instance.

* `server_id` - (Required) Specifies the ID of the source server. Changing this creates a new instance.

* `cluster_id` - (Optional) Specifies the ID of a storage pool. Changing this creates a new instance.

* `primary_subnet_id` - (Optional) Specifies the subnet ID of the primary NIC on the target server. Changing this creates a new instance.

* `primary_ip_address` - (Optional) Specifies the IP address of the primary NIC on the target server. Changing this creates a new instance.

* `delete_target_server` - (Optional) Specifies whether to delete the target server. The default value is false.. Changing this creates a new instance.

* `delete_target_eip` - (Optional) Specifies whether to delete the EIP of the target server. The default value is false. Changing this creates a new instance.


## Attributes Reference

The following attributes are exported:

* `id` -  ID of the protected instance.
* `target_server` -  ID of the target server.
