---
subcategory: "Elastic Load Balance (Dedicated ELB)"
description: ""
page_title: "flexibleengine_elb_ipgroup"
---

# flexibleengine_elb_ipgroup

Manages an ELB IP Group resource within FlexibleEngine.

## Example Usage

```hcl
resource "flexibleengine_elb_ipgroup" "basic" {
  name        = "basic"
  description = "basic example"

  ip_list {
    ip          = "192.168.10.10"
    description = "ECS01"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the ip group resource. If omitted, the
  provider-level region will be used. Changing this creates a new ip group.

* `name` - (Required, String) Specifies the name of the ip group.

* `description` - (Optional, String) Specifies the description of the ip group.

* `ip_list` - (Required, List) Specifies an array of one or more ip addresses. The [ip_list](#elb_ip_list) object
  structure is documented below.

<a name="elb_ip_list"></a>
The `ip_list` block supports:

* `ip` - (Required, String) IP address or CIDR block.

* `description` - (Optional, String) Human-readable description for the ip.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The uuid of the ip group.

## Import

ELB IP group can be imported using the IP group ID, e.g.

```shell
terraform import flexibleengine_elb_ipgroup.group_1 5c20fdad-7288-11eb-b817-0255ac10158b
```
