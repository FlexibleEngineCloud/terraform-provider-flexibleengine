---
subcategory: "Virtual Private Cloud (VPC)"
---

# flexibleengine_networking_secgroup_rule_v2

Use this data source to get the ID of an available FlexibleEngine security group rule.

## Example Usage

```hcl
data "flexibleengine_networking_secgroup_rule_v2" "secgroup_rule" {
  name = "tf_test_secgroup"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V2 Neutron client.
  A Neutron client is needed to retrieve security groups ids. If omitted, the
  `region` argument of the provider is used.

* `id` - (Optional) The ID of the security group.

* `description` - (Optional) The description of the security group rule.

* `security_group_id` - (Optional) The security group ID the rule belongs to.
  
* `ethertype` - (Optional) The layer 3 protocol type.

* `protocol` - (Optional) The layer 4 protocol type.

* `port_range_min` - (Optional) The lower part of the allowed port range.

* `port_range_max` - (Optional) The higher part of the allowd port range.

* `remote_ip_prefix` - (Optional) The remote CICR.

* `remote_group_id` - (Optional) The remote group id.

* `tenant_id` - (Optional) The owner of the security group.

## Attributes Reference

`id` is set to the ID of the found security group rule. In addition, the following
attributes are exported:

* `region` - See Argument Reference above.
* `description` - See Argument Reference above.
* `security_group_id` - See Argument Reference above.
* `ethertype` - See Argument Reference above.
* `protocol` - See Argument Reference above.
* `port_range_min` - See Argument Reference above.
* `port_range_max` - See Argument Reference above.
* `remote_ip_prefix` - See Argument Reference above.
* `remote_group_id` - See Argument Reference above.
* `remote ip_prefix` - See Argument Reference above.
  
