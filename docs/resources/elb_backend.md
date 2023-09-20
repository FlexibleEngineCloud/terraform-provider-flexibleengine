---
subcategory: "Deprecated"
description: ""
page_title: "flexibleengine_elb_backend"
---

# flexibleengine_elb_backend

!> **Warning:** Classic load balancers are no longer provided, using elastic load balancers instead.

Manages a **classic** lb backend resource within FlexibleEngine.

## Example Usage

```hcl
resource "flexibleengine_elb_backend" "backend" {
  listener_id = flexibleengine_elb_listener.listener.id
  server_id   = "8f7a32f1-f66c-4d13-9b17-3a13f9f0bb8d"
  address     = "192.168.0.211"
}
```

## Argument Reference

The following arguments are supported:

* `listener_id` - (Required, String, ForceNew) Specifies the listener ID.

* `server_id` - (Required, String, ForceNew) Specifies the backend member ID.

* `address` - (Required, String, ForceNew) Specifies the private IP address of the backend member.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.
