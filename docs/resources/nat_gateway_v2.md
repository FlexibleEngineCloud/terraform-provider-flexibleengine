---
subcategory: "NAT Gateway (NAT)"
---

# flexibleengine_nat_gateway_v2

Manages a V2 nat gateway resource within FlexibleEngine.

## Example Usage

```hcl
resource "flexibleengine_nat_gateway_v2" "nat_1" {
  name   = "nat_test"
  description = "test for terraform"
  spec = "3"
  router_id = "2c1fe4bd-ebad-44ca-ae9d-e94e63847b75"
  internal_network_id = "dc8632e2-d9ff-41b1-aa0c-d455557314a0"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 nat client.
    If omitted, the `region` argument of the provider is used. Changing this
    creates a new nat gateway.

* `name` - (Required) The name of the nat gateway.

* `description` - (Optional) The description of the nat gateway.

* `spec` - (Required) The specification of the nat gateway, valid values are "1",
    "2", "3", "4" (for Small, Medium, Large, Extra-Large)

* `router_id` - (Required) ID of the router/VPC this nat gateway belongs to. Changing
    this creates a new nat gateway.

* `internal_network_id` - (Required) ID of the subnet (!) this nat gateway connects to.
    Changing this creates a new nat gateway.

* `tenant_id` - (Optional) The target tenant/project ID in which to allocate the nat
    gateway. Changing this creates a new nat gateway .

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `spec` - See Argument Reference above.
* `tenant_id` - See Argument Reference above.
* `router_id` - See Argument Reference above.
* `internal_network_id` - See Argument Reference above.
