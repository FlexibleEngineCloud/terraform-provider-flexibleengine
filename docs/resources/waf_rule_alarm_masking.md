---
subcategory: "Web Application Firewall (WAF)"
description: ""
page_title: "flexibleengine_waf_rule_alarm_masking"
---

# flexibleengine_waf_rule_alarm_masking

Manages a WAF False Alarm Masking Rule resource within FlexibleEngine.

## Example Usage

```hcl
resource "flexibleengine_waf_policy" "policy_1" {
  name = "policy_1"
}

resource "flexibleengine_waf_rule_alarm_masking" "rule_1" {
  policy_id = flexibleengine_waf_policy.policy_1.id
  path      = "/a"
  event_id  = "3737fb122f2140f39292f597ad3b7e9a"
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, String, ForceNew) Specifies the WAF policy ID. Changing this creates a new rule.

* `path` - (Required, String) Specifies a misreported URL excluding a domain name.

* `event_id` - (Required, String) Specifies the event ID. It is the ID of a misreported event
  in Events whose type is not *Custom*.

## Attributes Reference

The following attributes are exported:

* `id` - The rule ID in UUID format.
* `event_type` - The event type.

## Import

Alarm Masking Rules can be imported using the policy ID and rule ID
separated by a slash, e.g.:

```sh
terraform import flexibleengine_waf_rule_alarm_masking.rule_1 44d887434169475794b2717438f7fa78/6cdc116040d444f6b3e1bf1dd629f1d0
```
