---
subcategory: "Web Application Firewall (WAF)"
---

# flexibleengine_waf_rule_data_masking

Manages a WAF Data Masking Rule resource within FlexibleEngine.

## Example Usage

```hcl
resource "flexibleengine_waf_policy" "policy_1" {
  name = "policy_1"
}

resource "flexibleengine_waf_rule_data_masking" "rule_1" {
  policy_id = flexibleengine_waf_policy.policy_1.id
  path      = "/login"
  field     = "params"
  subfield  = "password"
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, String, ForceNew) Specifies the WAF policy ID. Changing this creates a new rule.

* `path` - (Required, String) Specifies the URL to which the data masking rule applies (exact match by default).

* `field` - (Required, String) Specifies the masked field. The options are *params* and *header*.

* `subfield` - (Required, String) Specifies the masked subfield.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The rule ID in UUID format.

## Import

Data Masking Rules can be imported using the policy ID and rule ID
separated by a slash, e.g.:

```sh
terraform import flexibleengine_waf_rule_data_masking.rule_1 523083f4543c497faecd25fcfcc0b2a0/c6482bd0059148559b625f78e8ce92be
```
