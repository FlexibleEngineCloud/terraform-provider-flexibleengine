---
subcategory: "Elastic Load Balance (ELB)"
description: ""
page_title: "flexibleengine_lb_loadbalancer_v2"
---

# flexibleengine_lb_loadbalancer_v2

Manages an **enhanced** load balancer resource within FlexibleEngine.

## Example Usage

```hcl
resource "flexibleengine_lb_loadbalancer_v2" "lb_1" {
  vip_subnet_id = flexibleengine_vpc_subnet_v1.example_subnet.ipv4_subnet_id

  tags = {
    key = "value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the loadbalancer resource.
  If omitted, the `region` argument of the provider is used. Changing this creates a new loadbalancer.

* `vip_subnet_id` - (Required, String, ForceNew) The `ipv4_subnet_id` or `ipv6_subnet_id` of the
  VPC Subnet on which to allocate the loadbalancer's address.
  A tenant can only create Loadbalancers on networks authorized
  by policy (e.g. networks that belong to them or networks that
  are shared).  Changing this creates a new loadbalancer.

* `name` - (Optional, String) Human-readable name for the loadbalancer. Does not have to be unique.

* `description` - (Optional, String) Human-readable description for the loadbalancer.

* `vip_address` - (Optional, String, ForceNew) The ip address of the load balancer.
  Changing this creates a new loadbalancer.

* `tags` - (Optional, Map) The key/value pairs to associate with the loadbalancer.

* `loadbalancer_provider` - (Optional, String, ForceNew) The name of the provider. Currently, only
  vlb is supported. Changing this creates a new loadbalancer.

* `security_group_ids` - (Optional, List) A list of security group IDs to apply to the
  loadbalancer. The security groups must be specified by ID and not name (as
  opposed to how they are configured with the Compute Instance).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `vip_port_id` - The Port ID of the Load Balancer IP.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 5 minutes.

## Import

Loadbalancers can be imported using the `id`, e.g.

```shell
terraform import flexibleengine_lb_loadbalancer_v2.loadbalancer_1 3e3632db-36c6-4b28-a92e-e72e6562daa6
```
