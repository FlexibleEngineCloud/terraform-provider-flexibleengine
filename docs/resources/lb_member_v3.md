---
subcategory: "Elastic Load Balance (Dedicated ELB)"
description: ""
page_title: "flexibleengine_lb_member_v3"
---

# flexibleengine_lb_member_v3

Manages an ELB member resource within FlexibleEngine.

## Example Usage

```hcl
variable "address" {}
variable "pool_id" {}
variable "subnet_id" {}

resource "flexibleengine_lb_member_v3" "member_1" {
  address       = var.address
  protocol_port = 8080
  pool_id       = var.pool_id
  subnet_id     = var.subnet_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the ELB member resource.
  Changing this creates a new member.

* `pool_id` - (Required, String, ForceNew) Specifies the id of the pool that this member will be assigned to.

* `subnet_id` - (Optional, String, ForceNew) Specifies the subnet in which to access the member.
  The IPv4 or IPv6 subnet must be in the same VPC as the subnet of the load balancer.
  If this parameter is not passed, cross-VPC backend has been enabled for the load balancer. In this case,
  cross-VPC backend servers must use private IPv4 addresses, and the protocol of the backend server group
  must be TCP, HTTP, or HTTPS.

* `name` - (Optional, String) Specifies the name for the member.

* `address` - (Required, String, ForceNew) Specifies the IP address of the member to receive traffic from the
  load balancer. Changing this creates a new member.

* `protocol_port` - (Required, Int, ForceNew) Specifies the port on which to listen for client traffic.
  Changing this creates a new member.

* `weight` - (Optional, Int)  Specifies the positive integer value that indicates the relative portion of traffic
  that this member should receive from the pool. For example, a member with a weight of 10 receives five times as
  much traffic as a member with a weight of 2.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID for the member.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

ELB member can be imported using the pool ID and member ID separated by a slash, e.g.

```
$ terraform import flexibleengine_lb_member_v3.member_1 5c20fdad-7288-11eb-b817-0255ac10158b/e0bd694a-abbe-450e-b329-0931fd1cc5eb
```
