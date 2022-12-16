---
subcategory: "Network ACL"
---

# flexibleengine_network_acl

Manages a network ACL resource within FlexibleEngine.

## Example Usage

```hcl
data "flexibleengine_vpc_subnet_v1" "subnet" {
  name = "subnet-default"
}

resource "flexibleengine_network_acl_rule" "rule_1" {
  name             = "my-rule-1"
  description      = "drop TELNET traffic"
  action           = "deny"
  protocol         = "tcp"
  destination_port = "23"
  enabled          = "true"
}

resource "flexibleengine_network_acl_rule" "rule_2" {
  name             = "my-rule-2"
  description      = "drop NTP traffic"
  action           = "deny"
  protocol         = "udp"
  destination_port = "123"
  enabled          = "false"
}

resource "flexibleengine_network_acl" "fw_acl" {
  name          = "my-fw-acl"
  subnets       = [data.flexibleengine_vpc_subnet_v1.subnet.id]
  inbound_rules = [flexibleengine_network_acl_rule.rule_1.id,
    flexibleengine_network_acl_rule.rule_2.id]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the network ACL name. This parameter can contain a maximum of 64 characters,
    which may consist of letters, digits, underscores (_), and hyphens (-).

* `description` - (Optional) Specifies the supplementary information about the network ACL.
    This parameter can contain a maximum of 255 characters and cannot contain angle brackets (< or >).

* `inbound_rules` - (Optional)  A list of the IDs of ingress rules associated with the network ACL.

* `outbound_rules` - (Optional) A list of the IDs of egress rules associated with the network ACL.

* `subnets` - (Optional) A list of the IDs of networks associated with the network ACL.

## Attributes Reference

All of the argument attributes are also exported as result attributes:

* `id` - The ID of the network ACL.
* `inbound_policy_id` - The ID of the ingress firewall policy for the network ACL.
* `outbound_policy_id` - The ID of the egress firewall policy for the network ACL.
* `ports` - A list of the port IDs of the subnet gateway.
* `status` - The status of the network ACL.
