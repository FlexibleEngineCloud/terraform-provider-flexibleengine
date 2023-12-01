---
subcategory: "Virtual Private Cloud (VPC)"
description: ""
page_title: "flexibleengine_vpc_peering_connection_v2"
---

# flexibleengine_vpc_peering_connection_v2

Provides a resource to manage a VPC Peering Connection resource.

-> **Note:** If you create a VPC peering connection with another VPC of your own, the connection is created without
  the need for you to accept the connection.  
  For cross-tenant (requester's tenant differs from the accepter's tenant) VPC Peering Connections,
  use the `flexibleengine_vpc_peering_connection_v2` resource to manage the requester's side of the connection and
  use the `flexibleengine_vpc_peering_connection_accepter_v2` resource to manage the accepter's side of the connection.

## Example Usage

```hcl
variable "accepter_vpc_id" {}

resource "flexibleengine_vpc_v1" "example_vpc" {
  name = "example-vpc"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_peering_connection_v2" "peering" {
  name        = var.peer_conn_name
  vpc_id      = flexibleengine_vpc_v1.example_vpc.id
  peer_vpc_id = var.accepter_vpc_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the VPC peering connection.
  If omitted, the provider-level region will be used. Changing this creates a new VPC peering connection.

* `name` (Required, String) - Specifies the name of the VPC peering connection. The value can contain 1 to 64 characters.

* `vpc_id` (Required, String, ForceNew) - Specifies the ID of a VPC involved in a VPC peering connection.
  Changing this creates a new VPC peering connection.

* `peer_vpc_id` (Required, String, ForceNew) - Specifies the VPC ID of the accepter tenant.
  Changing this creates a new VPC peering connection.

* `peer_tenant_id` (Optional, String, ForceNew) - Specified the Tenant Id of the accepter tenant.
  Changing this creates a new VPC peering connection.

## Attributes Reference

All of the argument attributes are also exported as
result attributes:

* `id` - The VPC peering connection ID.

* `status` - The VPC peering connection status. The value can be PENDING_ACCEPTANCE, REJECTED, EXPIRED, DELETED, or ACTIVE.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minute.
* `delete` - Default is 10 minute.

## Import

VPC Peering resources can be imported using the `vpc peering id`, e.g.

```shell
terraform import flexibleengine_vpc_peering_connection_v2.test_connection 22b76469-08e3-4937-8c1d-7aad34892be1
```
