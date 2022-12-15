---
subcategory: "Elastic Load Balance (Dedicated ELB)"
---

# flexibleengine_lb_monitor_v3

Manages an ELB monitor resource within FlexibleEngine.

## Example Usage

```hcl
variable "pool_id" {}

resource "flexibleengine_lb_monitor_v3" "monitor_1" {
  protocol    = "HTTP"
  interval    = 30
  timeout     = 15
  max_retries = 10
  url_path    = "/api"
  port        = 8888
  pool_id     = var.pool_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the ELB monitor resource.
  If omitted, the provider-level region will be used. Changing this creates a new monitor.

* `pool_id` - (Required, String, ForceNew) Specifies the id of the pool that this monitor will be assigned to.

* `protocol` - (Required, String, ForceNew) Specifies the type of probe, which is TCP, HTTP, or HTTPS, that is
  sent by the load balancer to verify the member state. Changing this creates a new monitor.

* `domain_name` - (Optional, String) Specifies the Domain Name of the Monitor.

* `port` - (Optional, Int) Specifies the health check port. The value ranges from 1 to 65535.

* `interval` - (Required, Int) Specifies the time, in seconds, between sending probes to members.

* `timeout` - (Required, Int) Specifies the Maximum number of seconds for a monitor to wait for a ping reply before
  it times out. The value must be less than the delay value.

* `max_retries` - (Required, Int) Specifies the number of permissible ping failures before changing the member's
  status to INACTIVE. Must be a number between 1 and 10.

* `url_path` - (Optional, String) Specifies the required for HTTP(S) types. URI path that will be accessed if monitor
  type is HTTP or HTTPS.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID for the monitor.

## Import

ELB monitor can be imported using the monitor ID, e.g.

```
$ terraform import flexibleengine_lb_monitor_v3.monitor_1 5c20fdad-7288-11eb-b817-0255ac10158b
```
