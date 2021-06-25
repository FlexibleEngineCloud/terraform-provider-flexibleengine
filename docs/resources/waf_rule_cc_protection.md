---
subcategory: "Web Application Firewall (WAF)"
---

# flexibleengine_waf_rule_cc_protection

Manages a WAF CC Attack Protection Rule resource within FlexibleEngine.

## Example Usage

```hcl
resource "flexibleengine_waf_policy" "policy_1" {
  name = "policy_1"
}

resource "flexibleengine_waf_rule_cc_protection" "rule_1" {
  policy_id    = flexibleengine_waf_policy.policy_1.id
  path         = "/abc"
  limit_num    = 10
  limit_period = 60  
  mode         = "cookie"
  cookie       = "sessionid"

  action             = "block"
  block_time         = 10
  block_page_type    = "application/json"
  block_page_content = "{\"error\":\"forbidden\"}"
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, String, ForceNew) Specifies the WAF policy ID. Changing this creates a new rule.

* `path` - (Required, String) Specifies the URL to which the rule applies. The path ending with * indicates
  that the path is used as a prefix. For example, if the path to be protected is /admin/test.php or /adminabc,
  set Path to /admin*.

* `limit_num` - (Required, Int) Specifies the number of requests allowed from a web visitor in a rate limiting period.
  The value ranges from 0 to 2^32.

* `limit_period` - (Required, Int) Specifies the rate limiting period. The value ranges from 0 seconds to 2^32 seconds.

* `mode` - (Required, String) Specifies the rate limit mode. Valid Options are:
  * *ip* - A web visitor is identified by the IP address.
  * *cookie* - A web visitor is identified by the cookie key value.
  * *other* - A web visitor is identified by the Referer field(user-defined request source).

* `cookie` - (Optional, String) Specifies the cookie name. This field is mandatory when `mode` is set to *cookie*.

* `content` - (Optional, String) Specifies the category content. The format is as follows: http://www.example.com/path.
  This field is mandatory when `mode` is set to *other*.

* `action` - (Required, String) Specifies the action when the number of requests reaches the upper limit. Valid Options are:
  * *block* - block the requests.
  * *captcha* - Verification code. The user needs to enter the correct verification code after blocking to restore the correct access page.

  If `mode` is set to *other*, this parameter value can only be *block*.

* `block_time` - (Optional, Int) Specifies the lock duration. The value ranges from 0 seconds to 2^32 seconds.

* `block_page_type` - (Optional, String) Specifies the type of the returned page.
  The options are `application/json`, `text/html`, and `text/xml`.

* `block_page_content` - (Optional, String) Specifies the content of the returned page.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` -  The rule ID in UUID format.

## Import

CC Attack Protection Rules can be imported using the policy ID and rule ID
separated by a slash, e.g.:

```sh
terraform import flexibleengine_waf_rule_cc_protection.rule_1 523083f4543c497faecd25fcfcc0b2a0/dd3c14e91550453f81cff5fc3b7c3e89
```
