---
subcategory: "Cloud Container Engine (CCE)"
description: ""
page_title: "flexibleengine_cce_node_pool_v3"
---

# flexibleengine_cce_node_pool_v3

Add a node pool to a container cluster.

## Example Usage

```hcl
variable "cluster_id" { }
variable "key_pair" { }
variable "availability_zone" { }

resource "flexibleengine_cce_node_pool_v3" "node_pool" {
  cluster_id               = var.cluster_id
  name                     = "testpool"
  os                       = "EulerOS 2.5"
  initial_node_count       = 2
  flavor_id                = "s3.large.4"
  availability_zone        = var.availability_zone
  key_pair                 = var.keypair
  scale_enable             = true
  min_node_count           = 1
  max_node_count           = 10
  scale_down_cooldown_time = 100
  priority                 = 1
  type                     = "vm"

  root_volume {
    size       = 40
    volumetype = "SAS"
  }
  data_volumes {
    size       = 100
    volumetype = "SAS"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the CCE node pool resource.
  If omitted, the provider-level region will be used. Changing this will create a new CCE node pool resource.

* `cluster_id` - (Required, String, ForceNew) ID of the cluster. Changing this parameter will create a new resource.

* `name` - (Required, String) Node Pool Name.

* `initial_node_count` - (Required, Int) Specifies the initial number of expected nodes in the node pool.
  This parameter can be also used to manually scale the node count afterwards.

* `flavor_id` - (Required, String, ForceNew) Specifies the flavor id. Changing this parameter will create a new resource.

* `type` - (Optional, String, ForceNew) Node Pool type. Possible values are: "vm" and "ElasticBMS".

* `availability_zone` - (Optional, String, ForceNew) specify the name of the available partition (AZ).
    Default value is random to create nodes in a random AZ in the node pool.
    Changing this parameter will create a new resource.

* `os` - (Optional, String, ForceNew) Operating System of the node. The value can be EulerOS 2.5 and CentOS 7.6.
    Changing this parameter will create a new resource.

* `runtime` - (Optional, String, ForceNew) Specifies the runtime of the CCE node pool. Valid values are *docker* and
  *containerd*. Changing this creates a new resource.

* `key_pair` - (Optional, String, ForceNew) Key pair name when logging in to select the key pair mode.
    This parameter and `password` are alternative. Changing this parameter will create a new resource.

* `password` - (Optional, String, ForceNew) root password when logging in to select the password mode.
    This parameter must be **salted** and alternative to `key_pair`. Changing this parameter will create a new resource.

* `subnet_id` - (Optional, String, ForceNew) The ID of the VPC Subnet to which the NIC belongs.
    Changing this parameter will create a new resource.

* `security_groups` - (Optional, List, ForceNew) Specifies the list of custom security group IDs for the node pool.
  If specified, the nodes will be put in these security groups. When specifying a security group, do not modify
  the rules of the port on which CCE running depends. Changing this parameter will create a new resource.

* `ecs_group_id` - (Optional, String, ForceNew) Specifies the ECS group ID. If specified, the node will be created under
  the cloud server group. Changing this parameter will create a new resource.

* `max_pods` - (Optional, Int, ForceNew) The maximum number of instances a node is allowed to create.
    Changing this parameter will create a new resource.

* `preinstall` - (Optional, String, ForceNew) Script required before installation. The input value can be
    a Base64 encoded string or not. Changing this parameter will create a new resource.

* `postinstall` - (Optional, String, ForceNew) Script required after the installation. The input value can be
    a Base64 encoded string or not. Changing this parameter will create a new resource.

* `scale_enable` - (Optional, Bool) Whether to enable auto scaling. If Autoscaler is enabled, install the autoscaler
    add-on to use the auto scaling feature.

* `min_node_count` - (Optional, Int) Minimum number of nodes allowed if auto scaling is enabled.

* `max_node_count` - (Optional, Int) Maximum number of nodes allowed if auto scaling is enabled.

* `scale_down_cooldown_time` - (Optional, Int) Interval between two scaling operations, in minutes.

* `priority` - (Optional, Int) Weight of a node pool. A node pool with a higher weight has a higher priority during scaling.

* `labels` - (Optional, Map) Tags of a Kubernetes node, key/value pair format.

* `tags` - (Optional, Map) Tags of a VM node, key/value pair format.

* `root_volume` - (Required, List, ForceNew) It corresponds to the system disk related configuration.
    The object structure is documented below. Changing this parameter will create a new resource.

* `data_volumes` - (Required, List, ForceNew) Represents the data disk to be created.
    The object structure is documented below. Changing this parameter will create a new resource.

* `taints` - (Optional, List) You can add taints to created nodes to configure anti-affinity.
    The object structure is documented below.

The `root_volume` block supports:

* `size` - (Required, Int, ForceNew) Specifies the disk size in GB.
  Changing this will create a new CCE node pool resource.

* `volumetype` - (Required, String, ForceNew) Specifies the disk type.
  Changing this will create a new CCE node pool resource.

* `kms_key_id` - (Optional, String, ForceNew) Specifies the KMS key ID. This is used to encrypt the volume.
  Changing this will create a new CCE node pool resource.

  -> You need to create an agency (EVSAccessKMS) when disk encryption is used in the current project for the first time ever.
  The account and permission of the created agency are `op_svc_evs` and **KMS Administrator**, respectively.

* `extend_params` - (Optional, Map, ForceNew) Specifies the disk expansion parameters in key/value pair format.
  Changing this will create a new CCE node pool resource.

The `data_volumes` block supports:

* `size` - (Required, Int, ForceNew) Specifies the disk size in GB.
  Changing this will create a new CCE node pool resource.

* `volumetype` - (Required, String, ForceNew) Specifies the disk type.
  Changing this will create a new CCE node pool resource.

* `kms_key_id` - (Optional, String, ForceNew) Specifies the KMS key ID. This is used to encrypt the volume.
  Changing this will create a new CCE node pool resource.

  -> You need to create an agency (EVSAccessKMS) when disk encryption is used in the current project for the first time ever.
  The account and permission of the created agency are `op_svc_evs` and **KMS Administrator**, respectively.

* `extend_params` - (Optional, Map, ForceNew) Specifies the disk expansion parameters in key/value pair format.
  Changing this will create a new CCE node pool resource.

The `taints` block supports:

* `key` - (Required, String) A key must contain 1 to 63 characters starting with a letter or digit.
  Only letters, digits, hyphens (-), underscores (_), and periods (.) are allowed.
  A DNS subdomain name can be used as the prefix of a key.

* `value` - (Required, String) A value must start with a letter or digit and can contain a maximum of 63 characters,
  including letters, digits, hyphens (-), underscores (_), and periods (.).

* `effect` - (Required, String) Available options are *NoSchedule*, *PreferNoSchedule* and *NoExecute*.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.

* `status` -  Node status information.

* `current_node_count` - The current number of the nodes.

* `billing_mode` -  Billing mode of a node.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.
* `delete` - Default is 20 minutes.

## Import

Node_pool can be imported using the cluster ID and node_pool ID, e.g.

```shell
terraform import flexibleengine_cce_node_pool_v3.node_pool_1 <cluster_id>/<id>
```
