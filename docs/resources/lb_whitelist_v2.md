---
subcategory: "Deprecated"
description: ""
page_title: "flexibleengine_lb_whitelist_v2"
---

# flexibleengine_lb_whitelist_v2

Manages an **enhanced** load balancer whitelist resource within FlexibleEngine.

## Example Usage

```hcl
resource "flexibleengine_lb_whitelist_v2" "whitelist_1" {
  enable_whitelist = true
  whitelist        = "192.168.11.1,192.168.0.1/24,192.168.201.18/8"
  listener_id      = "d9415786-5f1a-428b-b35f-2f1523e146d2"
}
```

## Argument Reference

The following arguments are supported:

* `listener_id` - (Required, String, ForceNew) The Listener ID that the whitelist will be associated with.
  Changing this creates a new whitelist.

* `enable_whitelist` - (Optional, Bool) Specify whether to enable access control.

* `whitelist` - (Optional, String) Specifies the IP addresses in the whitelist. Use commas(,) to separate
  the multiple IP addresses.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID for the whitelist.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.
