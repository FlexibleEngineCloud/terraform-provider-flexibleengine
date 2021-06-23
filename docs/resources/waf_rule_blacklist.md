---
subcategory: "Web Application Firewall (WAF)"
---

# flexibleengine_waf_rule_blacklist

Manages a WAF blacklist and whitelist rule resource within FlexibleEngine.

## Example Usage

```hcl
resource "flexibleengine_waf_policy" "policy_1" {
  name = "policy_1"
}

resource "flexibleengine_waf_rule_blacklist" "rule_1" {
  policy_id = flexibleengine_waf_policy.policy_1.id
  address   = "192.168.0.0/24"
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, String, ForceNew) Specifies the WAF policy ID. Changing this creates a new rule.

* `address` - (Required, String) Specifies the IP address or range. For example, 192.168.0.125 or 192.168.0.0/24.

* `action` - (Optional, Int) Specifies the protective action. 1: Whitelist, 0: Blacklist.
  If you do not configure the parameter, the value is Blacklist by default.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The rule ID in UUID format.

## Import

Blacklist Rules can be imported using the policy ID and rule ID
separated by a slash, e.g.:

```sh
terraform import flexibleengine_waf_rule_blacklist.rule_1 523083f4543c497faecd25fcfcc0b2a0/e7f49f736bc74b828ce45e0e5c49d156
```
