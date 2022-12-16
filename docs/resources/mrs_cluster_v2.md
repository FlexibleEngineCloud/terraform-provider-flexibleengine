---
subcategory: "MapReduce Service (MRS)"
description: ""
page_title: "flexibleengine_mrs_cluster_v2"
---

# flexibleengine_mrs_cluster_v2

Manages a MRS cluster resource within FlexibleEngine.

## Example Usage

### Create an analysis cluster

```hcl
variable "mrs_az" {}
variable "cluster_name" {}
variable "password" {}
variable "keypair" {}
variable "vpc_id" {}
variable "subnet_id" {}

resource "flexibleengine_mrs_cluster_v2" "test" {
  availability_zone  = var.mrs_az
  name               = var.cluster_name
  version            = "MRS 2.0.1"
  type               = "ANALYSIS"
  component_list     = ["Hadoop", "Spark", "Hive", "Tez"]
  manager_admin_pwd  = var.password
  node_key_pair      = var.keypair
  vpc_id             = var.vpc_id
  subnet_id          = var.subnet_id

  master_nodes {
    flavor            = "c6.4xlarge.4.linux.mrs"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
  analysis_core_nodes {
    flavor            = "c6.4xlarge.4.linux.mrs"
    node_number       = 3
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

### Create a stream cluster

```hcl
variable "mrs_az" {}
variable "cluster_name" {}
variable "password" {}
variable "keypair" {}
variable "vpc_id" {}
variable "subnet_id" {}

resource "flexibleengine_mrs_cluster_v2" "test" {
  availability_zone  = var.mrs_az
  name               = var.cluster_name
  type               = "STREAMING"
  version            = "MRS 3.1.0-LTS.1"
  manager_admin_pwd  = var.password
  node_key_pair      = var.keypair
  vpc_id             = var.vpc_id
  subnet_id          = var.subnet_id
  component_list     = ["Ranger", "Kafka", "ZooKeeper"]

  master_nodes {
    flavor            = "c6.4xlarge.4.linux.mrs"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 480
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1
  }
  streaming_core_nodes {
    flavor            = "c6.4xlarge.4.linux.mrs"
    node_number       = 3
    root_volume_type  = "SAS"
    root_volume_size  = 480
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

### Create a hybrid cluster

```hcl
variable "mrs_az" {}
variable "cluster_name" {}
variable "password" {}
variable "keypair" {}
variable "vpc_id" {}
variable "subnet_id" {}

resource "flexibleengine_mrs_cluster_v2" "test" {
  availability_zone  = var.mrs_az
  name               = var.cluster_name
  version            = "MRS 2.0.1"
  type               = "MIXED"
  component_list     = ["Hadoop", "Spark", "Hive", "Tez", "Storm"]
  manager_admin_pwd  = var.password
  node_key_pair      = var.keypair
  vpc_id             = var.vpc_id
  subnet_id          = var.subnet_id

  master_nodes {
    flavor            = "c6.4xlarge.4.linux.mrs"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
  analysis_core_nodes {
    flavor            = "c6.4xlarge.4.linux.mrs"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
  streaming_core_nodes {
    flavor            = "c6.4xlarge.4.linux.mrs"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
  analysis_task_nodes {
    flavor            = "c6.4xlarge.4.linux.mrs"
    node_number       = 1
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
  streaming_task_nodes {
    flavor            = "c6.4xlarge.4.linux.mrs"
    node_number       = 1
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

### Create a custom cluster

```hcl
variable "mrs_az" {}
variable "cluster_name" {}
variable "password" {}
variable "keypair" {}
variable "vpc_id" {}
variable "subnet_id" {}

resource "flexibleengine_mrs_cluster_v2" "test" {
  availability_zone  = var.mrs_az
  name               = var.cluster_name
  version            = "MRS 3.1.0-LTS.1"
  type               = "CUSTOM"
  safe_mode          = true
  manager_admin_pwd  = var.password
  node_key_pair      = var.keypair
  vpc_id             = var.vpc_id
  subnet_id          = var.subnet_id
  template_id        = "mgmt_control_combined_v4"
  component_list     = ["DBService", "Hadoop", "ZooKeeper", "Ranger"]

  master_nodes {
    flavor            = "c6.4xlarge.4.linux.mrs"
    node_number       = 3
    root_volume_type  = "SAS"
    root_volume_size  = 480
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1
    assigned_roles = [
      "OMSServer:1,2",
      "SlapdServer:1,2",
      "KerberosServer:1,2",
      "KerberosAdmin:1,2",
      "quorumpeer:1,2,3",
      "NameNode:2,3",
      "Zkfc:2,3",
      "JournalNode:1,2,3",
      "ResourceManager:2,3",
      "JobHistoryServer:3",
      "DBServer:1,3",
      "HttpFS:1,3",
      "TimelineServer:3",
      "RangerAdmin:1,2",
      "UserSync:2",
      "TagSync:2",
      "KerberosClient",
      "SlapdClient",
      "meta"
    ]
  }

  custom_nodes {
    group_name        = "node_group_1"
    flavor            = "c6.4xlarge.4.linux.mrs"
    node_number       = 4
    root_volume_type  = "SAS"
    root_volume_size  = 480
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1
    assigned_roles = [
      "DataNode",
      "NodeManager",
      "KerberosClient",
      "SlapdClient",
      "meta"
    ]
  }
}

```

### Create an analysis cluster and bind public IP

```hcl
variable "mrs_az" {}
variable "cluster_name" {}
variable "password" {}
variable "keypair" {}
variable "vpc_id" {}
variable "subnet_id" {}
variable "public_ip" {}

resource "flexibleengine_mrs_cluster_v2" "test" {
  availability_zone  = var.mrs_az
  name               = var.cluster_name
  version            = "MRS 2.0.1"
  type               = "ANALYSIS"
  component_list     = ["Hadoop", "Hive", "Tez"]
  manager_admin_pwd  = var.password
  node_key_pair      = var.keypair
  vpc_id             = var.vpc_id
  subnet_id          = var.subnet_id
  public_ip          = var.public_ip

  master_nodes {
    flavor            = "c6.4xlarge.4.linux.mrs"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
  analysis_core_nodes {
    flavor            = "c6.4xlarge.4.linux.mrs"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
  analysis_task_nodes {
    flavor            = "c6.4xlarge.4.linux.mrs"
    node_number       = 1
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

<!--markdownlint-disable MD033-->

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the MRS cluster resource. If omitted, the
  provider-level region will be used. Changing this will create a new MRS cluster resource.

* `availability_zone` - (Required, String, ForceNew) Specifies the availability zone in which to create the cluster.
  Changing this will create a new MRS cluster resource.

* `name` - (Required, String, ForceNew) Specifies the name of the MRS cluster. The name can contain 2 to 64
  characters, which may consist of letters, digits, underscores (_) and hyphens (-). Changing this will create a new
  MRS cluster resource.

* `version` - (Required, String, ForceNew) Specifies the MRS cluster version. Currently, `MRS 1.8.9`,
  `MRS 2.0.1`, and `MRS 3.1.0-LTS.1` are supported. Changing this will create a new MRS cluster resource.

* `type` - (Optional, String, ForceNew) Specifies the type of the MRS cluster. The valid values are *ANALYSIS*,
  *STREAMING*, *MIXED* and *CUSTOM* (supported in MRS 3.x only), default to *ANALYSIS*.
  Changing this will create a new MRS cluster resource.

* `component_list` - (Required, List, ForceNew) Specifies the list of component names.
  Changing this will create a new MRS cluster resource. The supported components are as follows:

    <table border="2">
      <tr>
          <th>Cluster Version</th>
          <th>Cluster Type</th>
          <th>Components</th>
      </tr>
      <tr>
          <td rowspan="4">MRS 3.1.0-LTS.1</td>
          <td>analysis</td>
          <td>Hadoop, Spark2x, HBase, Hive, Hue, HetuEngine, Loader, Flink, Oozie, ZooKeeper, Ranger, and Tez</td>
      </tr>
      <tr>
          <td>streaming</td>
          <td>Kafka, Flume, ZooKeeper, and Ranger</td>
      </tr>
      <tr>
          <td>hybrid</td>
          <td>Hadoop, Spark2x, HBase, Hive, Hue, HetuEngine, Loader, Flink, Oozie, ZooKeeper, Ranger, Tez, Kafka, and Flume</td>
      </tr>
      <tr>
          <td>custom</td>
          <td>Hadoop, Spark2x, HBase, Hive, Hue, HetuEngine, Loader, Kafka, Flume, Flink, Oozie, ZooKeeper, Ranger, Tez,
          and ClickHouse</td>
      </tr>
      <tr>
          <td rowspan="2">MRS 2.0.1</td>
          <td>analysis</td>
          <td>Presto, Hadoop, Spark, HBase, Hive, Hue, Loader, and Tez</td>
      </tr>
      <tr>
          <td>streaming</td>
          <td>Kafka, Storm, and Flume</td>
      </tr>
      <tr>
          <td rowspan="2">MRS 1.8.9</td>
          <td>analysis</td>
          <td>Presto, Hadoop, Spark, HBase, Opentsdb, Hive, Hue, Loader, and Flink</td>
      </tr>
      <tr>
          <td>streaming</td>
          <td>Kafka, KafkaManager, Storm, and Flume</td>
      </tr>
    </table>

* `master_nodes` - (Required, List, ForceNew) Specifies a list of the informations about the master nodes in the
  MRS cluster.
  The nodes object structure of the `master_nodes` is documented below.
  Changing this will create a new MRS cluster resource.

* `vpc_id` - (Required, String, ForceNew) Specifies the ID of the VPC which bound to the MRS cluster.
  Changing this will create a new MRS cluster resource.

* `subnet_id` - (Required, String, ForceNew) Specifies the network ID of a subnet which bound to the MRS cluster.
  Changing this will create a new MRS cluster resource.

* `manager_admin_pwd` - (Required, String, ForceNew) Specifies the administrator password, which is used to login to
  the cluster management page. The password can contain 8 to 26 charactors and cannot be the username or the username
  spelled backwards. The password must contain lowercase letters, uppercase letters, digits, spaces and the special
  characters: `!?,.:-_{}[]@$^+=/`. Changing this will create a new MRS cluster resource.

* `node_key_pair` - (Required, String, ForceNew) Specifies the name of a key pair, which is used to login to the each
  nodes(ECSs). Changing this will create a new MRS cluster resource.

* `public_ip` - (Optional, String, ForceNew) Specifies the EIP address which bound to the MRS cluster.
  The EIP must have been created and must be in the same region as the cluster.
  Changing this will create a new MRS cluster resource.

* `eip_id` - (Optional, String, ForceNew) Specifies the EIP ID which bound to the MRS cluster.
  The EIP must have been created and must be in the same region as the cluster.
  Changing this will create a new MRS cluster resource.

* `log_collection` - (Optional, Bool, ForceNew) Specifies whether logs are collected when cluster installation fails.
  Default to true. If `log_collection` set true, the OBS buckets will be created and only used to collect logs that
  record MRS cluster creation failures. Changing this will create a new MRS cluster resource.

* `safe_mode` - (Optional, Bool, ForceNew) Specifies whether the running mode of the MRS cluster is secure,
  default to true.
  + true: enable Kerberos authentication.
  + false: disable Kerberos authentication. Changing this will create a new MRS cluster resource.

* `security_group_ids` - (Optional, List, ForceNew) Specifies an array of one or more security group ID to attach to the
  MRS cluster. If using the specified security group, the group need to open the specified port (9022) rules.

* `template_id` - (Optional, List, ForceNew) Specifies the template used for node deployment when the cluster type is
  CUSTOM.
  + mgmt_control_combined_v2: template for jointly deploying the management and control nodes. The management and
    control roles are co-deployed on the Master node, and data instances are deployed in the same node group. This
    deployment mode applies to scenarios where the number of control nodes is less than 100, reducing costs.
  + mgmt_control_separated_v2: The management and control roles are deployed on different master nodes, and data
    instances are deployed in the same node group. This deployment mode is applicable to a cluster with 100 to 500 nodes
    and delivers better performance in high-concurrency load scenarios.
  + mgmt_control_data_separated_v2: The management role and control role are deployed on different Master nodes,
    and data instances are deployed in different node groups. This deployment mode is applicable to a cluster with more
    than 500 nodes. Components can be deployed separately, which can be used for a larger cluster scale.

* `analysis_core_nodes` - (Optional, List) Specifies a list of the informations about the analysis core nodes in the
  MRS cluster.
  The nodes object structure of the `analysis_core_nodes` is documented below.

* `streaming_core_nodes` - (Optional, List) Specifies a list of the informations about the streaming core nodes in the
  MRS cluster.
  The nodes object structure of the `streaming_core_nodes` is documented below.

* `analysis_task_nodes` - (Optional, List) Specifies a list of the informations about the analysis task nodes in the
  MRS cluster.
  The nodes object structure of the `analysis_task_nodes` is documented below.

* `streaming_task_nodes` - (Optional, List) Specifies a list of the informations about the streaming task nodes in the
  MRS cluster.
  The nodes object structure of the `streaming_task_nodes` is documented below.

* `custom_nodes` - (Optional, List) Specifies a list of the informations about the custom nodes in the MRS cluster.
  The nodes object structure of the `custom_nodes` is documented below.
  Unlike other nodes, it needs to specify group_name.

* `tags` - (Optional, Map, ForceNew) Specifies the key/value pairs to associate with the cluster.

The `nodes` block supports:

* `flavor` - (Required, String, ForceNew) Specifies the instance specifications for each nodes in node group.
  Changing this will create a new MRS cluster resource.

* `node_number` - (Required, Int) Specifies the number of nodes for the node group.
  Only the core group and task group updations are allowed. The number of nodes after scaling cannot be
  less than the number of nodes originally created.

* `root_volume_type` - (Required, String, ForceNew) Specifies the system disk flavor of the nodes. Changing this will
  create a new MRS cluster resource.

* `root_volume_size` - (Required, Int, ForceNew) Specifies the system disk size of the nodes. Changing this will create
  a new MRS cluster resource.

* `data_volume_count` - (Required, Int, ForceNew) Specifies the data disk number of the nodes. The number configuration
  of each node are as follows:
  + master_nodes: 1.
  + analysis_core_nodes: minimum is one and the maximum is subject to the configuration of the corresponding flavor.
  + streaming_core_nodes: minimum is one and the maximum is subject to the configuration of the corresponding flavor.
  + analysis_task_nodes: minimum is zero and the maximum is subject to the configuration of the corresponding flavor.
  + streaming_task_nodes: minimum is zero and the maximum is subject to the configuration of the corresponding flavor.

  Changing this will create a new MRS cluster resource.
  
* `data_volume_type` - (Optional, String, ForceNew) Specifies the data disk flavor of the nodes.
  Required if `data_volume_count` is greater than zero. Changing this will create a new MRS cluster resource.
   The following disk types are supported:
  + `SATA`: common I/O disk
  + `SAS`: high I/O disk
  + `SSD`: ultra-high I/O disk

* `data_volume_size` - (Optional, Int, ForceNew) Specifies the data disk size of the nodes,in GB. The value range is 10
  to 32768. Required if `data_volume_count` is greater than zero. Changing this will create a new MRS cluster resource.

* `group_name` - (Optional, String, ForceNew) Specifies the name of nodes for the node group.
  This argument is mandatory when the cluster type is CUSTOM.

* `assigned_roles` - (Optional, List, ForceNew) Specifies the roles deployed in a node group. This argument is mandatory
  when the cluster type is CUSTOM. Each character string represents a role expression.

  **Role expression definition:**

   + If the role is deployed on all nodes in the node group, set this parameter to role_name, for example: `DataNode`.
   + If the role is deployed on a specified subscript node in the node group: role_name:index1,index2..., indexN,
     for example: `DataNode:1,2`. The subscript starts from 1.
   + Some roles support multi-instance deployment (that is, multiple instances of the same role are deployed on a node):
      role_name[instance_count], for example: `EsNode[9]`.
  
  [Mapping between roles and components](https://docs.prod-cloud-ocb.orange-business.com/api/mrs/mrs_02_0106.html)

  -> `DBService` is a basic component of a cluster. Components such as Hive, Hue, Oozie, Loader, and Redis, and Loader
   store their metadata in DBService, and provide the metadata backup and restoration functions by using DBService.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The cluster ID in UUID format.
* `total_node_number` - The total number of nodes deployed in the cluster.
* `master_node_ip` - The IP address of the master node.
* `private_ip` - The preferred private IP address of the master node.
* `status` - The cluster state, which include: running, frozen, abnormal and failed.
* `create_time` - The cluster creation time, in RFC-3339 format.
* `update_time` - The cluster update time, in RFC-3339 format.
* `charging_start_time` - The charging start time which is the start time of billing, in RFC-3339 format.
* `node` - all the nodes attributes: master_nodes/analysis_core_nodes/streaming_core_nodes/analysis_task_nodes
/streaming_task_nodes.
* `host_ips` - The host list of this nodes group in the cluster.

The `components` attributes:

* `id` - Component ID. For example, component_id of Hadoop is MRS 3.1.0-LTS.1_001,
  MRS 2.0.1_001, and MRS 1.8.9_001.
* `name` - Component name.
* `version` - Component version.
* `description` - Component description.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minute.
* `update` - Default is 180 minute.
* `delete` - Default is 40 minute.

## Import

Clusters can be imported by their `id`. For example,

```
terraform import flexibleengine_mrs_cluster_v2.test b11b407c-e604-4e8d-8bc4-92398320b847
```

Note that the imported state may not be identical to your resource definition, due to some attrubutes missing from the
API response, security or some other reason. The missing attributes include:
`manager_admin_pwd`, `template_id` and `assigned_roles`.
It is generally recommended running `terraform plan` after importing a cluster.
You can then decide if changes should be applied to the cluster, or the resource definition
should be updated to align with the cluster. Also you can ignore changes as below.

```hcl
resource "flexibleengine_mrs_cluster_v2" "test" {
    ...

  lifecycle {
    ignore_changes = [
      manager_admin_pwd,
    ]
  }
}
```
