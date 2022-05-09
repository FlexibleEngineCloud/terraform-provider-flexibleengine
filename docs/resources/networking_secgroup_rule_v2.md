---
subcategory: "Virtual Private Cloud (VPC)"
---

# flexibleengine_networking_secgroup_rule_v2

Manages a Security Group Rule resource within FlexibleEngine.

## Example Usage

```hcl
resource "flexibleengine_networking_secgroup_v2" "secgroup_1" {
  name        = "secgroup_1"
  description = "My neutron security group"
}

resource "flexibleengine_networking_secgroup_rule_v2" "secgroup_rule_1" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 22
  port_range_max    = 22
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = flexibleengine_networking_secgroup_v2.secgroup_1.id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 networking client.
    A networking client is needed to create a port. If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    security group rule.

* `security_group_id` - (Required) The security group ID the rule should belong
    to. Changing this creates a new security group rule.

* `direction` - (Required) The direction of the rule, valid values are __ingress__
    or __egress__. Changing this creates a new security group rule.

* `ethertype` - (Required) The layer 3 protocol type, valid values are __IPv4__
    or __IPv6__. Changing this creates a new security group rule.

* `protocol` - (Optional) The layer 4 protocol type, valid values are following.
    Changing this creates a new security group rule. This is required if you want to specify a port range.
  * __tcp__
  * __udp__
  * __icmp__
  * __ah__
  * __dccp__
  * __egp__
  * __esp__
  * __gre__
  * __igmp__
  * __ipv6-encap__
  * __ipv6-frag__
  * __ipv6-icmp__
  * __ipv6-nonxt__
  * __ipv6-opts__
  * __ipv6-route__
  * __ospf__
  * __pgm__
  * __rsvp__
  * __sctp__
  * __udplite__
  * __vrrp__

* `port_range_min` - (Optional) The lower part of the allowed port range, valid
    integer value needs to be between 1 and 65535. Changing this creates a new
    security group rule.

* `port_range_max` - (Optional) The higher part of the allowed port range, valid
    integer value needs to be between 1 and 65535. Changing this creates a new
    security group rule.

* `remote_ip_prefix` - (Optional) The remote CIDR, the value needs to be a valid
    CIDR (i.e. 192.168.0.0/16). Changing this creates a new security group rule.

* `remote_group_id` - (Optional) The remote group id, the value needs to be an
    FlexibleEngine ID of a security group in the same tenant. Changing this creates
    a new security group rule.

* `description` - (Optional) Specifies the supplementary information about the security group rule.
  This parameter can contain a maximum of 255 characters and cannot contain angle brackets (< or >).
  Changing this creates a new security group rule.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

## Import

Security Group Rules can be imported using the `id`, e.g.

```
$ terraform import flexibleengine_networking_secgroup_rule_v2.secgroup_rule_1 aeb68ee3-6e9d-4256-955c-9584a6212745
```
