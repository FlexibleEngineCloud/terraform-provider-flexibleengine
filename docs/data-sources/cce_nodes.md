---
subcategory: "Cloud Container Engine (CCE)"
---

# flexibleengine_cce_nodes

Use this data source to get a list of CCE nodes.

## Example Usage

```hcl
variable "cluster_id" {}
variable "node_name" {}

data "flexibleengine_cce_nodes" "node" {
  cluster_id = var.cluster_id
  name       = var.node_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the CCE nodes. If omitted, the provider-level
  region will be used.

* `cluster_id` - (Required, String) Specifies the ID of CCE cluster.

* `name` - (Optional, String) Specifies the name of the node.

* `node_id` - (Optional, String) Specifies the ID of the node.

* `status` - (Optional, String) Specifies the status of the node.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates a data source ID.

* `ids` - Indicates a list of IDs of all CCE nodes found.

* `nodes` - Indicates a list of CCE nodes found. The [nodes](#cce_nodes) object structure is documented below.

<a name="cce_nodes"></a>
The `nodes` block supports:

* `name` - The name of the node.

* `id` - The ID of the node.

* `status` - The status of the node.

* `flavor_id` - The flavor ID to be used.

* `availability_zone` - The available partitions where the node is located.

* `os` - The operating system of the node.

* `subnet_id` - The ID of the subnet to which the NIC belongs.

* `ecs_group_id` - The ID of ECS group to which the node belongs.

* `tags` - The tags of a VM node, key/value pair format.

* `key_pair` - The key pair name when logging in to select the key pair mode.

* `billing_mode` - The node's billing mode: The value is 0 (on demand).

* `server_id` - The node's virtual machine ID in ECS.

* `public_ip` - The elastic IP parameters of the node.

* `private_ip` - The private IP of the node.

* `root_volume` - The system disk related configuration. The [volume](#cce_volume) object structure is documented below.

* `data_volumes` - The data related configuration. The [volume](#cce_volume) object structure is documented below.

<a name="cce_volume"></a>
The `root_volume` and `data_volumes` blocks support:

* `size` - Disk size in GB.

* `volumetype` - Disk type.

* `extend_params` - Disk expansion parameters.
