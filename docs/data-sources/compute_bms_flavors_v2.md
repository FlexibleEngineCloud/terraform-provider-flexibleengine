---
subcategory: "Bare Metal Server (BMS)"
---

# Data Source: flexibleengine_compute_bms_flavors_v2

`flexibleengine_compute_bms_flavors_v2` used to query flavors of BMSs.

## Example Usage

```hcl
    
    variable "flavor_id" { }
    variable "disk_size" { }

    data "flexibleengine_compute_bms_flavors_v2" "Query_BMS_flavors" 
    {
        id = "${var.bms_id}",
        min_disk = "${var.disk_size}",
        sort_key = "id",
        sort_dir = "desc",
    }
    
```

## Argument Reference

The arguments of this data source act as filters for querying the BMSs details.

* `name` - (Optional) - The name of the BMS flavor.

* `id` (Optional) - The BMS flavor id.

* `min_ram` (Optional) - The minimum memory size in MB. Only the BMSs with the memory size greater than or equal to the minimum size can be queried.

* `min_disk` (Optional) - The minimum disk size in GB. Only the BMSs with a disk size greater than or equal to the minimum size can be queried.

* `sort_key` (Optional) - The sorting field. The default value is **flavorid**. The other values are **name**, **memory_mb**, **vcpus**, **root_gb**, or **flavorid**.

* `sort_dir` (Optional) - The sorting order, which can be **ascending** (**asc**) or **descending** (**desc**). The default value is **asc**.

## Attributes Reference

All of the argument attributes are also exported as result attributes. 

* `ram` - It is the memory size (in MB) of the flavor.

* `vcpus` - It is the number of CPU cores in the BMS flavor.

* `disk` - Specifies the disk size (GB) in the BMS flavor.

* `swap` -  This is a reserved attribute.

* `rx_tx_factor` - This is a reserved attribute.