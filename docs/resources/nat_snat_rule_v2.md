---
subcategory: "NAT Gateway (NAT)"
---

# flexibleengine\_nat\_snat\_rule_v2

Manages a V2 snat rule resource within FlexibleEngine Nat

## Example Usage

```hcl
resource "flexibleengine_nat_snat_rule_v2" "snat_1" {
  nat_gateway_id = "3c0dffda-7c76-452b-9dcc-5bce7ae56b17"
  network_id = "dc8632e2-d9ff-41b1-aa0c-d455557314a0"
  floating_ip_id = "0a166fc5-a904-42fb-b1ef-cf18afeeddca"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 nat client.
    If omitted, the `region` argument of the provider is used. Changing this
    creates a new snat rule.

* `nat_gateway_id` - (Required) ID of the nat gateway this snat rule belongs to.
    Changing this creates a new snat rule.

* `floating_ip_id` - (Required) ID of the floating ip this snat rule connets to.
    Changing this creates a new snat rule.

* `network_id` - (Optional) ID of the network this snat rule connects to.
    This parameter and `cidr` are alternative. Changing this creates a new snat rule.

* `cidr` - (Optional) Specifies CIDR, which can be in the format of a network segment or a host IP address.
    This parameter and `network_id` are alternative. Changing this creates a new snat rule.

* `source_type` - (Optional) Specifies the scenario. The valid value is 0 (VPC scenario) and 1 (Direct Connect scenario).
    Only `cidr` can be specified over a Direct Connect connection.
    If no value is entered, the default value 0 (VPC scenario) is used.
    Changing this creates a new snat rule.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `nat_gateway_id` - See Argument Reference above.
* `network_id` - See Argument Reference above.
* `cidr` - See Argument Reference above.
* `source_type` - See Argument Reference above.
* `floating_ip_id` - See Argument Reference above.
* `floating_ip_address` - The actual floating IP address.
* `status` -  The status of the SNAT rule.
