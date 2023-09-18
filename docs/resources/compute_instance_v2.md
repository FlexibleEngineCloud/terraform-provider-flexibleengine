---
subcategory: "Elastic Cloud Server (ECS)"
description: ""
page_title: "flexibleengine_compute_instance_v2"
---

# flexibleengine_compute_instance_v2

Manages a V2 VM instance resource within FlexibleEngine.

## Example Usage

### Basic Instance

```hcl
resource "flexibleengine_compute_instance_v2" "basic" {
  name            = "basic"
  image_id        = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id       = "s3.large.2"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  network {
    uuid = flexibleengine_vpc_subnet_v1.example_subnet.id
  }

  tags = {
    foo  = "bar"
    this = "that"
  }
}
```

### Instance With Attached Volume

```hcl
resource "flexibleengine_blockstorage_volume_v2" "myvol" {
  name = "myvol"
  size = 1
}

resource "flexibleengine_compute_instance_v2" "myinstance" {
  name            = "myinstance"
  image_id        = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id       = "s3.large.2"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  network {
    uuid = flexibleengine_vpc_subnet_v1.example_subnet.id
  }
}

resource "flexibleengine_compute_volume_attach_v2" "attached" {
  instance_id = flexibleengine_compute_instance_v2.myinstance.id
  volume_id   = flexibleengine_blockstorage_volume_v2.myvol.id
}
```

### Boot From Volume

```hcl
resource "flexibleengine_compute_instance_v2" "boot-from-volume" {
  name            = "boot-from-volume"
  flavor_id       = "s3.large.2"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  block_device {
    uuid                  = "<image-id>"
    source_type           = "image"
    volume_size           = 5
    boot_index            = 0
    destination_type      = "volume"
    delete_on_termination = true
    volume_type           = "SSD"
  }

  network {
    uuid = flexibleengine_vpc_subnet_v1.example_subnet.id
  }
}
```

### Boot From an Existing Volume

```hcl
resource "flexibleengine_blockstorage_volume_v2" "myvol" {
  name     = "myvol"
  size     = 5
  image_id = "<image-id>"
}

resource "flexibleengine_compute_instance_v2" "boot-from-volume" {
  name            = "bootfromvolume"
  flavor_id       = "s3.large.2"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  block_device {
    uuid                  = flexibleengine_blockstorage_volume_v2.myvol.id
    source_type           = "volume"
    boot_index            = 0
    destination_type      = "volume"
    delete_on_termination = true
  }

  network {
    uuid = flexibleengine_vpc_subnet_v1.example_subnet.id
  }
}
```

### Boot Instance, Create Volume, and Attach Volume as a Block Device

```hcl
resource "flexibleengine_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  image_id        = "<image-id>"
  flavor_id       = "s3.large.2"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  block_device {
    uuid                  = "<image-id>"
    source_type           = "image"
    destination_type      = "local"
    boot_index            = 0
    delete_on_termination = true
  }

  block_device {
    source_type           = "blank"
    destination_type      = "volume"
    volume_size           = 1
    boot_index            = 1
    delete_on_termination = true
  }
}
```

### Boot Instance and Attach Existing Volume as a Block Device

```hcl
resource "flexibleengine_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  size = 1
}

resource "flexibleengine_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  image_id        = "<image-id>"
  flavor_id       = "s3.large.2"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  block_device {
    uuid                  = "<image-id>"
    source_type           = "image"
    destination_type      = "local"
    boot_index            = 0
    delete_on_termination = true
  }

  block_device {
    uuid                  = flexibleengine_blockstorage_volume_v2.volume_1.id
    source_type           = "volume"
    destination_type      = "volume"
    boot_index            = 1
    delete_on_termination = true
  }
}
```

### Instance With Multiple Networks

```hcl
resource "flexibleengine_vpc_eip" "eip_1" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name       = "test"
    size       = 10
    share_type = "PER"
  }
}

resource "flexibleengine_compute_instance_v2" "multi-net" {
  name            = "multi-net"
  image_id        = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id       = "s3.large.2"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  network {
    uuid = flexibleengine_vpc_subnet_v1.example_subnet.id
  }

  network {
    uuid = flexibleengine_vpc_subnet_v1.example_subnet_2.id
  }
}

resource "flexibleengine_compute_floatingip_associate_v2" "myip" {
  floating_ip = flexibleengine_vpc_eip.eip_1.publicip.0.ip_address
  instance_id = flexibleengine_compute_instance_v2.multi-net.id
  fixed_ip    = flexibleengine_compute_instance_v2.multi-net.network.1.fixed_ip_v4
}
```

### Instance with Multiple Ephemeral Disks

```hcl
resource "flexibleengine_compute_instance_v2" "multi-eph" {
  name            = "multi_eph"
  image_id        = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id       = "s3.large.2"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  block_device {
    boot_index            = 0
    delete_on_termination = true
    destination_type      = "local"
    source_type           = "image"
    uuid                  = "<image-id>"
  }

  block_device {
    boot_index            = -1
    delete_on_termination = true
    destination_type      = "local"
    source_type           = "blank"
    volume_size           = 1
  }

  block_device {
    boot_index            = -1
    delete_on_termination = true
    destination_type      = "local"
    source_type           = "blank"
    volume_size           = 1
  }
}
```

### Instance with User Data (cloud-init)

```hcl
resource "flexibleengine_compute_instance_v2" "instance_1" {
  name            = "basic"
  image_id        = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id       = "s3.large.2"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]
  user_data       = "#cloud-config\nhostname: instance_1.example.com\nfqdn: instance_1.example.com"

  network {
    uuid = flexibleengine_vpc_subnet_v1.example_subnet.id
  }
}
```

`user_data` can come from a variety of sources: inline, read in from the `file`
function, or the `template_cloudinit_config` resource.

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the server instance.
  If omitted, the `region` argument of the provider is used. Changing this creates a new server.

* `name` - (Required, String) A unique name for the resource.

* `image_id` - (Optional, String, ForceNew) The image ID of the desired image for the server.
  It is **Required** if `image_name` is empty and not booting from a volume. Do not specify if booting from a volume.
  Changing this creates a new server.

* `image_name` - (Optional, String, ForceNew) The name of the desired image for the server.
  It is **Required** if `image_id` is empty and not booting from a volume. Do not specify if booting from a volume.
  Changing this creates a new server.

* `flavor_id` - (Optional, String) The flavor ID of the desired flavor for the server.
  It is **Required** if `flavor_name` is empty. Changing this resizes the existing server.

* `flavor_name` - (Optional, String) The name of the desired flavor for the server.
  It is **Required** if `flavor_id` is empty. Changing this resizes the existing server.

* `availability_zone` - (Optional, String, ForceNew) The availability zone in which to create the server.
  Changing this creates a new server.

* `security_groups` - (Optional, List) An array of one or more security group names
  to associate with the server. Changing this results in adding/removing
  security groups from the existing server. *Note*: When attaching the
  instance to networks using Ports, place the security groups on the Port and not the instance.

* `network` - (Optional, List, ForceNew) An array of one or more networks to attach to the
  instance. The [network](#ecs_arg_network) object structure is documented below. Changing this
  creates a new server.

* `user_data` - (Optional, String, ForceNew) The user data to provide when launching the instance.
  Changing this creates a new server.

* `metadata` - (Optional, Map) The key/value pairs to associate with the instance.

* `config_drive` - (Optional, Bool, ForceNew) Whether to use the config_drive feature to
  configure the instance. Changing this creates a new server.

* `admin_pass` - (Optional, String) The administrative password to assign to the server.

* `key_pair` - (Optional, String, ForceNew) The name of a key pair to put on the server. The key
  pair must already be created and associated with the tenant's account. Changing this creates a new server.

* `block_device` - (Optional, List) Configuration of block devices. The [block_device](#ecs_arg_block_device) object
  structure is documented below. You can specify multiple block devices which will create an instance with
  multiple disks. This configuration is very flexible, so please see the
  following [reference](http://docs.openstack.org/developer/nova/block_device_mapping.html) for more information.

* `scheduler_hints` - (Optional, List) Provide the Nova scheduler with hints on how
  the instance should be launched. The [scheduler_hints](#ecs_scheduler_hints) object structure is documented below.

* `stop_before_destroy` - (Optional, Bool) Whether to try stop instance gracefully
  before destroying it, thus giving chance for guest OS daemons to stop correctly.
  If instance doesn't stop within timeout, it will be destroyed anyway.

* `auto_recovery` - (Optional, Bool) Configures or deletes automatic recovery of an instance

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the instance.

<a name="ecs_arg_network"></a>
The `network` block supports:

* `uuid` - (Optional, String, ForceNew) The network UUID to attach to the server. It is **Required** unless `port` is
  provided. Changing this creates a new server.

* `port` - (Optional, String, ForceNew) The port UUID of a network to attach to the server.
  It is **Required** unless `uuid` is provided. Changing this creates a new server.

* `fixed_ip_v4` - (Optional, String, ForceNew) Specifies a fixed IPv4 address to be used on this
  network. Changing this creates a new server.

* `fixed_ip_v6` - (Optional, String, ForceNew) Specifies a fixed IPv6 address to be used on this
  network. Changing this creates a new server.

* `access_network` - (Optional, Bool) Specifies if this network should be used for
  provisioning access. Accepts true or false. Defaults to false.

<a name="ecs_arg_block_device"></a>
The `block_device` block supports:

* `uuid` - (Optional, String, ForceNew) The UUID of the image, volume, or snapshot.
  It is Required unless `source_type` is set to `"blank"`. Changing this creates a new server.

* `source_type` - (Required, String, ForceNew) The source type of the device. Must be one of
  "blank", "image", "volume", or "snapshot". Changing this creates a new server.

* `volume_size` - (Optional, Int, ForceNew)The size of the volume to create (in gigabytes). Required
  in the following combinations: source=image and destination=volume, source=blank and destination=local,
  and source=blank and destination=volume. Changing this creates a new server.

* `volume_type` - (Optional, String, ForceNew) Currently, the value can be `SSD` (ultra-I/O disk type),
  `SAS` (high I/O disk type), or `SATA` (common I/O disk type). Changing this creates a new server.

* `boot_index` - (Optional, Int, ForceNew) The boot index of the volume. It defaults to 0, which
  indicates that it's a system disk. Changing this creates a new server.

* `destination_type` - (Optional, String, ForceNew) The type that gets created. Possible values
  are "volume" and "local". Changing this creates a new server.

* `delete_on_termination` - (Optional, Bool, ForceNew) Delete the volume / block device upon
  termination of the instance. Defaults to false. Changing this creates a new server.

* `disk_bus` - (Optional, String) The low-level disk bus that will be used, for example, *virtio*, *scsi*.
  Most common thing is to leave this empty. Changing this creates a new server.

<a name="ecs_scheduler_hints"></a>
The `scheduler_hints` block supports:

* `group` - (Optional, String, ForceNew) Specifies the **anti-affinity** group ID.
  The instance will be placed into that group. Changing this creates a new server.

* `tenancy` - (Optional, String, ForceNew) Specifies whether the ECS is created on a Dedicated Host (DeH) or in a
  shared pool (default). The value can be **shared** or **dedicated**.

* `deh_id` - (Optional, String, ForceNew) Specifies the DeH ID.This parameter takes effect only when the value of tenancy
  is dedicated. If you do not specify this parameter, the system will automatically assign a DeH to you to deploy ECSs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `access_ip_v4` - The first detected Fixed IPv4 address *or* the Floating IP.

* `access_ip_v6` - The first detected Fixed IPv6 address.

* `network` - See Argument Reference above. The [network](#ecs_attr_network) object structure is documented below.

* `all_metadata` - Contains all instance metadata, even metadata not set by Terraform.

* `floating_ip` - The EIP address that is associate to the instance.

* `system_disk_id` - The system disk volume ID.

* `volume_attached` - An array of one or more disks to attach to the instance.
   The [volume_attached](#ecs_attr_volume_attached) object structure is documented below.

* `status` - The status of the instance.

<a name="ecs_attr_network"></a>
The `network` block supports:

* `mac` - The MAC address of the NIC on that network.

<a name="ecs_attr_volume_attached"></a>
The `volume_attached` block supports:

* `uuid` - The volume id on that attachment.

* `boot_index` - The volume boot index on that attachment.

* `size` - The volume size on that attachment.

* `type` - The volume type on that attachment.

* `pci_address` - The volume pci address on that attachment.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `update` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

Instances can be imported by their `id`. For example,

```shell
terraform import flexibleengine_compute_instance_v2.my_instance b11b407c-e604-4e8d-8bc4-92398320b847
```

Note that the imported state may not be identical to your resource definition, due to some attrubutes
missing from the API response, security or some other reason. The missing attributes include:
`admin_pass`, `config_drive`, `user_data`, `block_device`, `scheduler_hints`, `stop_before_destroy`,
`network/access_network` and arguments for pre-paid. It is generally recommended running
`terraform plan` after importing an instance. You can then decide if changes should
be applied to the instance, or the resource definition should be updated to align
with the instance. Also you can ignore changes as below.

```hcl
resource "flexibleengine_compute_instance_v2" "my_instance" {
    ...

  lifecycle {
    ignore_changes = [
      user_data, block_device,
    ]
  }
}
```

## Notes

### Multiple Ephemeral Disks

It's possible to specify multiple `block_device` entries to create an instance
with multiple ephemeral (local) disks. In order to create multiple ephemeral
disks, the sum of the total amount of ephemeral space must be less than or
equal to what the chosen flavor supports.

The following example shows how to create an instance with multiple ephemeral
disks:

```hcl
resource "flexibleengine_compute_instance_v2" "foo" {
  name            = "terraform-test"
  security_groups = ["default"]

  block_device {
    boot_index            = 0
    delete_on_termination = true
    destination_type      = "local"
    source_type           = "image"
    uuid                  = "<image uuid>"
  }

  block_device {
    boot_index            = -1
    delete_on_termination = true
    destination_type      = "local"
    source_type           = "blank"
    volume_size           = 1
  }

  block_device {
    boot_index            = -1
    delete_on_termination = true
    destination_type      = "local"
    source_type           = "blank"
    volume_size           = 1
  }
}
```

### Instances and Ports

Neutron Ports are a great feature and provide a lot of functionality. However,
there are some notes to be aware of when mixing Instances and Ports:

* When attaching an Instance to one or more networks using Ports, place the
security groups on the Port and not the Instance. If you place the security
groups on the Instance, the security groups will not be applied upon creation,
but they will be applied upon a refresh. This is a known FlexibleEngine bug.

* Network IP information is not available within an instance for networks that
are attached with Ports. This is mostly due to the flexibility Neutron Ports
provide when it comes to IP addresses. For example, a Neutron Port can have
multiple Fixed IP addresses associated with it. It's not possible to know which
single IP address the user would want returned to the Instance's state
information. Therefore, in order for a Provisioner to connect to an Instance
via it's network Port, customize the `connection` information:

```hcl
resource "flexibleengine_networking_port_v2" "port_1" {
  name           = "port_1"
  admin_state_up = "true"
  network_id     = flexibleengine_vpc_subnet_v1.example_subnet.id

  security_group_ids = [
    "2f02d20a-8dca-49b7-b26f-b6ce9fddaf4f",
    "ca1e5ed7-dae8-4605-987b-fadaeeb30461",
  ]
}

resource "flexibleengine_compute_instance_v2" "instance_1" {
  name = "instance_1"

  network {
    port = flexibleengine_networking_port_v2.port_1.id
  }

  connection {
    user        = "root"
    host        = flexibleengine_networking_port_v2.port_1.fixed_ip.0.ip_address
    private_key = "~/path/to/key"
  }

  provisioner "remote-exec" {
    inline = [
      "echo terraform executed > /tmp/foo",
    ]
  }
}
```
