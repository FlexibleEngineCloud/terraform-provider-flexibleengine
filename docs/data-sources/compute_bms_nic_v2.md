---
subcategory: "Bare Metal Server (BMS)"
---

# flexibleengine_compute_bms_nic_v2

`flexibleengine_compute_bms_nic_v2` used to query information about a BMS NIC based on the NIC ID.

## Example Usage

```hcl
variable "bms_id" {}
variable "nic_id" {}

data "flexibleengine_compute_bms_nic_v2" "nic" {
  server_id = var.bms_id
  id        = var.nic_id
}
```

## Argument Reference

The arguments of this data source act as filters for querying the BMSs details.

* `region` - (Optional, String) The region in which to query the data source. If omitted, the provider-level region
  will be used.

* `server_id` - (Required, String) - This is the unique BMS id.

* `id` - (Optional, String) - The ID of the NIC.

* `status` - (Optional, String) - The NIC port status.

## Attribute Reference

All of the argument attributes are also exported as result attributes.

* `mac_address` - It is NIC's mac address.

* `fixed_ips` - The NIC IP address.
    The [fixed_ips](#<a name="bms_fixed_ips"></a>) object structure is documented below.

* `network_id` - The ID of the network to which the NIC port belongs.

<a name="bms_fixed_ips"></a>
The `fixed_ips` block supports:

* `ip_address` - Specifies the NIC private IP address.

* `subnet_id` - Specifies the ID of the subnet (subnet_id) corresponding to the private IP address of the NIC.
