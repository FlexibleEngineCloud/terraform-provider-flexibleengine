---
subcategory: "Bare Metal Server (BMS)"
---

# Data Source: flexibleengine_compute_bms_nic_v2

`flexibleengine_compute_bms_nic_v2` used to query information about a BMS NIC based on the NIC ID.


## Example Usage

```hcl
    
    variable "bms_id" {}
    variable "nic_id" {}

    data "flexibleengine_compute_bms_nic_v2" "Query_BMS_Nic" 
    {
        server_id = "${var.bms_id}",
        id = "${var.nic_id}",
    }
       
```

## Argument Reference

The arguments of this data source act as filters for querying the BMSs details.

* `server_id` - (Required) - This is the unique BMS id.

* `id` - (Optional) - The ID of the NIC.

* `status` - (Optional) - The NIC port status.

## Attributes Reference

All of the argument attributes are also exported as result attributes. 

* `mac_address` - It is NIC's mac address.

* `fixed_ips` - The NIC IP address.

* `network_id` - The ID of the network to which the NIC port belongs.

