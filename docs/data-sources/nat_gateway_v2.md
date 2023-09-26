---
subcategory: "NAT Gateway (NAT)"
---

# flexibleengine_nat_gateway_v2

Use this data source to get the information of an available FlexibleEngine NAT gateway.

## Example Usage

```hcl
data "flexibleengine_nat_gateway_v2" "natgateway" {
  name = "test_natgateway"
}
```

## Argument Reference

* `region` - (Optional, String) The region in which to query the data source. If omitted, the provider-level region
  will be used.

* `id` - (Optional, String) Specifies the ID of the NAT gateway.

* `name` - (Optional, String) Specifies the NAT gateway name. The name can contain only digits, letters,
  underscores (_), and hyphens(-).

* `vpc_id` - (Optional, String) Specifies the ID of the VPC this NAT gateway belongs to.

* `subnet_id` - (Optional, String) Specifies the ID of the VPC Subnet of the downstream interface
  (the next hop of the DVR) of the NAT gateway.

* `spec` - (Optional, String) The NAT gateway type. The value can be:
  + `1`: small type, which supports up to 10,000 SNAT connections.
  + `2`: medium type, which supports up to 50,000 SNAT connections.
  + `3`: large type, which supports up to 200,000 SNAT connections.
  + `4`: extra-large type, which supports up to 1,000,000 SNAT connections.

* `description` - (Optional, String) Specifies the description of the NAT gateway. The value contains 0 to 255
  characters, and angle brackets (<) and (>) are not allowed.

* `status` - (Optional, String) Specifies the status of the NAT gateway.

## Attribute Reference

All the arguments above can also be exported attributes.
