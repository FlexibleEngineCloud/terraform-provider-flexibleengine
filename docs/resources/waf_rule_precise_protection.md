---
subcategory: "Web Application Firewall (WAF)"
description: ""
page_title: "flexibleengine_waf_rule_precise_protection"
---

# flexibleengine_waf_rule_precise_protection

Manages a WAF Precise Protection Rule resource within FlexibleEngine.

## Example Usage

### A rule takes effect immediately

```hcl
resource "flexibleengine_waf_rule_precise_protection" "rule_1" {
  policy_id = var.policy_id
  name      = "rule_1"
  priority  = 10

  conditions {
    field   = "path"
    logic   = "prefix"
    content = "/login"
  }
}
```

### A rule takes effect at the scheduled time

```hcl
resource "flexibleengine_waf_rule_precise_protection" "rule_1" {
  policy_id = var.policy_id
  name      = "rule_1"
  action     = "block"
  priority   = 20
  start_time = "2021-07-01 00:00:00"
  end_time   = "2021-12-31 23:59:59"

  conditions {
    field   = "ip"
    logic   = "equal"
    content = "192.168.1.1"
  }
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, String, ForceNew) Specifies the WAF policy ID. Changing this creates a new rule.

* `name` - (Required, String) Specifies the name of a precise protection rule. The maximum length is
  256 characters. Only digits, letters, underscores (_), and hyphens (-) are allowed.

* `action` - (Optional, String) Specifies the protective action after the precise protection rule is matched.
  The value can be *block* or *pass*. Defaults to *block*.

* `priority` - (Required, Int) Specifies the priority of a rule being executed. Smaller values correspond to higher priorities.
  If two rules are assigned with the same priority, the rule added earlier has higher priority, the rule added earlier
  has higher priority. The value ranges from 0 to 65535.

* `conditions` - (Required, List) Specifies the condition parameters. The object structure is documented below.

* `start_time` - (Optional, String) Specifies the UTC time when the precise protection rule takes effect.
  The time must be in "yyyy-MM-dd HH:mm:ss" format.
  If not specified, the rule takes effect immediately.

* `end_time` - (Optional, String) Specifies the UTC time when the precise protection rule expires.
  The time must be in "yyyy-MM-dd HH:mm:ss" format.

The `conditions` block supports:

* `field` - (Required, String) Specifies the matched field. The value can be *path*, *user-agent*, *ip*,
  *params*, *cookie*, *referer*, and *header*.

* `subfield` - (Optional, String) Specifies the matched subfield.
  - If `field` is set to *cookie*, subfield indicates cookie name.
  - If `field` is set to *params*, subfield indicates param name.
  - If `field` is set to *header*, subfield indicates an option in the header.

* `logic` - (Required, String) Specifies the logic relationship. The value can be *contain*, *not_contain*,
  *equal*, *not_equal*, *prefix*, *not_prefix*, *suffix*, and *not_suffix*.
  If `field` is set to *ip*, `logic` can only be *equal* or *not_equal*.

* `content` - (Required, String) Specifies the content matching the condition.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The rule ID in UUID format.

## Import

Precise Protection Rules can be imported using the policy ID and rule ID
separated by a slash, e.g.:

```sh
terraform import flexibleengine_waf_rule_precise_protection.rule_1 523083f4543c497faecd25fcfcc0b2a0/620801321b254f8fbc7dafa6bbebe652
```
