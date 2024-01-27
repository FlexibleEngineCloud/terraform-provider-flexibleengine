---
subcategory: "Elastic Load Balance (ELB)"
---

# flexibleengine_lb_loadbalancer

Use this data source to get a specific elb loadbalancer within FlexibleEngine.

## Example Usage

```hcl
variable "lb_name" {}

data "flexibleengine_lb_loadbalancer" "test" {
  name = var.lb_name
}
```

## Argument Reference

* `region` - (Optional, String) Specifies the region in which to obtain the load balancer. If omitted, the
  provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the load balancer.

* `id` - (Optional, String) Specifies the data source ID of the load balancer in UUID format.

* `status` - (Optional, String) Specifies the operating status of the load balancer. Valid values are *ONLINE* and
  *FROZEN*.

* `description` - (Optional, String) Specifies the supplementary information about the load balancer.

* `vip_address` - (Optional, String) Specifies the private IP address of the load balancer.

* `vip_subnet_id` - (Optional, String) Specifies the **IPv4 subnet ID** of the subnet where the load balancer works.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project id of the load balancer.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `tags` - The tags associated with the load balancer.

* `vip_port_id` - The ID of the port bound to the private IP address of the load balancer.

* `public_ip` - The EIP address that is associated to the Load Balancer instance.
