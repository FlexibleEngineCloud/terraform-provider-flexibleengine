---
subcategory: "Bare Metal Server (BMS)"
---


# flexibleengine_compute_bms_server_v2

`flexibleengine_compute_bms_server_v2` used to query a BMS or BMSs details.

## Example Usage

```hcl
variable "bms_name" {}

data "flexibleengine_compute_bms_server_v2" "server" {
  name = var.bms_name
}
```

## Argument Reference

The arguments of this data source act as filters for querying the BMSs details.

* `region` - (Optional, String) The region in which to query the data source. If omitted, the provider-level region
  will be used.

* `id` - (Optional, String) - The unique ID of the BMS.

* `user_id` (Optional, String) - The ID of the user to which the BMS belongs.

* `name` (Optional, String) - The name of BMS.

* `status` (Optional, String) - The BMS status.

* `host_status` (Optional, String) - The nova-compute status: **UP, UNKNOWN, DOWN, MAINTENANCE** and **Null**.

* `key_name` (Optional, String) - It is the SSH key name.

* `flavor_id` (Optional, String) - It gives the BMS flavor information.

* `image_id` (Optional, String) - The BMS image.

## Attribute Reference

All of the argument attributes are also exported as result attributes.

* `host_id` - It is the host ID of the BMS.

* `progress` - This is a reserved attribute.

* `metadata` -  The BMS metadata is specified.

* `access_ip_v4` -  This is a reserved attribute.

* `access_ip_v6` - This is a reserved attribute.  

* `security_groups` - The list of security groups to which the BMS belongs.
    The [security_groups](#<a name="bms_security_groups"></a>) object structure is documented below.

* `tags` - Specifies the BMS tag.

* `locked` -  It specifies whether a BMS is locked, true: The BMS is locked, false: The BMS is not locked.

* `config_drive` -  This is a reserved attribute.

* `availability_zone` - Specifies the AZ ID.

* `description` -  Provides supplementary information about the pool.

* `kernel_id` - The UUID of the kernel image when the AMI image is used.

* `hypervisor_hostname` -  It is the name of a host on the hypervisor.

* `instance_name` - Instance name is specified.

* `tenant_id` - Specifies the ID of the tenant owning the BMS. The value is in UUID format.
    This parameter specifies the same meaning as project_id.

<a name="bms_security_groups"></a>
The `security_groups` block supports:

* `name` - The name of security_groups.
