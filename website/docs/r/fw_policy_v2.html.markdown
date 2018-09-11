---
layout: "flexibleengine"
page_title: "FlexibleEngine: flexibleengine_fw_policy_v2"
sidebar_current: "docs-flexibleengine-resource-fw-policy-v1"
description: |-
  Manages a v1 firewall policy resource within FlexibleEngine.
---

# flexibleengine\_fw\_policy_v2

Manages a v1 firewall policy resource within FlexibleEngine.

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

  rules = ["${flexibleengine_fw_rule_v2.rule_1.id}",
    "${flexibleengine_fw_rule_v2.rule_2.id}",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the v1 networking client.
    A networking client is needed to create a firewall policy. If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    firewall policy.

* `name` - (Optional) A name for the firewall policy. Changing this
    updates the `name` of an existing firewall policy.

* `description` - (Optional) A description for the firewall policy. Changing
    this updates the `description` of an existing firewall policy.

* `rules` - (Optional) An array of one or more firewall rules that comprise
    the policy. Changing this results in adding/removing rules from the
    existing firewall policy.

* `audited` - (Optional) Audit status of the firewall policy
    (must be "true" or "false" if provided - defaults to "false").
    This status is set to "false" whenever the firewall policy or any of its
    rules are changed. Changing this updates the `audited` status of an existing
    firewall policy.

* `shared` - (Optional) Sharing status of the firewall policy (must be "true"
    or "false" if provided). If this is "true" the policy is visible to, and
    can be used in, firewalls in other tenants. Changing this updates the
    `shared` status of an existing firewall policy. Only administrative users
    can specify if the policy should be shared.

* `value_specs` - (Optional) Map of additional options.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `audited` - See Argument Reference above.
* `shared` - See Argument Reference above.

## Import

Firewall Policies can be imported using the `id`, e.g.

```
$ terraform import flexibleengine_fw_policy_v2.policy_1 07f422e6-c596-474b-8b94-fe2c12506ce0
```
