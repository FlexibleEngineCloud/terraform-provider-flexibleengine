---
subcategory: "Deprecated"
---

# flexibleengine_elb_loadbalancer

!> **Warning:** Classic load balancers are no longer provided, using elastic load balancers instead.

Manages a **classic** load balancer resource within FlexibleEngine.

## Example Usage

### External Load Balancer

```hcl
resource "flexibleengine_elb_loadbalancer" "elb" {
  type        = "External"
  name        = "elb-external"
  description = "external elb"
  vpc_id      = var.vpc_id
  bandwidth   = 5
}
```

### Internal Load Balancer

```hcl
resource "flexibleengine_elb_loadbalancer" "elb" {
  type              = "Internal"
  name              = "elb-internal"
  description       = "internal elb"
  az                = "eu-west-0"
  vpc_id            = var.vpc_id
  vip_subnet_id     = var.subnet_id
  security_group_id = var.sec_group
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to create the loadbalancer. If
    omitted, the `region` argument of the provider is used. Changing this
    creates a new loadbalancer.

* `name` - (Required) Specifies the load balancer name. The name is a string
    of 1 to 64 characters that consist of letters, digits, underscores (_),
    and hyphens (-).

* `type` - (Required) Specifies the load balancer type. The value can be
    Internal or External.

* `vpc_id` - (Required) Specifies the VPC ID.

* `description` - (Optional) Provides supplementary information about the
    listener. The value is a string of 0 to 128 characters and cannot be <>.

* `vip_address` - (Optional) Specifies the IP address provided by ELB.
    When type is set to External, the value of this parameter is the elastic
    IP address. When type is set to Internal, the value of this parameter is
    the private network IP address. You can select an existing elastic IP address
    and create a public network load balancer. When this parameter is configured,
    parameter `bandwidth` is invalid.

* `bandwidth` - (Optional) Specifies the bandwidth (Mbit/s). This parameter
    is valid when type is set to External, and it is invalid when type
    is set to Internal. The value ranges from 1 to 300.

* `vip_subnet_id` - (Optional) Specifies the ID of the private network
    to be added. This parameter is mandatory when type is set to Internal,
    and it is invalid when type is set to External.

* `security_group_id` - (Optional) Specifies the security group ID. The
    value is a string of 1 to 200 characters that consists of uppercase and
    lowercase letters, digits, and hyphens (-). This parameter is mandatory
    when type is set to Internal, and it is invalid when type is set to External.

* `az` - (Optional) Specifies the ID of the availability zone (AZ). This
    parameter is mandatory when type is set to Internal, and it is invalid
    when type is set to External.

* `admin_state_up` - (Optional) Specifies the status of the load balancer. Defaults to true.
    + true: indicates that the load balancer is running.
    + false: indicates that the load balancer is stopped.

* `tenantid` - (Optional) Specifies the tenant ID. This parameter is mandatory
    only when type is set to Internal.

## Attributes Reference

The following attributes are exported:

* `id` - Specifies the load balancer ID.
* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `vpc_id` - See Argument Reference above.
* `bandwidth` - See Argument Reference above.
* `type` - See Argument Reference above.
* `admin_state_up` - See Argument Reference above.
* `vip_subnet_id` - See Argument Reference above.
* `az` - See Argument Reference above.
* `security_group_id` - See Argument Reference above.
* `vip_address` - See Argument Reference above.
* `tenantid` - See Argument Reference above.
