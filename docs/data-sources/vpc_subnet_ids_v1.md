---
subcategory: "Virtual Private Cloud (VPC)"
---

# flexibleengine_vpc_subnet_ids_v1

Provides a list of subnet ids for a vpc_id

This data source can be useful for getting back a list of subnet ids for a vpc.

## Example Usage

The following example shows outputing all cidr blocks for every subnet id in a vpc.

 ```hcl
data "flexibleengine_vpc_subnet_ids_v1" "subnet_ids" {
  vpc_id = var.vpc_id
}

data "flexibleengine_vpc_subnet_v1" "subnet" {
  count = length(data.flexibleengine_vpc_subnet_ids_v1.subnet_ids.ids)
  id    = data.flexibleengine_vpc_subnet_ids_v1.subnet_ids.ids[count.index]
 }

output "subnet_cidr_blocks" {
  value = data.flexibleengine_vpc_subnet_v1.subnet.*.cidr
}
 ```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region in which to obtain the V1 VPC client. A VPC client is needed to retrieve VPCs.
  If omitted, the region argument of the provider is used.

* `vpc_id` (Required, String) - Specifies the VPC ID used as the query filter.

## Attributes Reference

The following attributes are exported:

* `id` - Specifies a data source ID in UUID format.

* `ids` - A list of all the subnet ids found. This data source will fail if none are found.
