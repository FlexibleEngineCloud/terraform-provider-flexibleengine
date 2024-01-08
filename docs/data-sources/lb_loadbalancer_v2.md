---
subcategory: "Deprecated"
description: ""
page_title: "flexibleengine_lb_loadbalancer_v2"
---

# flexibleengine_lb_loadbalancer_v2

Use this data source to get a specific elb loadbalancer within FlexibleEngine.

## Example Usage

```hcl
variable "lb_name" {}

data "flexibleengine_lb_loadbalancer_v2" "test" {
  name = var.lb_name
}
```

## Argument Reference

* `region` - (Optional, String) The region in which to query the data source. If omitted, the provider-level region
  will be used.

* `name` - (Optional, String) Specifies the name of the load balancer.

* `id` - (Optional, String) Specifies the data source ID of the load balancer in UUID format.

* `description` - (Optional, String) Specifies the supplementary information about the load balancer.

* `vip_subnet_id` - (Optional, String) Specifies the ID of the subnet where the load balancer works.

* `vip_address` - (Optional, String) Specifies the private IP address of the load balancer.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `vip_port_id` - The ID of the port bound to the private IP address of the load balancer.

* `status` - The operating status of the load balancer.

* `tags` - The tags associated with the load balancer.
