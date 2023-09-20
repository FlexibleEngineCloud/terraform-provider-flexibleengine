---
subcategory: "Elastic Load Balance (ELB)"
description: ""
page_title: "flexibleengine_lb_monitor_v2"
---

# flexibleengine_lb_monitor_v2

Manages an **enhanced** load balancer monitor resource within FlexibleEngine.

## Example Usage

```hcl
resource "flexibleengine_lb_monitor_v2" "monitor_1" {
  pool_id     = flexibleengine_lb_pool_v2.pool_1.id
  type        = "PING"
  delay       = 20
  timeout     = 10
  max_retries = 5
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to obtain the V2 Networking client.
    A Networking client is needed to be created. If omitted, the
    `region` argument of the provider is used. Changing this creates a new monitor.

* `pool_id` - (Required, String, ForceNew) The id of the pool that this monitor will be assigned to.
    Changing this creates a new monitor.

* `type` - (Required, String, ForceNew) The type of probe, which is PING, TCP, HTTP, or HTTPS,
    that is sent by the load balancer to verify the member state. Changing this creates a new monitor.

* `delay` - (Required, Int) The time, in seconds, between sending probes to members.

* `timeout` - (Required, Int) Maximum number of seconds for a monitor to wait for a
    ping reply before it times out. The value must be less than the delay
    value.

* `max_retries` - (Required, Int) Number of permissible ping failures before
    changing the member's status to INACTIVE. Must be a number between 1
    and 10.

* `name` - (Optional, String) The Name of the Monitor.

* `url_path` - (Optional, String) Required for HTTP(S) types. URI path that will be
    accessed if monitor type is HTTP or HTTPS.

* `http_method` - (Optional, String) Required for HTTP(S) types. The HTTP method used
    for requests by the monitor. If this attribute is not specified, it
    defaults to "GET".

* `expected_codes` - (Optional, String) Required for HTTP(S) types. Expected HTTP codes
    for a passing HTTP(S) monitor. You can either specify a single status like
    "200", or a range like "200-202".

* `port` - (Optional, Int) Specifies the health check port. The value ranges from 1 to 65536.

* `admin_state_up` - (Optional, Bool) The administrative state of the monitor.
    A valid value is true (UP) or false (DOWN).

* `tenant_id` - (Optional, String, ForceNew) The UUID of the tenant who owns the monitor.
    Only administrative users can specify a tenant UUID other than their own.
    Changing this creates a new monitor.

## Attribute Reference

The following attributes are exported:

* `id` - The unique ID for the monitor.
* `tenant_id` - See Argument Reference above.
* `type` - See Argument Reference above.
* `delay` - See Argument Reference above.
* `timeout` - See Argument Reference above.
* `max_retries` - See Argument Reference above.
* `url_path` - See Argument Reference above.
* `http_method` - See Argument Reference above.
* `expected_codes` - See Argument Reference above.
* `admin_state_up` - See Argument Reference above.
* `port` - See Argument Reference above.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.
