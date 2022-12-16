---
subcategory: "Elastic Load Balance (ELB)"
description: ""
page_title: "flexibleengine_lb_member_v2"
---

# flexibleengine_lb_member_v2

Manages an **enhanced** load balancer member resource within FlexibleEngine.

## Example Usage

```hcl
resource "flexibleengine_lb_member_v2" "member_1" {
  address       = "192.168.199.23"
  protocol_port = 8080
  pool_id       = POOL_ID
  subnet_id     = SUBNET_ID
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Networking client.
    A Networking client is needed to create an . If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    member.

* `pool_id` - (Required) The id of the pool that this member will be
    assigned to.

* `subnet_id` - (Required) The subnet in which to access the member

* `name` - (Optional) Human-readable name for the member.

* `address` - (Required) The IP address of the member to receive traffic from
    the load balancer. Changing this creates a new member.

* `protocol_port` - (Required) The port on which to listen for client traffic.
    Changing this creates a new member.

* `weight` - (Optional)  A positive integer value that indicates the relative
    portion of traffic that this member should receive from the pool. For
    example, a member with a weight of 10 receives five times as much traffic
    as a member with a weight of 2.

* `admin_state_up` - (Optional) The administrative state of the member.
    A valid value is true (UP) or false (DOWN).

* `tenant_id` - (Optional) The UUID of the tenant who owns the member.
    Only administrative users can specify a tenant UUID other than their own.
    Changing this creates a new member.

## Attributes Reference

The following attributes are exported:

* `id` - The unique ID for the member.
* `name` - See Argument Reference above.
* `weight` - See Argument Reference above.
* `admin_state_up` - See Argument Reference above.
* `tenant_id` - See Argument Reference above.
* `subnet_id` - See Argument Reference above.
* `pool_id` - See Argument Reference above.
* `address` - See Argument Reference above.
* `protocol_port` - See Argument Reference above.
