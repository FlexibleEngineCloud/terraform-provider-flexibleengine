---
subcategory: "Web Application Firewall (WAF)"
description: ""
page_title: "flexibleengine_waf_policy"
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

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the WAF policy resource.
  If omitted, the provider-level region will be used. Changing this will create a new WAF policy resource.

* `name` - (Required, String) Specifies the policy name. The maximum length is 256 characters.
  Only digits, letters, underscores(_), and hyphens(-) are allowed.

* `protection_mode` - (Optional, String) Specifies the protective action after a rule is matched. Valid values are:
  + *block*: WAF blocks and logs detected attacks.
  + *log*: WAF logs detected attacks only.

* `level` - (Optional, Int) Specifies the protection level. Valid values are:
  + *1*: low
  + *2*: medium
  + *3*: high

* `full_detection` - (Optional, Bool) Specifies the detection mode in Precise Protection. Valid values are:
  + *true*: full detection, Full detection finishes all threat detections before blocking requests that
    meet Precise Protection specified conditions.
  + *false*: instant detection. Instant detection immediately ends threat detection after blocking a request that
    meets Precise Protection specified conditions.

* `domains` - (Optional, List) An array of domain IDs.

* `protection_status` - (Optional, List) Specifies the protection switches. The [protection_status](#waf_protection_status)
  object structure is documented below.

<a name="waf_protection_status"></a>
The `protection_status` block supports:

* `basic_web_protection` - (Optional, Bool) Specifies whether Basic Web Protection is enabled.

* `general_check` - (Optional, Bool) Specifies whether General Check in Basic Web Protection is enabled.

* `crawler_engine` - (Optional, Bool) Specifies whether the Search Engine switch in Basic Web Protection is enabled.

* `crawler_scanner` - (Optional, Bool) Specifies whether the Scanner switch in Basic Web Protection is enabled.

* `crawler_script` - (Optional, Bool) Specifies whether the Script Tool switch in Basic Web Protection is enabled.

* `crawler_other` - (Optional, Bool) Specifies whether detection of other crawlers in Basic Web Protection is enabled.

* `webshell` - (Optional, Bool) Specifies whether webshell detection in Basic Web Protection is enabled.

* `cc_protection` - (Optional, Bool) Specifies whether CC Attack Protection is enabled.

* `precise_protection` - (Optional, Bool) Specifies whether Precise Protection is enabled.

* `blacklist` - (Optional, Bool) Specifies whether Blacklist and Whitelist is enabled.

* `data_masking` - (Optional, Bool) Specifies whether Data Masking is enabled.

* `false_alarm_masking` - (Optional, Bool) Specifies whether False Alarm Masking is enabled.

* `web_tamper_protection` - (Optional, Bool) Specifies whether Web Tamper Protection is enabled.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The policy ID in UUID format.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

Policies can be imported using the `id`, e.g.

```sh
terraform import flexibleengine_waf_policy.policy_1 c5946141e52441d9b13c5e9d4e9560c7
```
