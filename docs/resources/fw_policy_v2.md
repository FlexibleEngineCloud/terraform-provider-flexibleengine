---
subcategory: "Network ACL"
description: ""
page_title: "flexibleengine_fw_policy_v2"
---

# flexibleengine_fw_policy_v2

Manages a v2 firewall policy resource within FlexibleEngine.

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

  rules = [
    flexibleengine_fw_rule_v2.rule_1.id,
    flexibleengine_fw_rule_v2.rule_2.id,
  ]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `name` - (Optional, String) A name for the firewall policy. Changing this
  updates the `name` of an existing firewall policy.

* `description` - (Optional, String) A description for the firewall policy. Changing
  this updates the `description` of an existing firewall policy.

* `rules` - (Optional, List) An array of one or more firewall rules that comprise
  the policy. Changing this results in adding/removing rules from the
  existing firewall policy.

* `audited` - (Optional, Bool) Audit status of the firewall policy
  (must be "true" or "false" if provided - defaults to "false").
  This status is set to "false" whenever the firewall policy or any of its
  rules are changed. Changing this updates the `audited` status of an existing
  firewall policy.

* `shared` - (Optional, Bool) Sharing status of the firewall policy (must be "true"
  or "false" if provided). If this is "true" the policy is visible to, and
  can be used in, firewalls in other tenants. Changing this updates the
  `shared` status of an existing firewall policy. Only administrative users
  can specify if the policy should be shared.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the policy.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.

## Import

Firewall Policies can be imported using the `id`, e.g.

```shell
terraform import flexibleengine_fw_policy_v2.policy_1 07f422e6-c596-474b-8b94-fe2c12506ce0
```
