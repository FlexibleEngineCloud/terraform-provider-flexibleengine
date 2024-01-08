---
subcategory: "Elastic Load Balance (ELB)"
description: ""
page_title: "flexibleengine_lb_whitelist"
---

# flexibleengine_lb_whitelist

Manages an **enhanced** load balancer whitelist resource within FlexibleEngine.

## Example Usage

```hcl
resource "flexibleengine_lb_listener" "listener_1" {
  name            = "listener_1"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = var.loadbalancer_id
}

resource "flexibleengine_lb_whitelist" "whitelist_1" {
  enable_whitelist = true
  whitelist        = "192.168.11.1,192.168.0.1/24,192.168.201.18/8"
  listener_id      = flexibleengine_lb_listener.listener_1.id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the ELB whitelist resource. If omitted, the
  provider-level region will be used. Changing this creates a new whitelist.

* `listener_id` - (Required, String, ForceNew) The Listener ID that the whitelist will be associated with. Changing this
  creates a new whitelist.

* `enable_whitelist` - (Optional, Bool) Specify whether to enable access control.

* `whitelist` - (Optional, String) Specifies the IP addresses in the whitelist. Use commas(,) to separate the multiple
  IP addresses.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID for the whitelist.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

ELB whitelist can be imported using the whitelist ID, e.g.

```shell
terraform import flexibleengine_lb_whitelist.whitelist_1 5c20fdad-7288-11eb-b817-0255ac10158b
```
