---
subcategory: "Cloud Container Engine (CCE)"
---

# flexibleengine_cce_node_v3

To get the specified CCE node in a cluster.

## Example Usage

```hcl
variable "cluster_id" {}
variable "node_name" {}
  
data "flexibleengine_cce_node_v3" "node" {
  cluster_id = var.cluster_id
  name       = var.node_name
}
```

## Argument Reference

The following arguments are supported:
 
* `cluster_id` - (Required) The id of container cluster.

* `name` - (Optional) - Name of the node.

* `node_id` - (Optional) - The id of the node.

* `status` - (Optional) - The state of the node.

## Attributes Reference

All above argument parameters can be exported as attribute parameters along with attribute reference:

* `flavor_id` - The flavor id to be used. 

* `availability_zone` - Available partitions where the node is located. 

* `key_pair` - Key pair name when logging in to select the key pair mode.

* `billing_mode` - Node's billing mode: The value is 0 (on demand).

* `eip_ids` - List of existing elastic IP IDs.
 
* `server_id` - The node's virtual machine ID in ECS.

* `private_ip` - Private IP of the node

* `public_ip` - Elastic IP parameters of the node.

* `ip_type` - Elastic IP address type.

* `share_type` - Bandwidth sharing type.
* `bandwidth_size` - Bandwidth (Mbit/s), in the range of [1, 2000].
* `charge_mode` - Bandwidth billing type.

**root_volumes**

  * `disk_size` - Disk size in GB.
  * `volume_type` - Disk type.

**data_volumes**

  * `disk_size` - Disk size in GB.
  * `volume_type` - Disk type.
