---
subcategory: "Bare Metal Server (BMS)"
---

# Data Source: flexibleengine_compute_bms_flavors_v2

Use this data source to get an available BMS Flavor.

## Example Usage

```hcl
data "flexibleengine_compute_bms_flavors_v2" "BMS_flavor" {
  vcpus = 32
}
```

## Argument Reference

The arguments of this data source act as filters for querying the BMSs details.

* `name` (Optional) - Specifies the name of the BMS flavor.

* `vcpus` (Optional) - Specifies the number of CPU cores in the BMS flavor.

* `min_ram` (Optional) - Specifies the minimum memory size in MB. Only the BMSs with the memory size
  greater than or equal to the minimum size can be queried.

* `min_disk` (Optional) - Specifies the minimum disk size in GB. Only the BMSs with a disk size
  greater than or equal to the minimum size can be queried.

* `sort_key` (Optional) - The sorting field. The default value is **flavorid**.
  The available values are **name**, **memory_mb**, **vcpus**, **root_gb**, or **flavorid**.

* `sort_dir` (Optional) - The sorting order, which can be **asc** (ascending) or **desc** (descending).
  The default value is **asc**.

## Attributes Reference

All of the argument attributes are also exported as result attributes.

* `id` - The BMS flavor id.

* `ram` - The memory size (in MB) of the BMS flavor.

* `disk` - The disk size (GB) in the BMS flavor.

* `swap` -  This is a reserved attribute.

* `rx_tx_factor` - This is a reserved attribute.
