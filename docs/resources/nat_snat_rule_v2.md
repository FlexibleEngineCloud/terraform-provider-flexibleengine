---
subcategory: "NAT Gateway (NAT)"
description: ""
page_title: "flexibleengine_nat_snat_rule_v2"
---

# flexibleengine_nat_snat_rule_v2

Manages a V2 SNAT rule resource within FlexibleEngine.

## Example Usage

### SNAT rule in VPC scenario

```hcl
resource "flexibleengine_nat_snat_rule_v2" "snat_1" {
  nat_gateway_id = var.natgw_id
  floating_ip_id = var.publicip_id
  subnet_id      = var.subent_id
}
```

### SNAT rule in Direct Connect scenario

```hcl
resource "flexibleengine_nat_snat_rule_v2" "snat_2" {
  nat_gateway_id = var.natgw_id
  floating_ip_id = var.publicip_id
  source_type    = 1
  cidr           = "192.168.10.0/24"
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

* `subnet_id` - (Optional) ID of the subnet this snat rule connects to.
    This parameter and `cidr` are alternative. Changing this creates a new snat rule.

* `cidr` - (Optional) Specifies CIDR, which can be in the format of a network segment or a host IP address.
    This parameter and `subnet_id` are alternative. Changing this creates a new snat rule.

* `source_type` - (Optional) Specifies the scenario. The valid value is 0 (VPC scenario) and 1 (Direct Connect scenario).
    Only `cidr` can be specified over a Direct Connect connection.
    If no value is entered, the default value 0 (VPC scenario) is used.
    Changing this creates a new snat rule.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
* `floating_ip_address` - The actual floating IP address.
* `status` - The status of the snat rule.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minute.
* `delete` - Default is 10 minute.

## Import

SNAT rules can be imported using the following format:

```shell
terraform import flexibleengine_nat_snat_rule_v2.snat_1 9e0713cb-0a2f-484e-8c7d-daecbb61dbe4
```
