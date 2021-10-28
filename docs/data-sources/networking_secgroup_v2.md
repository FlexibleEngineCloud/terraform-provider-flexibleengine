---
subcategory: "Virtual Private Cloud (VPC)"
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

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `description`- The description of the security group.
