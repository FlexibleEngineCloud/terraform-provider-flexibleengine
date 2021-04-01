---
subcategory: "Elastic Load Balance (ELB)"
---

# flexibleengine\_elb\_backend

Manages an elastic loadbalancer backend resource within FlexibleEngine.

## Example Usage

```hcl
resource "flexibleengine_elb_loadbalancer" "elb" {
  name = "elb"
  type = "External"
  description = "test elb"
  vpc_id = "e346dc4a-d9a6-46f4-90df-10153626076e"
  admin_state_up = 1
  bandwidth = 5
}

resource "flexibleengine_elb_listener" "listener" {
  name = "test-elb-listener"
  description = "great listener"
  protocol = "TCP"
  backend_protocol = "TCP"
  protocol_port = 12345
  backend_port = 8080
  lb_algorithm = "roundrobin"
  loadbalancer_id = "${flexibleengine_elb_loadbalancer.elb.id}"
  timeouts {
	create = "5m"
	update = "5m"
	delete = "5m"
  }
}

resource "flexibleengine_elb_backend" "backend" {
  address = "192.168.0.211"
  listener_id = "${flexibleengine_elb_listener.listener.id}"
  server_id = "8f7a32f1-f66c-4d13-9b17-3a13f9f0bb8d"
}
```

## Argument Reference

The following arguments are supported:

* `listener_id` - (Required) Specifies the listener ID.

* `server_id` - (Required) Specifies the backend member ID.

* `address` - (Required) Specifies the private IP address of the backend member.

## Attributes Reference

The following attributes are exported:

* `listener_id` - See Argument Reference above.
* `server_id` - See Argument Reference above.
* `address` - See Argument Reference above.
* `server_address` - Specifies the floating IP address assigned to the backend member.
* `id` - Specifies the backend member ID.
* `status` - Specifies the backend ECS status. The value is ACTIVE, PENDING,
    or ERROR.
* `health_status` - Specifies the health check status. The value is NORMAL,
    ABNORMAL, or UNAVAILABLE.
* `update_time` - Specifies the time when information about the backend member
    was updated.
* `create_time` - Specifies the time when the backend member was created.
* `server_name` - Specifies the backend member name.
* `listeners` - Specifies the listener to which the backend member belongs.
