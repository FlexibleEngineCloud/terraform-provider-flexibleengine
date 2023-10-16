---
subcategory: "Elastic Load Balance (ELB)"
description: ""
page_title: "flexibleengine_lb_member_v2"
---

# flexibleengine_lb_member_v2

Manages an **enhanced** load balancer member resource within FlexibleEngine.

## Example Usage

```hcl
resource "flexibleengine_lb_member_v2" "example_member" {
  address       = "192.168.199.23"
  protocol_port = 8080
  pool_id       = flexibleengine_lb_pool_v2.example_pool.id
  subnet_id     = flexibleengine_vpc_subnet_v1.example_subnet.ipv4_subnet_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to obtain the V2 Networking client.
  A Networking client is needed to be created. If omitted, the `region` argument of the provider is used.
  Changing this creates a new member.

* `pool_id` - (Required, String, ForceNew) The id of the pool that this member will be
  assigned to. Changing this creates a new member.

* `subnet_id` - (Required, String, ForceNew) The `ipv4_subnet_id` or `ipv6_subnet_id` of the
  VPC Subnet in which to access the member. Changing this creates a new member.

* `address` - (Required, String, ForceNew) The IP address of the member to receive traffic from
  the load balancer. Changing this creates a new member.

* `protocol_port` - (Required, Int, ForceNew) The port on which to listen for client traffic.
  Changing this creates a new member.

* `name` - (Optional, String) Human-readable name for the member.

* `weight` - (Optional, Int)  A positive integer value that indicates the relative
  portion of traffic that this member should receive from the pool. For
  example, a member with a weight of 10 receives five times as much traffic
  as a member with a weight of 2.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID for the member.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.
