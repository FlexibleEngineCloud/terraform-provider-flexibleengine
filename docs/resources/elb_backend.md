---
subcategory: "Deprecated"
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

* `listener_id` - (Required) Specifies the listener ID.

* `server_id` - (Required) Specifies the backend member ID.

* `address` - (Required) Specifies the private IP address of the backend member.

## Attributes Reference

The following attributes are exported:

* `id` - Specifies the backend member ID.
* `listener_id` - See Argument Reference above.
* `server_id` - See Argument Reference above.
* `address` - See Argument Reference above.
* `server_address` - Specifies the floating IP address assigned to the backend member.
* `status` - Specifies the backend ECS status. The value is ACTIVE, PENDING,
    or ERROR.
* `health_status` - Specifies the health check status. The value is NORMAL,
    ABNORMAL, or UNAVAILABLE.
* `update_time` - Specifies the time when information about the backend member
    was updated.
* `create_time` - Specifies the time when the backend member was created.
* `server_name` - Specifies the backend member name.
* `listeners` - Specifies the listener to which the backend member belongs.
