---
subcategory: "Web Application Firewall (WAF)"
---

# flexibleengine_waf_rule_web_tamper_protection

Manages a WAF Web Tamper Protection Rule resource within FlexibleEngine.

## Example Usage

```hcl
resource "flexibleengine_waf_policy" "policy_1" {
  name = "policy_1"
}

resource "flexibleengine_waf_rule_web_tamper_protection" "rule_1" {
  policy_id = flexibleengine_waf_policy.policy_1.id
  domain    = "www.abc.com"
  path      = "/a"
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, String, ForceNew) Specifies the WAF policy ID. Changing this creates a new rule.

* `domain` - (Required, String, ForceNew) Specifies the domain name. Changing this creates a new rule.

* `path` - (Required, String, ForceNew) Specifies the URL protected by the web tamper protection rule,
  excluding a domain name. Changing this creates a new rule.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The rule ID in UUID format.

## Import

Web Tamper Protection Rules can be imported using the policy ID and rule ID
separated by a slash, e.g.:

```sh
terraform import flexibleengine_waf_rule_web_tamper_protection.rule_1 523083f4543c497faecd25fcfcc0b2a0/5b3b07fedc3642d18e424b2e45aebc8a
```
