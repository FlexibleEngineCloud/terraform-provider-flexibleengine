---
layout: "flexibleengine"
page_title: "FlexibleEngine: flexibleengine_networking_secgroup_v2"
sidebar_current: "docs-flexibleengine-datasource-networking-secgroup-v2"
description: |-
  Get information on an FlexibleEngine Security Group.
---

# flexibleengine\_networking\_secgroup\_v2

Use this data source to get the ID of an available FlexibleEngine security group.

## Example Usage

```hcl
data "flexibleengine_networking_secgroup_v2" "secgroup" {
  name = "tf_test_secgroup"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V2 Neutron client.
  A Neutron client is needed to retrieve security groups ids. If omitted, the
  `region` argument of the provider is used.

* `secgroup_id` - (Optional) The ID of the security group.

* `name` - (Optional) The name of the security group.

* `tenant_id` - (Optional) The owner of the security group.

## Attributes Reference

`id` is set to the ID of the found security group. In addition, the following
attributes are exported:

* `name` - See Argument Reference above.
* `description`- The description of the security group.
* `region` - See Argument Reference above.
