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

* `region` - (Optional, String, ForceNew) The region in which to obtain the V2 Networking client.
    A Networking client is needed to create an . If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    pool.

* `name` - (Optional, String) Human-readable name for the pool.

* `description` - (Optional, String) Human-readable description for the pool.

* `protocol` = (Required) The protocol - can either be TCP, UDP or HTTP.

    + When the protocol used by the listener is UDP, the protocol of the backend pool must be UDP.
    + When the protocol used by the listener is TCP, the protocol of the backend pool must be TCP.
    + When the protocol used by the listener is HTTP or TERMINATED_HTTPS, the protocol of the backend pool must be HTTP.

    Changing this creates a new pool.

* `loadbalancer_id` - (Optional, String, ForceNew) The load balancer on which to provision this
    pool. Changing this creates a new pool.
    Note: One of LoadbalancerID or ListenerID must be provided.

* `listener_id` - (Optional, String, ForceNew) The Listener on which the members of the pool
    will be associated with. Changing this creates a new pool.
    Note: One of LoadbalancerID or ListenerID must be provided.

* `lb_method` - (Required, String) The load balancing algorithm to
    distribute traffic to the pool's members. Must be one of
    ROUND_ROBIN, LEAST_CONNECTIONS, or SOURCE_IP.

* `persistence` - (Optional, List, ForceNew) Omit this field to prevent session persistence.  Indicates
    whether connections in the same session will be processed by the same Pool
    member or not. Changing this creates a new pool.

* `admin_state_up` - (Optional, Bool) The administrative state of the pool.
    A valid value is true (UP) or false (DOWN).

The `persistence` argument supports:

* `type` - (Required, String, ForceNew) The type of persistence mode. The current specification
    supports SOURCE_IP, HTTP_COOKIE, and APP_COOKIE.

* `cookie_name` - (Optional, String, ForceNew) The name of the cookie if persistence mode is set
    appropriately. Required if `type = APP_COOKIE`.

* `timeout` - (Optional, Int, ForceNew) Specifies the sticky session timeout duration in minutes. This parameter is
  invalid when type is set to APP_COOKIE. The value range varies depending on the protocol of the backend server group:

  + When the protocol of the backend server group is TCP or UDP, the value ranges from 1 to 60.
  + When the protocol of the backend server group is HTTP or HTTPS, the value ranges from 1 to 1440.

## Attributes Reference

The following attributes are exported:

* `id` - The unique ID for the pool.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `protocol` - See Argument Reference above.
* `lb_method` - See Argument Reference above.
* `persistence` - See Argument Reference above.
* `admin_state_up` - See Argument Reference above.
