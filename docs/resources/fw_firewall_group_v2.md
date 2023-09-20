---
subcategory: "Network ACL"
description: ""
page_title: "flexibleengine_fw_firewall_group_v2"
---

# flexibleengine_fw_firewall_group_v2

Manages a v2 firewall group resource within FlexibleEngine.

## Example Usage

```hcl
resource "flexibleengine_fw_rule_v2" "rule_1" {
  name             = "my-rule-1"
  description      = "drop TELNET traffic"
  action           = "deny"
  protocol         = "tcp"
  destination_port = "23"
  enabled          = "true"
}

resource "flexibleengine_fw_rule_v2" "rule_2" {
  name             = "my-rule-2"
  description      = "drop NTP traffic"
  action           = "deny"
  protocol         = "udp"
  destination_port = "123"
  enabled          = "false"
}

resource "flexibleengine_fw_policy_v2" "policy_1" {
  name = "my-policy"

  rules = [flexibleengine_fw_rule_v2.rule_1.id,
    flexibleengine_fw_rule_v2.rule_2.id,
  ]
}

resource "flexibleengine_fw_firewall_group_v2" "firewall_group_1" {
  name      = "my-firewall-group"
  ingress_policy_id = flexibleengine_fw_policy_v2.policy_1.id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to obtain the v2 networking client.
  A networking client is needed to create a firewall group. If omitted, the
  `region` argument of the provider is used. Changing this creates a new
  firewall group.

* `ingress_policy_id` - (Optional, String) The ingress policy resource id for the firewall group. Changing
  this updates the `ingress_policy_id` of an existing firewall group.

* `egress_policy_id` - (Optional, String) The egress policy resource id for the firewall group. Changing
  this updates the `egress_policy_id` of an existing firewall group.

* `name` - (Optional, String) A name for the firewall group. Changing this
  updates the `name` of an existing firewall group.

* `description` - (Optional, String) A description for the firewall group. Changing this
  updates the `description` of an existing firewall group.

* `ports` - (Optional, List) Port(s) to associate this firewall group instance
  with. Must be a list of strings. Changing this updates the associated routers
  of an existing firewall group.

## Attribute Reference

All the arguments above can also be exported attributes.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

Firewall Groups can be imported using the `id`, e.g.

```shell
terraform import flexibleengine_fw_firewall_group_v2.firewall_group_1 c9e39fb2-ce20-46c8-a964-25f3898c7a97
```
