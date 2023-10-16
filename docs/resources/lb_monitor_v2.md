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

* `region` - (Optional, String, ForceNew) The region in which to create the resources.
  If omitted, the `region` argument of the provider is used. Changing this creates a new resource.

* `pool_id` - (Required, String, ForceNew) The id of the pool that this monitor will be assigned to.
  Changing this creates a new monitor.

* `type` - (Required, String, ForceNew) The type of probe, which is PING, TCP, HTTP, or HTTPS,
  that is sent by the load balancer to verify the member state. Changing this creates a new monitor.

* `delay` - (Required, Int) The time, in seconds, between sending probes to members.

* `timeout` - (Required, Int) Maximum number of seconds for a monitor to wait for a
  ping reply before it times out. The value must be less than the delay value.

* `max_retries` - (Required, Int) Number of permissible ping failures before
  changing the member's status to INACTIVE. Must be a number between 1 and 10.

* `name` - (Optional, String) The Name of the Monitor.

* `url_path` - (Optional, String) Required for HTTP(S) types. URI path that will be
  accessed if monitor type is HTTP or HTTPS.

* `http_method` - (Optional, String) Required for HTTP(S) types. The HTTP method used
  for requests by the monitor. If this attribute is not specified, it defaults to "GET".

* `expected_codes` - (Optional, String) Required for HTTP(S) types. Expected HTTP codes
  for a passing HTTP(S) monitor. You can either specify a single status like "200", or a range like "200-202".

* `port` - (Optional, Int) Specifies the health check port. The value ranges from 1 to 65536.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID for the monitor.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.
