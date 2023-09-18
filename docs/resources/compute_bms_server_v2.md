---
subcategory: "Bare Metal Server (BMS)"
description: ""
page_title: "flexibleengine_compute_bms_server_v2"
---

# flexibleengine_compute_bms_server_v2

Manages a BMS Server resource within FlexibleEngine.

## Example Usage

### Basic Instance

```hcl
variable "image_id" {}
variable "flavor_id" {}
variable "keypair_name" {}
variable "availability_zone" {}

resource "flexibleengine_compute_bms_server_v2" "basic" {
  name              = "basic"
  image_id          = var.image_id
  flavor_id         = var.flavor_id
  key_pair          = var.keypair_name
  security_groups   = ["default"]
  availability_zone = var.availability_zone

  metadata = {
    this = "that"
  }

  network {
    uuid = flexibleengine_vpc_subnet_v1.example_subnet.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the bms server instance. If
    omitted, the `region` argument of the provider is used. Changing this will create a new bms server.

* `name` - (Required, String) The name of the BMS.

* `image_id` - (Optional, String, ForceNew) Changing this creates a new bms server. It is Required if `image_name` is
  empty. Changing this creates a new bms server.

* `image_name` - (Optional, String, ForceNew) The name of the desired image for the bms server.
    Changing this creates a new bms server. It is Required if `image_id` is empty.

* `flavor_id` - (Optional, String) The flavor ID of
    the desired flavor for the bms server. Changing this resizes the existing bms server.
    It is Required if `flavor_name` is empty.

* `flavor_name` - (Optional, String) The name of the
    desired flavor for the bms server. Changing this resizes the existing bms server.
    It is Required if `flavor_id` is empty.

* `user_data` - (Optional, String, ForceNew) The user data to provide when launching the instance.
    Changing this creates a new bms server.

* `security_groups` - (Optional, List) An array of one or more security group names
    to associate with the bms server. Changing this results in adding/removing
    security groups from the existing bms server.

* `availability_zone` - (Required, String, ForceNew) The availability zone in which to create
    the bms server. Changing this will create a new bms server resource.

* `network` - (Optional, List, ForceNew) An array of one or more networks to attach to the
    bms instance. Changing this creates a new bms server.
  The [network](#<a name="bms_network"></a>) object structure is documented below.

* `metadata` - (Optional, Map) Metadata key/value pairs to make available from
    within the instance. Changing this updates the existing bms server metadata.

* `admin_pass` - (Optional, String) The administrative password to assign to the bms server.
    Changing this changes the root password on the existing server.

* `key_pair` - (Optional, String, ForceNew) The name of a key pair to put on the bms server. The key
    pair must already be created and associated with the tenant's account.
    Changing this creates a new bms server.

* `stop_before_destroy` - (Optional, Bool) Whether to try stop instance gracefully
    before destroying it, thus giving chance for guest OS daemons to stop correctly.
    If instance doesn't stop within timeout, it will be destroyed anyway.

<a name="bms_network"></a>
The `network` block supports:

* `uuid` - (Optional, String, ForceNew) The network UUID to attach to the bms server.
    Changing this creates a new bms server. It is Required unless `port`  or `name` is provided

* `name` - (Optional, String, ForceNew) The human-readable name of the network. Changing this creates a new bms server.
    It is Required unless `uuid` or `port` is provided.

* `port` - (Optional, String, ForceNew) The port UUID of a network to attach to the bms server.
    Changing this creates a new server.It is Required unless `uuid` or `name` is provided

* `fixed_ip_v4` - (Optional, String, ForceNew) Specifies a fixed IPv4 address to be used on this
    network. Changing this creates a new bms server.

* `fixed_ip_v6` - (Optional, String, ForceNew) Specifies a fixed IPv6 address to be used on this
    network. Changing this creates a new bms server.

* `access_network` - (Optional, Bool) Specifies if this network should be used for
    provisioning access. Accepts true or false. Defaults to false.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The id of the bms server.

* `config_drive` - Whether to use the config_drive feature to configure the instance.

* `kernel_id` - The UUID of the kernel image when the AMI image is used.

* `user_id` - The ID of the user to which the BMS belongs.

* `host_status` - The nova-compute status: **UP, UNKNOWN, DOWN, MAINTENANCE** and **Null**.

* `access_ip_v4` - This is a reserved attribute.

* `access_ip_v6` - This is a reserved attribute.

* `host_id` - Specifies the host ID of the BMS.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `update` - Default is 30 minutes.
* `delete` - Default is 30 minutes.
