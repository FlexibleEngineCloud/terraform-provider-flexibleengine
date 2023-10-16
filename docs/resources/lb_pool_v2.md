---
subcategory: "Elastic Load Balance (ELB)"
description: ""
page_title: "flexibleengine_lb_pool_v2"
---

# flexibleengine_lb_pool_v2

Manages an **enhanced** load balancer pool resource within FlexibleEngine.

## Example Usage

```hcl
resource "flexibleengine_lb_pool_v2" "pool_1" {
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = "d9415786-5f1a-428b-b35f-2f1523e146d2"

  persistence {
    type        = "HTTP_COOKIE"
    cookie_name = "testCookie"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the `region` argument of the provider is used. Changing this creates a new pool.

* `protocol` - (Required, String, ForceNew) The protocol - can either be TCP, UDP or HTTP.

  + When the protocol used by the listener is UDP, the protocol of the backend pool must be UDP.
  + When the protocol used by the listener is TCP, the protocol of the backend pool must be TCP.
  + When the protocol used by the listener is HTTP or TERMINATED_HTTPS, the protocol of the backend pool must be HTTP.

  Changing this creates a new pool.

* `lb_method` - (Required, String) The load balancing algorithm to
  distribute traffic to the pool's members. Must be one of
  ROUND_ROBIN, LEAST_CONNECTIONS, or SOURCE_IP.

* `name` - (Optional, String) Human-readable name for the pool.

* `description` - (Optional, String) Human-readable description for the pool.

* `loadbalancer_id` - (Optional, String, ForceNew) The load balancer on which to provision this
  pool. Changing this creates a new pool.
  Note: One of LoadbalancerID or ListenerID must be provided.

* `listener_id` - (Optional, String, ForceNew) The Listener on which the members of the pool
  will be associated with. Changing this creates a new pool.
  Note: One of LoadbalancerID or ListenerID must be provided.

* `persistence` - (Optional, List, ForceNew) Omit this field to prevent session persistence. Indicates
  whether connections in the same session will be processed by the same Pool member or not.
  The [persistence](#lb_persistence) object structure is documented below.
  Changing this creates a new pool.

<a name="lb_persistence"></a>
The `persistence` block supports:

* `type` - (Required, String, ForceNew) The type of persistence mode. The current specification
  supports SOURCE_IP, HTTP_COOKIE, and APP_COOKIE. Changing this will create a new resource.

* `cookie_name` - (Optional, String, ForceNew) The name of the cookie if persistence mode is set
  appropriately. It is Required if `type = APP_COOKIE`. Changing this will create a new resource.

* `timeout` - (Optional, Int, ForceNew) Specifies the sticky session timeout duration in minutes. This parameter is
  invalid when type is set to APP_COOKIE. Changing this will create a new resource.
  The value range varies depending on the protocol of the backend server group:

  + When the protocol of the backend server group is TCP or UDP, the value ranges from 1 to 60.
  + When the protocol of the backend server group is HTTP or HTTPS, the value ranges from 1 to 1440.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID for the pool.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

ELB pool can be imported using the ELB pool ID, e.g.

```shell
terraform import flexibleengine_lb_pool_v2.pool_1 3e3632db-36c6-4b28-a92e-e72e6562daa6
```
