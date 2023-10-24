---
subcategory: "Storage Disaster Recovery Service (SDRS)"
description: ""
page_title: "flexibleengine_sdrs_protectedinstance_v1"
---

# flexibleengine_sdrs_protectedinstance_v1

Manages a SDRS protected instance resource within FlexibleEngine.

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

resource "flexibleengine_sdrs_protectedinstance_v1" "instance_1" {
  group_id = flexibleengine_sdrs_protectiongroup_v1.group_1.id
  server_id = "{{ server_id }}"
  name = "instance_1"
  description = "test description"
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) The name of a protected instance.

* `description` - (Optional, String, ForceNew) The description of a protected instance. Changing this creates a new
  instance.

* `group_id` - (Required, String, ForceNew) Specifies the ID of the protection group where a protected instance is
  added. Changing this creates a new instance.

* `server_id` - (Required, String, ForceNew) Specifies the ID of the source server. Changing this creates a new instance.

* `cluster_id` - (Optional, String, ForceNew) Specifies the ID of a storage pool. Changing this creates a new instance.

* `primary_subnet_id` - (Optional, String, ForceNew) Specifies the `ipv4_subnet_id` or `ipv6_subnet_id` of the
  VPC Subnet of the primary NIC on the target server. Changing this creates a new instance.

* `primary_ip_address` - (Optional, String, ForceNew) Specifies the IP address of the primary NIC on the target server.
  Changing this creates a new instance.

* `delete_target_server` - (Optional, Bool, ForceNew) Specifies whether to delete the target server. The default
  value is false. Changing this creates a new instance.

* `delete_target_eip` - (Optional, Bool, ForceNew) Specifies whether to delete the EIP of the target server.
  The default value is false. Changing this creates a new instance.

## Attribute Reference

The following attributes are exported:

* `id` -  ID of the protected instance.

* `target_server` -  ID of the target server.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.
