---
subcategory: "Web Application Firewall (WAF)"
---

# flexibleengine_waf_certificate

Get the certificate in the WAF, including the one pushed from SCM.

## Example Usage

```hcl
data "flexibleengine_waf_certificate" "certificate_1" {
  name = "certificate name"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region in which to query the WAF. If omitted, the provider-level region will be
  used.

* `name` - (Required, String) The name of certificate. The value is case-sensitive and supports fuzzy matching.
  The certificate name is not unique. Only returns the last created one when matched multiple certificates.

* `expire_status` - (Optional, Int) The expired status of certificate. Defaults is **0**. The value can be:
  + **0**: not expire.
  + **1**: has expired.
  + **2**: will be expired soon.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The certificate ID in UUID format.

* `expiration` - Indicates the time when the certificate expires.
