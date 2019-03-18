---
layout: "flexibleengine"
page_title: "FlexibleEngine: flexibleengine_nat_dnat_rule_v2"
sidebar_current: "docs-flexibleengine-resource-nat-dnat-rule-v2"
description: |-

---

# flexibleengine\_nat\_dnat


## Example Usage

### Dnat

```hcl
resource "flexibleengine_nat_dnat_rule_v2" "dnat" {
  floating_ip_id = "bf99c679-9f41-4dac-8513-9c9228e713e1"
  nat_gateway_id = "bf99c679-9f41-4dac-8513-9c9228e713e1"
  internal_service_port = 993
  protocol = "tcp"
  external_service_port = 242
}
```

## Argument Reference

The following arguments are supported:

* `floating_ip_id` - (Required) Specifies the ID of the floating IP address.

* `internal_service_port` - (Required) Specifies port used by ECSs or BMSs
  to provide services for external systems.

* `nat_gateway_id` - (Required) Specifies the ID of the NAT gateway.
  nat gateway id

* `port_id` - (Required) Specifies the port ID of an ECS or a BMS.
  This parameter and private_ip are alternative.

* `private_ip` - (Required) Specifies the private IP address of a
  user, for example, the IP address of a VPC for dedicated connection.
  This parameter and port_id is alternative.

* `protocol` - (Required) Specifies the protocol type. Currently,
  TCP, UDP, and ANY are supported. The protocol number of TCP, UDP,
  and ANY is 6, 17, and 0, respectively.

* `internal_service_port` - (Required) Specifies port used by ECSs or
  BMSs to provide services for external systems.


## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `created_at` - Dnat rule creation time.

* `status` - Dnat rule status.

## Import

Dnat can be imported using the following format:

```
$ terraform import flexibleengine_nat_dnat_rule_v2.default {{ resource id}}
```
