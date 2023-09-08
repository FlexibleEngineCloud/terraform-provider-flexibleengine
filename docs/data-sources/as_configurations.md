---
subcategory: "Auto Scaling (AS)"
---

# flexibleengine_as_configurations

Use this data source to get a list of AS configurations.

```hcl
data "flexibleengine_as_configurations" "configurations" {
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region in which to obtain the AS configurations.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the AS configuration name. Supports fuzzy search.

* `image_id` - (Optional, String) Specifies the image ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the list.

* `configurations` - A list of AS configurations.
  The [configurations](#instance_configurations) object structure is documented below.

<a name="instance_configurations"></a>
The `configurations` block supports:

* `scaling_configuration_name` - The AS configuration name.

* `instance_config` - The list of information about instance configurations.
  The [instance_config](#instance_config_object) structure is documented below.

* `status` - The AS configuration status, the value can be **Bound** or **Unbound**.

<a name="instance_config_object"></a>
The `instance_config` block supports:

* `instance_id` - The ECS instance ID when using its specification as the template to create AS configurations.

* `flavor` - The ECS flavor name.

* `image` - The ECS image ID.

* `disk` - The list of disk group information. The [disk](#instance_config_disk_object) structure is documented below.

* `key_name` - The name of the SSH key pair used to log in to the instance.

* `security_group_ids` - An array of one or more security group IDs.

* `charging_mode` - The billing mode for ECS, the value can be **postPaid** or **spot**.

* `flavor_priority_policy` - The priority policy used when there are multiple flavors
  and instances to be created using an AS configuration. The value can be `PICK_FIRST` and `COST_FIRST`.

* `ecs_group_id` - The ECS group ID.

* `user_data` - The user data to provide when launching the instance.

* `public_ip` - The EIP list of the ECS instance.
  The [public_ip](#instance_config_public_ip_object) structure is documented below.

* `metadata` - The key/value pairs to make available from within the instance.

* `personality` - The list of information about the injected file.
  The [personality](#instance_config_personality_object) structure is documented below.

<a name="instance_config_disk_object"></a>
The `disk` block supports:

* `size` - The disk size. The unit is GB.

* `volume_type` - The volume type.

* `disk_type` - The disk type.

* `kms_id` - The encryption KMS ID of the **DATA** disk.

<a name="instance_config_public_ip_object"></a>
The `public_ip` block supports:

* `eip` - The list of EIP configuration that will be automatically assigned to the instance.
  The [eip](#instance_eip) object structure is documented below.

<a name="instance_eip"></a>
The `eip` block supports:

* `ip_type` - The EIP type.

* `bandwidth` - The list of bandwidth information.
  The [bandwidth](#instance_bandwidth) object structure is documented below.

<a name="instance_bandwidth"></a>
The `bandwidth` block supports:

* `share_type` - The bandwidth sharing type.

* `charging_mode` - The bandwidth billing mode, the value can be **traffic** or **bandwidth**.

* `size` - The bandwidth (Mbit/s).

<a name="instance_config_personality_object"></a>
The `personality` block supports:

* `path` - The path of the injected file.

* `content` - The content of the injected file.
