---
layout: "flexibleengine"
page_title: "FlexibleEngine: flexibleengine_compute_bms_server_v2"
sidebar_current: "docs-flexibleengine-resource-compute-bms-server-v2"
description: |-
  Manages a BMS server resource within FlexibleEngine.
---

# flexibleengine_compute_bms_server_v2

Manages a BMS Server resource within FlexibleEngine.

## Example Usage

### Basic Instance

```hcl
variable "image_id" {}
variable "flavor_id" {}
variable "keypair_name" {}
variable "network_id" {}
variable "availability_zone" {}

resource "flexibleengine_compute_bms_server_v2" "basic" {
  name            = "basic"
  image_id        = "${var.image_id}"
  flavor_id       = "${var.flavor_id}"
  key_pair        = "${var.keypair_name}"
  security_groups = ["default"]
  availability_zone = "${var.availability_zone}"

  metadata = {
    this = "that"
  }

  network {
    uuid = "${var.network_id}"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to create the bms server instance. If
    omitted, the `region` argument of the provider is used. Changing this
    creates a new bms server.

* `name` - (Required) The name of the BMS.

* `image_id` - (Optional; Required if `image_name` is empty.) Changing this creates a new bms server.

* `image_name` - (Optional; Required if `image_id` is empty.) The name of the
    desired image for the bms server. Changing this creates a new bms server.

* `flavor_id` - (Optional; Required if `flavor_name` is empty) The flavor ID of
    the desired flavor for the bms server. Changing this resizes the existing bms server.

* `flavor_name` - (Optional; Required if `flavor_id` is empty) The name of the
    desired flavor for the bms server. Changing this resizes the existing bms server.

* `user_data` - (Optional) The user data to provide when launching the instance.
    Changing this creates a new bms server.

* `security_groups` - (Optional) An array of one or more security group names
    to associate with the bms server. Changing this results in adding/removing
    security groups from the existing bms server.

* `availability_zone` - (Required) The availability zone in which to create
    the bms server.

* `network` - (Optional) An array of one or more networks to attach to the
    bms instance. Changing this creates a new bms server.

* `metadata` - (Optional) Metadata key/value pairs to make available from
    within the instance. Changing this updates the existing bms server metadata.

* `admin_pass` - (Optional) The administrative password to assign to the bms server.
    Changing this changes the root password on the existing server.

* `key_pair` - (Optional) The name of a key pair to put on the bms server. The key
    pair must already be created and associated with the tenant's account.
    Changing this creates a new bms server.

* `stop_before_destroy` - (Optional) Whether to try stop instance gracefully
    before destroying it, thus giving chance for guest OS daemons to stop correctly.
    If instance doesn't stop within timeout, it will be destroyed anyway.

The `network` block supports:

* `uuid` - (Required unless `port`  or `name` is provided) The network UUID to
    attach to the bms server. Changing this creates a new bms server.

* `name` - (Required unless `uuid` or `port` is provided) The human-readable
    name of the network. Changing this creates a new bms server.

* `port` - (Required unless `uuid` or `name` is provided) The port UUID of a
    network to attach to the bms server. Changing this creates a new server.

* `fixed_ip_v4` - (Optional) Specifies a fixed IPv4 address to be used on this
    network. Changing this creates a new bms server.

* `fixed_ip_v6` - (Optional) Specifies a fixed IPv6 address to be used on this
    network. Changing this creates a new bms server.

* `access_network` - (Optional) Specifies if this network should be used for
    provisioning access. Accepts true or false. Defaults to false.

# Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The id of the bms server.

* `config_drive` - Whether to use the config_drive feature to configure the instance.

* `kernel_id` - The UUID of the kernel image when the AMI image is used.

* `user_id` - The ID of the user to which the BMS belongs.

* `host_status` - The nova-compute status: **UP, UNKNOWN, DOWN, MAINTENANCE** and **Null**.
