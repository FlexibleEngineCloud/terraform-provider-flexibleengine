---
subcategory: "Cloud Container Engine (CCE)"
---

# flexibleengine_cce_node_v3

Add a node to a container cluster.

## Example Usage

```hcl
variable "cluster_id" { }
variable "ssh_key" { }
variable "availability_zone" { }

resource "flexibleengine_cce_node_v3" "node_1" {
  cluster_id        = var.cluster_id
  name              = "node1"
  flavor_id         = "s1.medium"
  availability_zone = var.availability_zone
  key_pair          = var.ssh_key
  iptype            = "5_bgp"
  sharetype         = "PER"
  bandwidth_size    = 100

  root_volume {
    size       = 40
    volumetype = "SATA"
  }
  data_volumes {
    size       = 100
    volumetype = "SATA"
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) ID of the cluster. Changing this parameter will create a new resource.

* `name` - (Optional) Node Name.

* `flavor_id` - (Required) Specifies the flavor id. Changing this parameter will create a new resource.

* `availability_zone` - (Required) specify the name of the available partition (AZ).
  Changing this parameter will create a new resource.

* `key_pair` - (Required) Key pair name when logging in to select the key pair mode.
  Changing this parameter will create a new resource.

* `os` - (Optional) Operating System of the node, possible values are EulerOS 2.2 and CentOS 7.6. Defaults to EulerOS 2.2.
    Changing this parameter will create a new resource.

* `labels` - (Optional) Tags of a Kubernetes node, key/value pair format. Changing this parameter will create a new resource.

* `tags` - (Optional) VM tag, key/value pair format.

* `annotations` - (Optional) Node annotation, key/value pair format. Changing this parameter will create a new resource.

* `eip_ids` - (Optional) List of existing elastic IP IDs. Changing this parameter will create a new resource.

  -> If the `eip_ids` parameter is configured, you do not need to configure the `eip_count` and bandwidth parameters:
  `iptype`, `bandwidth_charge_mode`, `bandwidth_size` and `share_type`.

* `eip_count` - (Optional) Number of elastic IPs to be dynamically created. Changing this parameter will create a new resource.

* `iptype` - (Optional) Elastic IP type.

* `bandwidth_charge_mode` - (Optional) Bandwidth billing type. Changing this parameter will create a new resource.

* `sharetype` - (Optional) Bandwidth sharing type. Changing this parameter will create a new resource.

* `bandwidth_size` - (Optional) Bandwidth size. Changing this parameter will create a new resource.

* `ecs_performance_type` - (Optional) Classification of cloud server specifications.
    Changing this parameter will create a new resource.

* `product_id` - (Optional) The Product ID. Changing this parameter will create a new resource.

* `max_pods` - (Optional) The maximum number of instances a node is allowed to create.
    Changing this parameter will create a new resource.

* `public_key` - (Optional) The Public key. Changing this parameter will create a new resource.

* `preinstall` - (Optional) Script required before installation. The input value can be a Base64 encoded string or not.
    Changing this parameter will create a new resource.

* `postinstall` - (Optional) Script required after installation. The input value can be a Base64 encoded string or not.
   Changing this parameter will create a new resource.

* `extend_param` - (Optional, Map, ForceNew) Extended parameter. Changing this parameter will create a new resource.
  Availiable keys:

  + `agency_name` - Specifies the agency name to provide temporary credentials for CCE node to access other cloud services.
  + `dockerBaseSize` - The available disk space of a single docker container on the node in device mapper mode.
  + `DockerLVMConfigOverride` - Docker data disk configurations. The following is an example default configuration:

    ```hcl
    extend_param = {
      DockerLVMConfigOverride = "dockerThinpool=vgpaas/90%VG;kubernetesLV=vgpaas/10%VG;diskType=evs;lvType=linear"
    }
    ```

**root_volume** **- (Required)** It corresponds to the system disk related configuration.
  Changing this parameter will create a new resource.

  * `size` - (Required) Specifies the disk size in GB.
  * `volumetype` - (Required) Specifies the disk type.
  * `extend_params` - (Optional) Specifies the disk expansion parameters in key/value pair format.

**data_volumes** **- (Required)** Represents the data disk to be created.
  Changing this parameter will create a new resource.

  * `size` - (Required) Specifies the disk size in GB.
  * `volumetype` - (Required) Specifies the disk type.
  * `extend_params` - (Optional) Specifies the disk expansion parameters in key/value pair format.
  * `kms_key_id` - (Optional) Specifies the ID of a KMS key. This is used to encrypt the volume.

  -> You need to create an agency (EVSAccessKMS) when disk encryption is used in the current project for the first time ever.
  The account and permission of the created agency are `op_svc_evs` and **KMS Administrator**, respectively.

**taints** **- (Optional)** You can add taints to created nodes to configure anti-affinity.
  Changing this parameter will create a new resource.
  Each taint contains the following parameters:

  * `key` - (Required) A key must contain 1 to 63 characters starting with a letter or digit. Only letters, digits,
    hyphens (-), underscores (_), and periods (.) are allowed. A DNS subdomain name can be used as the prefix of a key.
  * `value` - (Required) A value must start with a letter or digit and can contain a maximum of 63 characters,
    including letters, digits, hyphens (-), underscores (_), and periods (.).
  * `effect` - (Required) Available options are NoSchedule, PreferNoSchedule, and NoExecute.

## Attributes Reference

All above argument parameters can be exported as attribute parameters along with attribute reference.

* `id` - The resource ID in UUID format.
* `status` -  Node status information.
* `server_id` - ID of the ECS instance associated with the node.
* `private_ip` - Private IP of the CCE node.
* `public_ip` - Public IP of the CCE node.
