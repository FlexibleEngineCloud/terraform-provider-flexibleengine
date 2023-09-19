---
subcategory: "Deprecated"
description: ""
page_title: "flexibleengine_elb_loadbalancer"
---

# flexibleengine_elb_loadbalancer

!> **Warning:** Classic load balancers are no longer provided, using elastic load balancers instead.

Manages a **classic** load balancer resource within FlexibleEngine.

## Example Usage

### External Load Balancer

```hcl
resource "flexibleengine_vpc_v1" "example_vpc" {
  name = "example-vpc"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_elb_loadbalancer" "elb" {
  type        = "External"
  name        = "elb-external"
  description = "external elb"
  vpc_id      = flexibleengine_vpc_v1.example_vpc.id
  bandwidth   = 5
}
```

### Internal Load Balancer

```hcl
resource "flexibleengine_vpc_v1" "example_vpc" {
  name = "example-vpc"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_elb_loadbalancer" "elb" {
  type              = "Internal"
  name              = "elb-internal"
  description       = "internal elb"
  az                = "eu-west-0"
  vpc_id            = flexibleengine_vpc_v1.example_vpc.id
  vip_subnet_id     = var.subnet_id
  security_group_id = var.sec_group
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the loadbalancer. If
  omitted, the `region` argument of the provider is used. Changing this creates a new loadbalancer.

* `name` - (Optional, String) Specifies the load balancer name. The name is a string
  of 1 to 64 characters that consist of letters, digits, underscores (_), and hyphens (-).

* `type` - (Required, String, ForceNew) Specifies the load balancer type. The value can be
  Internal or External. Changing this creates a new loadbalancer.

* `vpc_id` - (Required, String, ForceNew) Specifies the VPC ID. Changing this creates a new loadbalancer.

* `description` - (Optional, String) Provides supplementary information about the
  listener. The value is a string of 0 to 128 characters and cannot be <>.

* `vip_address` - (Optional, String, ForceNew) Specifies the IP address provided by ELB.
  When type is set to External, the value of this parameter is the elastic
  IP address. When type is set to Internal, the value of this parameter is
  the private network IP address. You can select an existing elastic IP address
  and create a public network load balancer. When this parameter is configured,
  parameter `bandwidth` is invalid. Changing this creates a new loadbalancer.

* `bandwidth` - (Optional, Int) Specifies the bandwidth (Mbit/s). This parameter
  is valid when type is set to External, and it is invalid when type
  is set to Internal. The value ranges from 1 to 300.

* `vip_subnet_id` - (Optional, String, ForceNew) Specifies the ID of the private network
  to be added. This parameter is mandatory when type is set to Internal,
  and it is invalid when type is set to External. Changing this creates a new loadbalancer.

* `security_group_id` - (Optional, String, ForceNew) Specifies the security group ID. The
  value is a string of 1 to 200 characters that consists of uppercase and
  lowercase letters, digits, and hyphens (-). This parameter is mandatory
  when type is set to Internal, and it is invalid when type is set to External.
  Changing this creates a new loadbalancer.

* `az` - (Optional, String, ForceNew) Specifies the ID of the availability zone (AZ). This
  parameter is mandatory when type is set to Internal, and it is invalid
  when type is set to External. Changing this creates a new loadbalancer.

* `tenantid` - (Optional, String, ForceNew) Specifies the tenant ID. This parameter is mandatory
  only when type is set to Internal. Changing this creates a new loadbalancer.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies the load balancer ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 5 minutes.
