---
subcategory: "Deprecated"
description: ""
page_title: "flexibleengine_mrs_cluster_v1"
---

# flexibleengine_mrs_cluster_v1

Manages a MRS cluster resource within FlexibleEngine.

!> **Warning:** It has been deprecated, please use `flexibleengine_mrs_cluster_v2` instead.

## Example Usage:  Creating a MRS cluster

```hcl
resource "flexibleengine_vpc_v1" "example_vpc" {
  name = "example-vpc"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "example_subnet" {
  name       = "example-vpc-subnet"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = flexibleengine_vpc_v1.example_vpc.id
}

resource "flexibleengine_mrs_cluster_v1" "cluster1" {
  region            = "eu-west-0"
  available_zone_id = "eu-west-0a"
  cluster_name      = "mrs-cluster-test"
  cluster_type      = 0
  cluster_version   = "MRS 2.0.1"

  master_node_num  = 2
  core_node_num    = 3
  master_node_size = "s3.2xlarge.4.linux.mrs"
  core_node_size   = "s3.xlarge.4.linux.mrs"
  volume_type      = "SATA"
  volume_size      = 100
  vpc_id           = flexibleengine_vpc_v1.example_vpc.id
  subnet_id        = flexibleengine_vpc_subnet_v1.example_subnet.id

  safe_mode             = 0
  cluster_admin_secret  = "{{password_of_mrs_manager}}"
  node_public_cert_name = "KeyPair-ci"

  component_list {
    component_name = "Hadoop"
  }
  component_list {
    component_name = "Spark"
  }
  component_list {
    component_name = "Hive"
  }
  component_list {
    component_name = "Tez"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) Cluster region information. Obtain the value from
    Regions and Endpoints.

* `available_zone_id` - (Required) ID or Name of an available zone. Obtain the value
    from Regions and Endpoints.

* `cluster_name` - (Required) Cluster name, which is globally unique and contains
    only 1 to 64 letters, digits, hyphens (-), and underscores (_).

* `master_node_num` - (Required) Number of Master nodes The value is 2.

* `master_node_size` - (Required) Best match based on several years of commissioning
    experience. MRS supports specifications of hosts, and host specifications are
    determined by CPUs, memory, and disks space.
    + Master nodes support s1.4xlarge and s1.8xlarge, c3.2xlarge.2, c3.xlarge.4, c3.2xlarge.4, c3.4xlarge.2, c3.4xlarge.4,
      c3.8xlarge.4, c3.15xlarge.4.
    + Core nodes of a streaming cluster support s1.xlarge, c2.2xlarge, s1.2xlarge, s1.4xlarge, s1.8xlarge, d1.8xlarge,
      c3.2xlarge.2, c3.xlarge.4, c3.2xlarge.4, c3.4xlarge.2, c3.4xlarge.4, c3.8xlarge.4, c3.15xlarge.4.
    + Core nodes of an analysis cluster support all specifications c2.2xlarge, s1.xlarge, s1.4xlarge, s1.8xlarge,
      d1.xlarge, d1.2xlarge, d1.4xlarge, d1.8xlarge, , c3.2xlarge.2, c3.xlarge.4, c3.2xlarge.4, c3.4xlarge.2,
      c3.4xlarge.4, c3.8xlarge.4, c3.15xlarge.4, d2.xlarge.8, d2.2xlarge.8, d2.4xlarge.8, d2.8xlarge.8.

    The following provides specification details.

    node_size | CPU(core) | Memory(GB) | System Disk | Data Disk
    --- | --- | --- | --- | ---
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

* `core_node_num` - (Required) Number of Core nodes Value range: 3 to 500. A
    maximum of 500 Core nodes are supported by default. If more than 500 Core nodes
    are required, contact technical support engineers or invoke background APIs
    to modify the database.

* `core_node_size` - (Required) Instance specification of a Core node Configuration
    method of this parameter is identical to that of master_node_size.

* `vpc_id` - (Required) ID of the VPC where the subnet locates Obtain the VPC
    ID from the management console as follows: Register an account and log in to
    the management console. Click Virtual Private Cloud and select Virtual Private
    Cloud from the left list. On the Virtual Private Cloud page, obtain the VPC
    ID from the list.

* `subnet_id` - (Required) Specifies the ID of the VPC Subnet which bound to the MRS cluster.
    Changing this will create a new MRS cluster resource.

* `volume_type` - (Required) Type of disks SATA and SSD are supported. SATA:
    common I/O SSD: super high-speed I/O

* `volume_size` - (Required) Data disk storage space of a Core node Users can
    add disks to expand storage capacity when creating a cluster. There are the
    following scenarios: Separation of data storage and computing: Data is stored
    in the OBS system. Costs of clusters are relatively low but computing performance
    is poor. The clusters can be deleted at any time. It is recommended when data
    computing is not frequently performed. Integration of data storage and computing:
    Data is stored in the HDFS system. Costs of clusters are relatively high but
    computing performance is good. The clusters cannot be deleted in a short term.
    It is recommended when data computing is frequently performed. Value range:
    100 GB to 32000 GB

* `node_public_cert_name` - (Required) Name of a key pair You can use a key
    to log in to the Master node in the cluster.

* `safe_mode` - (Required) MRS cluster running mode.
    + **0**: common mode The value indicates that the Kerberos authentication is disabled.
      Users can use all functions provided by the cluster.
    + **1**: safe mode The value indicates that the Kerberos authentication is enabled.
      Common users cannot use the file management and job management functions of an MRS cluster or
      view cluster resource usage and the job records of Hadoop and Spark. To use these functions,
      the users must obtain the relevant permissions from the MRS Manager administrator.

* `cluster_admin_secret` - (Required) Indicates the password of the MRS Manager
    administrator.
    + Must contain 8 to 32 characters.
    + Must contain at least three of the following: Lowercase letters, Uppercase letters,
      Digits and Special characters: `~!@#$%^&*()-_=+\|[{}];:'",<.>/? and space
    + Cannot be the username or the username spelled backwards.

* `cluster_version` - (Optional) Version of the clusters. Possible values are as follows:
    MRS 1.8.9, MRS 2.0.1, MRS 2.1.0 and MRS 3.1.0-LTS.1. The latest version of MRS is used by default.

* `cluster_type` - (Optional) Type of clusters. 0: analysis cluster; 1: streaming cluster.
   The default value is 0.

* `log_collection` - (Optional) Indicates whether logs are collected when cluster
    installation fails. 0: not collected 1: collected The default value is 0. If
    log_collection is set to 1, OBS buckets will be created to collect the MRS logs.
    These buckets will be charged.

* `component_list` - (Required) Service component list. The object structure is documented below.

* `add_jobs` - (Optional) You can submit a job when you create a cluster to save time and use MRS easily.
    Only one job can be added. The object structure is documented below.

The `component_list` block supports:

* `component_name` - (Required) the Component name.
    + MRS 3.1.0-LTS.1 supports the following components:
      - The analysis cluster contains the following components: Hadoop, Spark2x, HBase, Hive, Hue, HetuEngine,
        Loader, Flink, Oozie, ZooKeeper, Ranger, and Tez.
      - The streaming cluster contains the following components: Kafka, Flume, ZooKeeper, and Ranger.
    + MRS 2.0.1 supports the following components:
      - The analysis cluster contains the following components: Presto, Hadoop, Spark, HBase, Hive, Hue, Loader, and Tez
      - The streaming cluster contains the following components: Kafka, Storm, and Flume.
    + MRS 1.8.9 supports the following components:
      - The analysis cluster contains the following components: Presto, Hadoop, Spark, HBase, Opentsdb, Hive, Hue, Loader,
        and Flink.
      - The streaming cluster contains the following components: Kafka, KafkaManager, Storm, and Flume.

The `add_jobs` block supports:

* `job_type` - (Required) Job type. 1: MapReduce 2: Spark 3: Hive Script 4: HiveQL
    (not supported currently) 5: DistCp, importing and exporting data (not supported
    in this API currently). 6: Spark Script 7: Spark SQL, submitting Spark SQL statements
    (not supported in this API currently). NOTE: Spark and Hive jobs can be added
    to only clusters including Spark and Hive components.

* `job_name` - (Required) Job name It contains only 1 to 64 letters, digits,
    hyphens (-), and underscores (_). NOTE: Identical job names are allowed but
    not recommended.

* `jar_path` - (Required) Path of the .jar file or .sql file for program execution
    The parameter must meet the following requirements: Contains a maximum of 1023
    characters, excluding special characters such as ;|&><'$. The address cannot
    be empty or full of spaces. Starts with / or s3a://. Spark Script must end with
    .sql; while MapReduce and Spark Jar must end with .jar. sql and jar are case-insensitive.

* `arguments` - (Optional) Key parameter for program execution The parameter
    is specified by the function of the user's program. MRS is only responsible
    for loading the parameter. The parameter contains a maximum of 2047 characters,
    excluding special characters such as ;|&>'<$, and can be empty.

* `input` - (Optional) Path for inputting data, which must start with / or s3a://.
    A correct OBS path is required. The parameter contains a maximum of 1023 characters,
    excluding special characters such as ;|&>'<$, and can be empty.

* `output` - (Optional) Path for outputting data, which must start with / or
    s3a://. A correct OBS path is required. If the path does not exist, the system
    automatically creates it. The parameter contains a maximum of 1023 characters,
    excluding special characters such as ;|&>'<$, and can be empty.

* `job_log` - (Optional) Path for storing job logs that record job running status.
    This path must start with / or s3a://. A correct OBS path is required. The parameter
    contains a maximum of 1023 characters, excluding special characters such as
    ;|&>'<$, and can be empty.

* `shutdown_cluster` - (Optional) Whether to delete the cluster after the jobs
    are complete true: Yes false: No

* `file_action` - (Optional) Data import and export import export

* `submit_job_once_cluster_run` - (Required) true: A job is submitted when a
    cluster is created. false: A job is submitted separately. The parameter is set
    to true in this example.

* `hql` - (Optional) HiveQL statement

* `hive_script_path` - (Optional) SQL program path This parameter is needed
    by Spark Script and Hive Script jobs only and must meet the following requirements:
    Contains a maximum of 1023 characters, excluding special characters such as
    ;|&><'$. The address cannot be empty or full of spaces. Starts with / or s3a://.
    Ends with .sql. sql is case-insensitive.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
* `cluster_id` - Cluster ID.
* `available_zone_name` - Name of an availability zone.
* `instance_id` - Instance ID.
* `hadoop_version` - Hadoop version.
* `master_node_ip` - IP address of a Master node.
* `externalIp` - Internal IP address.
* `private_ip_first` - Primary private IP address.
* `external_ip` - External IP address.
* `slave_security_groups_id` - Standby security group ID.
* `security_groups_id` - Security group ID.
* `external_alternate_ip` - Backup external IP address.
* `master_node_spec_id` - Specification ID of a Master node.
* `core_node_spec_id` - Specification ID of a Core node.
* `master_node_product_id` - Product ID of a Master node.
* `core_node_product_id` - Product ID of a Core node.
* `vnc` - URI address for remote login of the elastic cloud server.
* `fee` - Cluster creation fee, which is automatically calculated.
* `deployment_id` - Deployment ID of a cluster.
* `cluster_state` - Cluster status. Valid values include: starting, running, terminated, failed, abnormal,
    terminating, frozen, scaling-out and scaling-in.
* `order_id` - Order ID for creating clusters.
* `tenant_id` - Project ID.
* `create_at` - Cluster creation time.
* `update_at` - Cluster update time.
* `duration` - Cluster subscription duration.
* `charging_start_time` - Time when charging starts.
* `remark` - Remarks of a cluster.
* `error_info` - Error information.
* `component_list` - See Argument Reference below.

The component_list attributes:

* `component_id` - Component ID. For example, component_id of Hadoop is MRS 3.1.0-LTS.1_001, MRS 2.1.0_001,
    MRS 2.0.1_001, and MRS 1.8.9_001.
* `component_name` - Component name.
* `component_version` - Component version.
* `component_desc` - Component description.
