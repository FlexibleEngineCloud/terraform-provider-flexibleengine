---
layout: "flexibleengine"
page_title: "FlexibleEngine: flexibleengine_lb_loadbalancer_v2"
sidebar_current: "docs-flexibleengine-resource-lb-loadbalancer-v2"
description: |-
  Manages a V2 loadbalancer resource within FlexibleEngine.
---

# flexibleengine\_lb\_loadbalancer\_v2

Manages a V2 loadbalancer resource within FlexibleEngine.

## Example Usage

```hcl
resource "flexibleengine_lb_loadbalancer_v2" "lb_1" {
  vip_subnet_id = "d9415786-5f1a-428b-b35f-2f1523e146d2"

  tags = {
    key = "value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to create the loadbalancer resource.
    If omitted, the `region` argument of the provider is used.
    Changing this creates a new loadbalancer.

* `vip_subnet_id` - (Required) The network on which to allocate the
    loadbalancer's address. A tenant can only create Loadbalancers on networks
    authorized by policy (e.g. networks that belong to them or networks that
    are shared).  Changing this creates a new loadbalancer.

* `name` - (Optional) Human-readable name for the loadbalancer. Does not have
    to be unique.

* `description` - (Optional) Human-readable description for the loadbalancer.

* `tenant_id` - (Optional) Required for admins. The UUID of the tenant who owns
    the loadbalancer.  Only administrative users can specify a tenant UUID
    other than their own.  Changing this creates a new loadbalancer.

* `vip_address` - (Optional) The ip address of the load balancer.
    Changing this creates a new loadbalancer.

* `admin_state_up` - (Optional) The administrative state of the loadbalancer.
    A valid value is true (UP) or false (DOWN).

* `tags` - (Optional) The key/value pairs to associate with the loadbalancer.

* `flavor` - (Optional) The UUID of a flavor. Currently, this is not supported.
    Changing this creates a new loadbalancer.

* `loadbalancer_provider` - (Optional) The name of the provider. Currently, only
    vlb is supported. Changing this creates a new loadbalancer.

* `security_group_ids` - (Optional) A list of security group IDs to apply to the
    loadbalancer. The security groups must be specified by ID and not name (as
    opposed to how they are configured with the Compute Instance).

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `vip_subnet_id` - See Argument Reference above.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `tenant_id` - See Argument Reference above.
* `vip_address` - See Argument Reference above.
* `admin_state_up` - See Argument Reference above.
* `tags` - See Argument Reference above.
* `flavor` - See Argument Reference above.
* `loadbalancer_provider` - See Argument Reference above.
* `security_group_ids` - See Argument Reference above.
* `vip_port_id` - The Port ID of the Load Balancer IP.
