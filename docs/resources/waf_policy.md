---
subcategory: "Web Application Firewall (WAF)"
---

# flexibleengine_waf_policy

Manages a WAF policy resource within FlexibleEngine.

## Example Usage

```hcl
resource "flexibleengine_waf_policy" "policy_1" {
  name = "policy_1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies the policy name. The maximum length is 256 characters.
  Only digits, letters, underscores(_), and hyphens(-) are allowed.

* `protection_mode` - (Optional, String) Specifies the protective action after a rule is matched. Valid values are:
  - *block*: WAF blocks and logs detected attacks.
  - *log*: WAF logs detected attacks only.

* `level` - (Optional, Int) Specifies the protection level. Valid values are:
  - *1*: low
  - *2*: medium
  - *3*: high

* `full_detection` - (Optional, Bool) Specifies the detection mode in Precise Protection. Valid values are:
  * *true*: full detection, Full detection finishes all threat detections before blocking requests that
    meet Precise Protection specified conditions.
  * *false*: instant detection. Instant detection immediately ends threat detection after blocking a request that
    meets Precise Protection specified conditions.

* `domains` - (Optional, List) An array of domain IDs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The policy ID in UUID format.

* `protection_status` - The protection switches. The object structure is documented below.

The `protection_status` block supports:

* `basic_web_protection` - Indicates whether Basic Web Protection is enabled.

* `general_check` - Indicates whether General Check in Basic Web Protection is enabled.

* `crawler_engine` - Indicates whether the Search Engine switch in Basic Web Protection is enabled.

* `crawler_scanner` - Indicates whether the Scanner switch in Basic Web Protection is enabled.

* `crawler_script` - Indicates whether the Script Tool switch in Basic Web Protection is enabled.

* `crawler_other` - Indicates whether detection of other crawlers in Basic Web Protection is enabled.

* `webshell` - Indicates whether webshell detection in Basic Web Protection is enabled.

* `cc_protection` - Indicates whether CC Attack Protection is enabled.

* `precise_protection` - Indicates whether Precise Protection is enabled.

* `blacklist` - Indicates whether Blacklist and Whitelist is enabled.

* `data_masking` - Indicates whether Data Masking is enabled.

* `false_alarm_masking` - Indicates whether False Alarm Masking is enabled.

* `web_tamper_protection` - Indicates whether Web Tamper Protection is enabled.

## Import

Policies can be imported using the `id`, e.g.

```sh
terraform import flexibleengine_waf_policy.policy_1 c5946141e52441d9b13c5e9d4e9560c7
```
