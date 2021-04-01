---
subcategory: "Bare Metal Server (BMS)"
---


# Data Source: flexibleengine_compute_bms_server_v2

`flexibleengine_compute_bms_server_v2` used to query a BMS or BMSs details.

## Example Usage

```hcl

    variable "bms_id" {}
    variable "bms_name" {}

    data "flexibleengine_compute_bms_server_v2" "Query_BMS" 
    {
        id = "${var.bms_id}",
        name = "${var.bms_name}"     
    }

```

## Argument Reference

The arguments of this data source act as filters for querying the BMSs details.

* `id` - (Optional) - The unique ID of the BMS.

* `user_id` (Optional) - The ID of the user to which the BMS belongs.

* `name` (Optional) - The name of BMS.

* `status` (Optional) - The BMS status.

* `host_status` (Optional) - The nova-compute status: **UP, UNKNOWN, DOWN, MAINTENANCE** and **Null**.

* `key_name` (Optional) - It is the SSH key name.

* `flavor_id` (Optional) - It gives the BMS flavor information.

* `image_id` (Optional) - The BMS image.


## Attributes Reference

All of the argument attributes are also exported as result attributes. 

* `host_id` - 	It is the host ID of the BMS.

* `progress` - This is a reserved attribute.

* `metadata` -  The BMS metadata is specified.

* `access_ip_v4` -  This is a reserved attribute.

* `access_ip_v6` - This is a reserved attribute.  

* `addresses` - It gives the BMS network address.

* `security_groups` - The list of security groups to which the BMS belongs.

* `tags` - Specifies the BMS tag.

* `locked` -  It specifies whether a BMS is locked, true: The BMS is locked, false: The BMS is not locked.

* `config_drive` -  This is a reserved attribute.

* `availability_zone` - Specifies the AZ ID.

* `description` -  Provides supplementary information about the pool.

* `kernel_id` - The UUID of the kernel image when the AMI image is used.

* `hypervisor_hostname` -  It is the name of a host on the hypervisor.

* `instance_name` - Instance name is specified.