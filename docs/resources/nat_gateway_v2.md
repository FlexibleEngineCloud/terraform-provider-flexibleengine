---
subcategory: "NAT Gateway (NAT)"
description: ""
page_title: "flexibleengine_nat_gateway_v2"
---

# flexibleengine_nat_gateway_v2

Manages a V2 nat gateway resource within FlexibleEngine.

## Example Usage

```hcl
resource "flexibleengine_nat_gateway_v2" "nat_1" {
  name        = "nat_test"
  description = "test for terraform"
  spec        = "3"
  vpc_id      = flexibleengine_vpc_v1.example_vpc.id
  subnet_id   = flexibleengine_vpc_subnet_v1.example_subnet.id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the Nat gateway resource.
  If omitted, the provider-level region will be used. Changing this creates a new nat gateway.

* `name` - (Required, String) Specifies the nat gateway name. The name can contain only digits, letters,
  underscores (_), and hyphens(-).

* `spec` - (Required, String) Specifies the nat gateway type. The value can be:
  + `1`: small type, which supports up to 10,000 SNAT connections.
  + `2`: medium type, which supports up to 50,000 SNAT connections.
  + `3`: large type, which supports up to 200,000 SNAT connections.
  + `4`: extra-large type, which supports up to 1,000,000 SNAT connections.

* `vpc_id` - (Required, String, ForceNew) Specifies the ID of the VPC this nat gateway belongs to.
  Changing this creates a new nat gateway.

* `subnet_id` - (Required, String, ForceNew) Specifies the ID of the VPC Subnet of the downstream interface
  (the next hop of the DVR) of the NAT gateway. Changing this creates a new nat gateway.

* `description` - (Optional, String) Specifies the description of the nat gateway.
  The value contains 0 to 255 characters, and angle brackets (<) and (>) are not allowed.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `status` - The status of the nat gateway.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minute.
* `delete` - Default is 10 minute.

## Import

Nat gateway can be imported using the following format:

```shell
terraform import flexibleengine_nat_gateway_v2.nat_1 d126fb87-43ce-4867-a2ff-cf34af3765d9
```
