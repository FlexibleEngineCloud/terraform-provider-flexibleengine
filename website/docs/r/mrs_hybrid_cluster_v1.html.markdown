---
layout: "flexibleengine"
page_title: "FlexibleEngine: flexibleengine_mrs_hybrid_cluster_v1"
sidebar_current: "docs-flexibleengine-resource-mrs-hybrid-cluster-v1"
description: |-
  Manages resource cluster within FlexibleEngine MRS.
---

# flexibleengine\_mrs\_hybrid\_cluster\_v1

Manages resource cluster within FlexibleEngine MRS.

## Example Usage:  Creating a MRS hybrid cluster

```hcl
resource "flexibleengine_vpc_v1" "vpc_2" {
  name = "terraform_provider_vpc2"
  cidr= "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "subnet_1" {
  name = "flexibleengine_subnet"
  cidr = "192.168.0.0/16"
  gateway_ip = "192.168.0.1"
  vpc_id = flexibleengine_vpc_v1.vpc_2.id
}
resource "flexibleengine_mrs_hybrid_cluster_v1" "cluster1" {
  available_zone  = "eu-west-0a"
  cluster_name    = "mrs-hybrid-cluster-acc"
  cluster_version = "MRS 2.0.1"
  cluster_admin_secret  = "Cluster@123"
  master_node_key_pair = "KeyPair-ci"
  vpc_id = flexibleengine_vpc_v1.vpc_2.id
  subnet_id = flexibleengine_vpc_subnet_v1.subnet_1.id
  component_list = ["Hadoop", "Storm", "Spark", "Hive"]
  master_nodes {
    node_number = 1
    flavor = "s3.2xlarge.4.linux.mrs"
    data_volume_type = "SATA"
    data_volume_size = 100
    data_volume_count = 1
  }
  analysis_core_nodes {
    node_number = 1
    flavor = "s3.xlarge.4.linux.mrs"
    data_volume_type = "SATA"
    data_volume_size = 100
    data_volume_count = 1
  }
  streaming_core_nodes {
    node_number = 1
    flavor = "s3.xlarge.4.linux.mrs"
    data_volume_type = "SATA"
    data_volume_size = 100
    data_volume_count = 1
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) Cluster region information. Obtain the value from
    Regions and Endpoints.

* `available_zone` - (Required) ID or Name of an available zone. Obtain the value
    from Regions and Endpoints.

* `cluster_name` - (Required) Cluster name, which is globally unique and contains
    only 1 to 64 letters, digits, hyphens (-), and underscores (_).

* `cluster_version` - (Required) Version of the clusters. Currently, MRS 1.6.3, MRS 1.8.9, 
    and MRS 2.0.1 are supported. The latest version of MRS is used by default.
    Currently, the latest version is MRS 2.0.1.

* `vpc_id` - (Required) Specifies the id of the VPC.

* `subnet_id` - (Required) Specifies the id of the subnet.

* `safe_mode` - (Optional) MRS cluster running mode 
    - 0: common mode

        The value indicates that the Kerberos authentication is disabled. 
        Users can use all functions provided by the cluster. 
    - 1: safe mode (by default)

        The value indicates that the Kerberos authentication is enabled. 
        Common users cannot use the file management or job management functions of an MRS cluster 
        and cannot view cluster resource usage or the job records of Hadoop and Spark. To use these 
        functions, the users must obtain the relevant permissions from the MRS Manager administrator. 
        The request has the cluster_admin_secret parameter only when safe_mode is set to 1.

* `cluster_admin_secret` - (Required) Indicates the password of the MRS Manager administrator.
    - Must contain 8 to 32 characters.
    - Must contain at least three types of the following: Lowercase letters, Uppercase letters,
      Digits, Special characters of `~!@#$%^&*()-_=+\|[{}];:'",<.>/? and Spaces.
    - Must be different from the username.
    - Must be different from the username written in reverse order.
    For versions earlier than MRS 2.0.1, this parameter is mandatory only when safe_mode is set to 1.
    For MRS 2.0.1 or later, this parameter is mandatory no matter which value safe_mode is set to.

* `master_node_key_pair` - (Required) Name of a key pair You can use a key
    to log in to the Master node in the cluster.

* `security_group_id` - (Optional) Specifies the id of the security group which the cluster
    belongs to. If this parameter is empty, MRS automatically creates a security group, whose
    name starts with mrs_{cluster_name}.

* `log_collection` - (Optional) Indicates whether logs are collected when cluster
    installation fails. 0: not collected 1: collected The default value is 0. If
    log_collection is set to 1, OBS buckets will be created to collect the MRS logs.
    These buckets will be charged.

* `component_list` - (Required) Component name
    - Presto, Hadoop, Spark, HBase, Hive, Tez, Hue, Loader, Flume, Kafka and Storm are supported by MRS 2.0.1 or later.
    - Presto, Hadoop, Spark, HBase, Opentsdb, Hive, Hue, Loader, Flink, Flume, Kafka, KafkaManager and Storm are supported by MRS 1.8.9.
    - Hadoop, Spark, HBase, Hive, Hue, Loader, Flume, Kafka and Storm are supported by versions earlier than MRS 1.8.9.


* `master_nodes` - (Required) Specifies the master nodes information.

* `analysis_core_nodes` - (Required) Specifies the analysis core nodes information.

* `streaming_core_nodes` - (Required) Specifies the streaming core nodes information.

* `analysis_task_nodes` - (Optional) Specifies the analysis task nodes information.

* `streaming_task_nodes` - (Optional) Specifies the streaming task nodes information.


The `master_nodes`, `analysis_core_nodes`, `streaming_core_nodes`, `analysis_task_nodes`, `streaming_task_nodes` block supports:

* `flavor` - (Required) Best match based on several years of commissioning
    experience. MRS supports specifications of hosts, and host specifications are
    determined by CPUs, memory, and disks space.
    - Master nodes support s1.4xlarge and s1.8xlarge, c3.2xlarge.2, c3.xlarge.4, c3.2xlarge.4, c3.4xlarge.2, c3.4xlarge.4, c3.8xlarge.4, c3.15xlarge.4.
    - Core nodes of a streaming cluster support s1.xlarge, c2.2xlarge, s1.2xlarge, s1.4xlarge, s1.8xlarge, d1.8xlarge, , c3.2xlarge.2, c3.xlarge.4, c3.2xlarge.4, c3.4xlarge.2, c3.4xlarge.4, c3.8xlarge.4, c3.15xlarge.4.
    - Core nodes of an analysis cluster support all specifications c2.2xlarge, s1.xlarge, s1.4xlarge, s1.8xlarge, d1.xlarge, d1.2xlarge, d1.4xlarge, d1.8xlarge, , c3.2xlarge.2, c3.xlarge.4, c3.2xlarge.4, c3.4xlarge.2, c3.4xlarge.4, c3.8xlarge.4, c3.15xlarge.4, d2.xlarge.8, d2.2xlarge.8, d2.4xlarge.8, d2.8xlarge.8.

    The following provides specification details.

    node_size | CPU(core) | Memory(GB) | System Disk | Data Disk
    | --- | --- | --- | --- | --- |
    c2.2xlarge.linux.mrs | 8  | 16  | 40 | -
    cc3.xlarge.4.linux.mrs | 4  | 16  | 40 | -
    cc3.2xlarge.4.linux.mrs | 8  | 32  | 40 | -
    cc3.4xlarge.4.linux.mrs | 16 | 64  | 40 | -
    cc3.8xlarge.4.linux.mrs | 32 | 128 | 40 | -
    s1.xlarge.linux.mrs  | 4  | 16  | 40 | -
    s1.4xlarge.linux.mrs | 16 | 64  | 40 | -
    s1.8xlarge.linux.mrs | 32 | 128 | 40 | -
    s3.xlarge.4.linux.mrs| 4  | 16  | 40 | -
    s3.2xlarge.4.linux.mrs| 8 | 32  | 40 | -
    s3.4xlarge.4.linux.mrs| 16 | 64  | 40 | -
    d1.xlarge.linux.mrs  | 6  | 55  | 40 | 1.8 TB x 3 HDDs
    d1.2xlarge.linux.mrs | 12 | 110 | 40 | 1.8 TB x 6 HDDs
    d1.4xlarge.linux.mrs | 24 | 220 | 40 | 1.8 TB x 12 HDDs
    d1.8xlarge.linux.mrs | 48 | 440 | 40 | 1.8 TB x 24 HDDs
    d2.xlarge.linux.mrs  | 4  | 32  | 40 | -
    d2.2xlarge.linux.mrs | 8 | 64 | 40 | -
    d2.4xlarge.linux.mrs | 16 | 128 | 40 | 1.8TB*8HDDs
    d2.8xlarge.linux.mrs | 32 | 256 | 40 | 1.8TB*16HDDs
* `node_number` - (Required) Number of nodes. The value ranges from 0 to 500 and the default value is 0. 
    The total number of Core and Task nodes cannot exceed 500.
* `data_volume_type` - (Required) Data disk storage type of the node, supporting SATA and SSD currently
    - SATA: common I/O
    - SSD: Ultrahigh-speed I/O
* `data_volume_size` - (Required) Data disk size of the node
    Value range: 100 GB to 32000 GB
* `data_volume_count` - (Required) Number of data disks of the node
    Value range: 0 to 10


## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `available_zone` - See Argument Reference above.
* `cluster_name` - See Argument Reference above.
* `cluster_version` - See Argument Reference above.  
* `safe_mode` - See Argument Reference above.
* `cluster_admin_secret` - See Argument Reference above.
* `master_node_key_pair` - See Argument Reference above.
* `vpc_id` - See Argument Reference above.
* `subnet_id` - See Argument Reference above.
* `security_group_id` - See Argument Reference above.
* `log_collection` - See Argument Reference above.
* `master_nodes` - See Argument Reference above.
* `analysis_core_nodes` - See Argument Reference above.
* `streaming_core_nodes` - See Argument Reference above.
* `analysis_task_nodes` - See Argument Reference above.
* `streaming_task_nodes` - See Argument Reference above.
* `component_list` - See Argument Reference above.
* `billing_type` - The value is Metered, indicating on-demand payment.
* `total_node_number` - Total node number.
* `master_node_ip` - IP address of a Master node.
* `internal_ip` - Iternal IP address.
* `private_ip_first` - Primary private IP address.
* `external_ip` - External IP address.
* `external_alternate_ip` - Backup external IP address.
* `vnc` - URI address for remote login of the elastic cloud server.
* `state` - Cluster creation fee, which is automatically calculated.
* `create_at` - Cluster creation time.
* `update_at` - Cluster update time.
* `charging_start_time` - Time when charging starts.

The components attributes:

* `component_name` - Component name
* `component_id` - Component ID Component IDs supported by MRS 1.5.0 include:
    MRS 1.5.0_001: Hadoop MRS 1.5.0_002: Spark MRS 1.5.0_003: HBase MRS 1.5.0_004:
    Hive MRS 1.5.0_005: Hue MRS 1.5.0_006: Kafka MRS 1.5.0_007: Storm MRS 1.5.0_008:
    Loader MRS 1.5.0_009: Flume Component IDs supported by MRS 1.3.0 include: MRS
    1.3.0_001: Hadoop MRS 1.3.0_002: Spark MRS 1.3.0_003: HBase MRS 1.3.0_004: Hive
    MRS 1.3.0_005: Hue MRS 1.3.0_006: Kafka MRS 1.3.0_007: Storm For example, the
    component ID of Hadoop is MRS 1.5.0_001, or MRS 1.3.0_001.
* `component_version` - Component version MRS 1.5.0 supports the following component
    version: Component version of an analysis cluster: Hadoop: 2.7.2 Spark: 2.1.0
    HBase: 1.0.2 Hive: 1.2.1 Hue: 3.11.0 Loader: 2.0.0 Component version of a streaming
    cluster: Kafka: 0.10.0.0 Storm: 1.0.2 Flume: 1.6.0 MRS 1.3.0 supports the following
    component version: Component version of an analysis cluster: Hadoop: 2.7.2 Spark:
    1.5.1 HBase: 1.0.2 Hive: 1.2.1 Hue: 3.11.0 Component version of a streaming
    cluster: Kafka: 0.10.0.0 Storm: 1.0.2
* `component_desc` - Component description
