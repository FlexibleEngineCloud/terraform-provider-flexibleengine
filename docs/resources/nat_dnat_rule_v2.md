---
subcategory: "NAT Gateway (NAT)"
---

# flexibleengine_nat_dnat_rule_v2

Manages a DNAT rule resource within FlexibleEngine.

## Example Usage

### DNAT rule in VPC scenario
```hcl
resource "flexibleengine_compute_instance_v2" "instance_1" {
  ...
}

resource "flexibleengine_nat_dnat_rule_v2" "dnat_1" {
  nat_gateway_id        = var.natgw_id
  floating_ip_id        = var.publicip_id
  port_id               = flexibleengine_compute_instance_v2.instance_1.network[0].port
  protocol              = "tcp"
  internal_service_port = 23
  external_service_port = 8023
}
```

### DNAT rule in Direct Connect scenario
```hcl
resource "flexibleengine_nat_dnat_rule_v2" "dnat_2" {
  nat_gateway_id        = var.natgw_id
  floating_ip_id        = var.publicip_id
  private_ip            = "10.0.0.12"
  protocol              = "tcp"
  internal_service_port = 80
  external_service_port = 8080
}
```

## Argument Reference

The following arguments are supported:

* `nat_gateway_id` - (Required) ID of the nat gateway this dnat rule belongs to.
   Changing this creates a new dnat rule.

* `floating_ip_id` - (Required) Specifies the ID of the floating IP address.
  Changing this creates a new resource.

* `internal_service_port` - (Required) Specifies port used by ECSs or BMSs
  to provide services for external systems. Changing this creates a new resource.

* `external_service_port` - (Required) Specifies port used by ECSs or
  BMSs to provide services for external systems.
  Changing this creates a new dnat rule.

* `port_id` - (Optional) Specifies the port ID of an ECS or a BMS.
  This parameter and `private_ip` are alternative. Changing this creates a
  new dnat rule.

* `private_ip` - (Optional) Specifies the private IP address of a
  user, for example, the IP address of a VPC for dedicated connection.
  This parameter and `port_id` are alternative.
  Changing this creates a new dnat rule.

* `protocol` - (Required) Specifies the protocol type. Currently,
  TCP, UDP, and ANY are supported. The protocol number of TCP, UDP,
  and ANY is 6, 17, and 0, respectively.
  Changing this creates a new dnat rule.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `created_at` - DNAT rule creation time.

* `status` - DNAT rule status.

## Import

DNAT can be imported using the following format:

```
$ terraform import flexibleengine_nat_dnat_rule_v2.dnat_1 f4f783a7-b908-4215-b018-724960e5df4a
```
