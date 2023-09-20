---
subcategory: "Network ACL"
description: ""
page_title: "flexibleengine_fw_rule_v2"
---

# flexibleengine_fw_rule_v2

Manages a v2 firewall rule resource within FlexibleEngine.

## Example Usage

```hcl
resource "flexibleengine_fw_rule_v2" "rule_1" {
  name             = "my_rule"
  description      = "drop TELNET traffic"
  action           = "deny"
  protocol         = "tcp"
  destination_port = "23"
  enabled          = "true"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to obtain the v2 networking client.
  A Compute client is needed to create a firewall rule. If omitted, the
  `region` argument of the provider is used. Changing this creates a new
  firewall rule.

* `name` - (Optional, String) A unique name for the firewall rule. Changing this
  updates the `name` of an existing firewall rule.

* `description` - (Optional, String) A description for the firewall rule. Changing this
  updates the `description` of an existing firewall rule.

* `protocol` - (Required, String) The protocol type on which the firewall rule operates.
  Valid values are: `tcp`, `udp`, `icmp`, and `any`. Changing this updates the
  `protocol` of an existing firewall rule.

* `action` - (Required, String) Action to be taken ( must be "allow" or "deny") when the
  firewall rule matches. Changing this updates the `action` of an existing
  firewall rule.

* `ip_version` - (Optional, Int) IP version, either 4 (default) or 6. Changing this
  updates the `ip_version` of an existing firewall rule.

* `source_ip_address` - (Optional, String) The source IP address on which the firewall
  rule operates. Changing this updates the `source_ip_address` of an existing
  firewall rule.

* `destination_ip_address` - (Optional, String) The destination IP address on which the
  firewall rule operates. Changing this updates the `destination_ip_address`
  of an existing firewall rule.

* `source_port` - (Optional, String) The source port on which the firewall
  rule operates. Changing this updates the `source_port` of an existing
  firewall rule.

* `destination_port` - (Optional, String) The destination port on which the firewall
  rule operates. Changing this updates the `destination_port` of an existing
  firewall rule.

* `enabled` - (Optional, Bool) Enabled status for the firewall rule (must be "true"
  or "false" if provided - defaults to "true"). Changing this updates the
  `enabled` status of an existing firewall rule.

## Attribute Reference

All the arguments above can also be exported attributes.

## Import

Firewall Rules can be imported using the `id`, e.g.

```shell
terraform import flexibleengine_fw_rule_v2.rule_1 8dbc0c28-e49c-463f-b712-5c5d1bbac327
```
