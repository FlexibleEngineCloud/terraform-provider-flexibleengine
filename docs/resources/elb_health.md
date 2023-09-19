---
subcategory: "Deprecated"
description: ""
page_title: "flexibleengine_elb_health"
---

# flexibleengine_elb_health

!> **Warning:** Classic load balancers are no longer provided, using elastic load balancers instead.

Manages a **classic** lb health check resource within FlexibleEngine.

## Example Usage

```hcl
resource "flexibleengine_elb_health" "healthcheck" {
  listener_id              = flexibleengine_elb_listener.listener.id
  healthcheck_protocol     = "TCP"
  healthcheck_connect_port = 22
  healthy_threshold        = 5
  healthcheck_timeout      = 25
  healthcheck_interval     = 3
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the elb health. If
    omitted, the `region` argument of the provider is used. Changing this
    creates a new elb health.

* `listener_id` - (Required, String, ForceNew) Specifies the ID of the listener to which the health
    check belongs.

* `healthcheck_protocol` - (Optional, String) Specifies the protocol used for the health
    check. The value can be HTTP or TCP (case-insensitive).

* `healthcheck_uri` - (Optional, String) Specifies the URI for health check. This parameter
    is valid when healthcheck_ protocol is HTTP. The value is a string of 1 to 80
    characters that must start with a slash (/) and can only contain letters, digits,
    and special characters, such as -/.%?#&.

* `healthcheck_connect_port` - (Optional, Int) Specifies the port used for the health
    check. The value ranges from 1 to 65535.

* `healthy_threshold` - (Optional, Int) Specifies the threshold at which the health
    check result is success, that is, the number of consecutive successful health
    checks when the health check result of the backend server changes from fail
    to success. The value ranges from 1 to 10.

* `unhealthy_threshold` - (Optional, Int) Specifies the threshold at which the health
    check result is fail, that is, the number of consecutive failed health checks
    when the health check result of the backend server changes from success to fail.
    The value ranges from 1 to 10.

* `healthcheck_timeout` - (Optional, Int) Specifies the maximum timeout duration
    (s) for the health check. The value ranges from 1 to 50.

* `healthcheck_interval` - (Optional, Int) Specifies the maximum interval (s) for
    health check. The value ranges from 1 to 5.

## Attribute Reference

The following attributes are exported:

* `id` - Specifies the health check ID.
* `region` - See Argument Reference above.
* `listener_id` - See Argument Reference above.
* `healthcheck_protocol` - See Argument Reference above.
* `healthcheck_uri` - See Argument Reference above.
* `healthcheck_connect_port` - See Argument Reference above.
* `healthy_threshold` - See Argument Reference above.
* `unhealthy_threshold` - See Argument Reference above.
* `healthcheck_timeout` - See Argument Reference above.
* `healthcheck_interval` - See Argument Reference above.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.
