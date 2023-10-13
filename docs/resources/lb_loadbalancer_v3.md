---
subcategory: "Elastic Load Balance (Dedicated ELB)"
description: ""
page_title: "flexibleengine_lb_loadbalancer_v3"
---

# flexibleengine_lb_loadbalancer_v3

Manages a **Dedicated** Load Balancer resource within FlexibleEngine.

## Example Usage

### Basic Loadbalancer

```hcl
data "flexibleengine_elb_flavors" "l7_flavors" {
  type            = "L7"
}

data "flexibleengine_elb_flavors" "l4_flavors" {
  type            = "L4"
}

resource "flexibleengine_vpc_v1" "example_vpc" {
  name = "example-vpc"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "example_subnet" {
  name       = "example-vpc-subnet"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = flexibleengine_vpc_v1.example_vpc.id
}

resource "flexibleengine_lb_loadbalancer_v3" "basic" {
  name              = "basic"
  description       = "basic example"
  cross_vpc_backend = true

  vpc_id         = flexibleengine_vpc_v1.example_vpc.id
  ipv4_subnet_id = flexibleengine_vpc_subnet_v1.example_subnet.ipv4_subnet_id

  l4_flavor_id = data.flexibleengine_elb_flavors.l4_flavors.ids[0]
  l7_flavor_id = data.flexibleengine_elb_flavors.l7_flavors.ids[0]

  availability_zone = [
    "eu-west-0a",
    "eu-west-0b",
  ]
}
```

### Loadbalancer With Existing EIP

```hcl
data "flexibleengine_elb_flavors" "l7_flavors" {
  type            = "L7"
}

data "flexibleengine_elb_flavors" "l4_flavors" {
  type            = "L4"
}

resource "flexibleengine_vpc_v1" "example_vpc" {
  name = "example-vpc"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "example_subnet" {
  name       = "example-vpc-subnet"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = flexibleengine_vpc_v1.example_vpc.id
}

resource "flexibleengine_lb_loadbalancer_v3" "basic" {
  name              = "basic"
  description       = "basic example"
  cross_vpc_backend = true

  vpc_id            = flexibleengine_vpc_v1.example_vpc.id
  ipv6_network_id   = flexibleengine_vpc_subnet_v1.example_subnet_ipv6.id
  ipv6_bandwidth_id = "{{ ipv6_bandwidth_id }}"
  ipv4_subnet_id    = flexibleengine_vpc_subnet_v1.example_subnet.ipv4_subnet_id

  l4_flavor_id = data.flexibleengine_elb_flavors.l4_flavors.ids[0]
  l7_flavor_id = data.flexibleengine_elb_flavors.l7_flavors.ids[0]

  availability_zone = [
    "eu-west-0a",
    "eu-west-0b",
  ]

  ipv4_eip_id = flexibleengine_vpc_eip.example_eip.id
}
```

### Loadbalancer With EIP

```hcl
data "flexibleengine_elb_flavors" "l7_flavors" {
  type            = "L7"
}

data "flexibleengine_elb_flavors" "l4_flavors" {
  type            = "L4"
}

resource "flexibleengine_vpc_v1" "example_vpc" {
  name = "example-vpc"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "example_subnet" {
  name       = "example-vpc-subnet"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = flexibleengine_vpc_v1.example_vpc.id
}

resource "flexibleengine_lb_loadbalancer_v3" "basic" {
  name              = "basic"
  description       = "basic example"
  cross_vpc_backend = true

  vpc_id            = flexibleengine_vpc_v1.example_vpc.id
  ipv6_network_id   = flexibleengine_vpc_subnet_v1.example_subnet_ipv6.id
  ipv6_bandwidth_id = "{{ ipv6_bandwidth_id }}"
  ipv4_subnet_id    = flexibleengine_vpc_subnet_v1.example_subnet.ipv4_subnet_id

  l4_flavor_id = data.flexibleengine_elb_flavors.l4_flavors.ids[0]
  l7_flavor_id = data.flexibleengine_elb_flavors.l7_flavors.ids[0]

  availability_zone = [
    "eu-west-0a",
    "eu-west-0b",
  ]

  iptype                = "5_bgp"
  bandwidth_charge_mode = "traffic"
  sharetype             = "PER"
  bandwidth_size        = 10
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the loadbalancer resource. If omitted, the
  provider-level region will be used. Changing this creates a new loadbalancer.

* `availability_zone` - (Required, List, ForceNew) Specifies the list of AZ names. Changing this parameter will create a
  new resource.

* `name` - (Required, String) Human-readable name for the loadbalancer.

* `description` - (Optional, String) Human-readable description for the loadbalancer.

* `cross_vpc_backend` - (Optional, Bool) Enable this if you want to associate the IP addresses of backend servers with
  your load balancer. Can only be true when updating.

* `vpc_id` - (Optional, String, ForceNew) The vpc on which to create the loadbalancer. Changing this creates a new
  loadbalancer.

* `ipv4_subnet_id` - (Optional, String) The **IPv4 subnet ID** of the subnet on which to allocate the loadbalancer's
  ipv4 address.

* `ipv6_network_id` - (Optional, String) The network on which to allocate the loadbalancer's ipv6 address.

* `ipv6_bandwidth_id` - (Optional, String) The ipv6 bandwidth id. Only support shared bandwidth.

* `ipv4_address` - (Optional, String) The ipv4 address of the load balancer.

* `ipv4_eip_id` - (Optional, String, ForceNew) The ID of the EIP. Changing this parameter will create a new resource.

-> **NOTE:** If the ipv4_eip_id parameter is configured, you do not need to configure the bandwidth parameters:
`iptype`, `bandwidth_charge_mode`, `bandwidth_size`, `share_type` and `bandwidth_id`.

* `iptype` - (Optional, String, ForceNew) Elastic IP type. Changing this parameter will create a new resource.

* `bandwidth_charge_mode` - (Optional, String, ForceNew) Bandwidth billing type. Value options:
  + **bandwidth**: Billed by bandwidth.
  + **traffic**: Billed by traffic.

  It is mandatory when `iptype` is set and `bandwidth_id` is empty.
  Changing this parameter will create a new resource.

* `sharetype` - (Optional, String, ForceNew) Bandwidth sharing type. Value options:
  + **PER**: Dedicated bandwidth.
  + **WHOLE**: Shared bandwidth.

  It is mandatory when `iptype` is set and `bandwidth_id` is empty.
  Changing this parameter will create a new resource.

* `bandwidth_size` - (Optional, Int, ForceNew) Bandwidth size. It is mandatory when `iptype` is set and `bandwidth_id`
  is empty. Changing this parameter will create a new resource.

* `bandwidth_id` - (Optional, String, ForceNew) Bandwidth ID of the shared bandwidth. It is mandatory when `sharetype`
  is **WHOLE**. Changing this parameter will create a new resource.

  -> **NOTE:** If the `bandwidth_id` parameter is configured, you can not configure the parameters:
  `bandwidth_charge_mode`, `bandwidth_size`.

* `l4_flavor_id` - (Optional, String) The L4 flavor id of the load balancer.

* `l7_flavor_id` - (Optional, String) The L7 flavor id of the load balancer.

* `backend_subnets` - (Optional, List) The IDs of subnets on the downstream plane.
  + If this parameter is not specified, select subnets as follows:
    - If IPv6 is enabled for a load balancer, the ID of subnet specified in `ipv6_network_id` will be used.
    - If IPv4 is enabled for a load balancer, the ID of subnet specified in `ipv4_subnet_id` will be used.
    - If only public network is available for a load balancer, the ID of any subnet in the VPC where the load balancer
      resides will be used. Subnets with more IP addresses are preferred.
  + If there is more than one subnet, the first subnet in the list will be used, and the subnets must be in the VPC
    where the load balancer resides.

* `tags` - (Optional, Map) The key/value pairs to associate with the loadbalancer.

* `autoscaling_enabled` - (Optional, Bool) Specifies whether autoscaling is enabled. Valid values are **true** and
  **false**.

* `min_l7_flavor_id` - (Optional, String) Specifies the ID of the minimum Layer-7 flavor for elastic scaling.
  This parameter cannot be left blank if there are HTTP or HTTPS listeners.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project id of the resource.
  Changing this will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `ipv4_port_id` - The ID of the port bound to the private IPv4 address of the loadbalancer.
* `ipv4_eip` - The ipv4 eip address of the Load Balancer.
* `ipv6_eip` - The ipv6 eip address of the Load Balancer.
* `ipv6_eip_id` - The ipv6 eip id of the Load Balancer.
* `ipv6_address` - The ipv6 address of the Load Balancer.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 5 minutes.

## Import

ELB loadbalancer can be imported using the loadbalancer ID, e.g.

```shell
terraform import flexibleengine_lb_loadbalancer_v3.loadbalancer_1 5c20fdad-7288-11eb-b817-0255ac10158b
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `ipv6_bandwidth_id`, `iptype`,
`bandwidth_charge_mode`, `sharetype`,  `bandwidth_size` and `bandwidth_id`.
It is generally recommended running `terraform plan` after importing a loadbalancer.
You can then decide if changes should be applied to the loadbalancer, or the resource
definition should be updated to align with the loadbalancer. Also you can ignore changes as below.

```hcl
resource "flexibleengine_lb_loadbalancer_v3" "loadbalancer_1" {
  ...
lifecycle {
  ignore_changes = [
    ipv6_bandwidth_id, iptype, bandwidth_charge_mode, sharetype, bandwidth_size, bandwidth_id,
  ]
}
}
```
