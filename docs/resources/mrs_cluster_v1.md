---
subcategory: "MapReduce Service (MRS)"
---

# flexibleengine\_mrs\_cluster\_v1

Manages resource cluster within FlexibleEngine MRS.

## Example Usage:  Creating a MRS cluster

```hcl
resource "flexibleengine_mrs_cluster_v1" "cluster1" {
  region = "eu-west-0"
  available_zone_id = "eu-west-0a"
  billing_type    = 12
  cluster_name    = "mrs-cluster-test"
  cluster_type    = 0
  cluster_version = "MRS 2.0.1"

  master_node_num  = 2
  core_node_num    = 3
  master_node_size = "s3.2xlarge.4.linux.mrs"
  core_node_size   = "s3.xlarge.4.linux.mrs"
  volume_type = "SATA"
  volume_size = 100
  vpc_id = "51edfb75-f9f0-4bbc-b4dc-21466b93f60d"
  subnet_id = "1d7a8646-43ee-455a-a3ab-40da87a1304c"

  safe_mode = 0
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

* `billing_type` - (Required) The value is 12, indicating on-demand payment.

* `region` - (Required) Cluster region information. Obtain the value from
    Regions and Endpoints.

* `master_node_num` - (Required) Number of Master nodes The value is 2.

* `master_node_size` - (Required) Best match based on several years of commissioning
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

* `core_node_num` - (Required) Number of Core nodes Value range: 3 to 500. A
    maximum of 500 Core nodes are supported by default. If more than 500 Core nodes
    are required, contact technical support engineers or invoke background APIs
    to modify the database.

* `core_node_size` - (Required) Instance specification of a Core node Configuration
    method of this parameter is identical to that of master_node_size.

* `available_zone_id` - (Required) ID of an available zone. Obtain the value
    from Regions and Endpoints.

* `cluster_name` - (Required) Cluster name, which is globally unique and contains
    only 1 to 64 letters, digits, hyphens (-), and underscores (_).

* `vpc_id` - (Required) ID of the VPC where the subnet locates Obtain the VPC
    ID from the management console as follows: Register an account and log in to
    the management console. Click Virtual Private Cloud and select Virtual Private
    Cloud from the left list. On the Virtual Private Cloud page, obtain the VPC
    ID from the list.

* `subnet_id` - (Required) Subnet ID Obtain the subnet ID from the management
    console as follows: Register an account and log in to the management console.
    Click Virtual Private Cloud and select Virtual Private Cloud from the left list.
    On the Virtual Private Cloud page, obtain the subnet ID from the list.

* `cluster_version` - (Optional) Version of the clusters. Currently, MRS 1.3.0, MRS 1.5.0,
    MRS 1.6.3, MRS 1.8.9, and MRS 2.0.1 are supported. The latest version of MRS is used by default.
   Currently, the latest version is MRS 2.0.1.

* `cluster_type` - (Optional) Type of clusters. 0: analysis cluster; 1: streaming cluster.
   The default value is 0.

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

* `safe_mode` - (Required) MRS cluster running mode 0: common mode The value
    indicates that the Kerberos authentication is disabled. Users can use all functions
    provided by the cluster. 1: safe mode The value indicates that the Kerberos
    authentication is enabled. Common users cannot use the file management or job
    management functions of an MRS cluster and cannot view cluster resource usage
    or the job records of Hadoop and Spark. To use these functions, the users must
    obtain the relevant permissions from the MRS Manager administrator. The request
    has the cluster_admin_secret parameter only when safe_mode is set to 1.

* `cluster_admin_secret` - (Optional) Indicates the password of the MRS Manager
    administrator. The password for MRS 1.5.0: Must contain 6 to 32 characters.
    Must contain at least two types of the following: Lowercase letters Uppercase
    letters Digits Special characters of `~!@#$%^&*()-_=+\|[{}];:'",<.>/? Spaces
    Must be different from the username. Must be different from the username written
    in reverse order. The password for MRS 1.3.0: Must contain 8 to 64 characters.
    Must contain at least four types of the following: Lowercase letters Uppercase
    letters Digits Special characters of `~!@#$%^&*()-_=+\|[{}];:'",<.>/? Spaces
    Must be different from the username. Must be different from the username written
    in reverse order. This parameter needs to be configured only when safe_mode
    is set to 1.

* `log_collection` - (Optional) Indicates whether logs are collected when cluster
    installation fails. 0: not collected 1: collected The default value is 0. If
    log_collection is set to 1, OBS buckets will be created to collect the MRS logs.
    These buckets will be charged.

* `component_list` - (Required) Service component list.

* `add_jobs` - (Optional) You can submit a job when you create a cluster to
    save time and use MRS easily. Only one job can be added.


The `component_list` block supports:

* `component_name` - (Required) Component name Currently, Hadoop, Spark, HBase,
    Hive, Hue, Loader, Flume, Kafka and Storm are supported. Loader and Flume are
    not supported by MRS 1.3.0.
* `componentId` - Component ID Component IDs supported by MRS 1.5.0 include:
    MRS 1.5.0_001: Hadoop MRS 1.5.0_002: Spark MRS 1.5.0_003: HBase MRS 1.5.0_004:
    Hive MRS 1.5.0_005: Hue MRS 1.5.0_006: Kafka MRS 1.5.0_007: Storm MRS 1.5.0_008:
    Loader MRS 1.5.0_009: Flume Component IDs supported by MRS 1.3.0 include: MRS
    1.3.0_001: Hadoop MRS 1.3.0_002: Spark MRS 1.3.0_003: HBase MRS 1.3.0_004: Hive
    MRS 1.3.0_005: Hue MRS 1.3.0_006: Kafka MRS 1.3.0_007: Storm For example, the
    component ID of Hadoop is MRS 1.5.0_001, or MRS 1.3.0_001.
* `componentName` - Component name Currently, Hadoop, Spark, HBase, Hive, Hue,
    Loader, Flume, Kafka and Storm are supported. Loader and Flume are not supported
    by MRS 1.3.0.
* `componentVersion` - Component version MRS 1.5.0 supports the following component
    version: Component version of an analysis cluster: Hadoop: 2.7.2 Spark: 2.1.0
    HBase: 1.0.2 Hive: 1.2.1 Hue: 3.11.0 Loader: 2.0.0 Component version of a streaming
    cluster: Kafka: 0.10.0.0 Storm: 1.0.2 Flume: 1.6.0 MRS 1.3.0 supports the following
    component version: Component version of an analysis cluster: Hadoop: 2.7.2 Spark:
    1.5.1 HBase: 1.0.2 Hive: 1.2.1 Hue: 3.11.0 Component version of a streaming
    cluster: Kafka: 0.10.0.0 Storm: 1.0.2
* `componentDesc` - Component description


The `add_jobs` block supports:
* `job_type` - (Required) Job type 1: MapReduce 2: Spark 3: Hive Script 4: HiveQL
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

The following attributes are exported:

* `billing_type` - See Argument Reference above.
* `data_center` - See Argument Reference above.
* `master_node_num` - See Argument Reference above.
* `master_node_size` - See Argument Reference above.
* `core_node_num` - See Argument Reference above.
* `core_node_size` - See Argument Reference above.
* `available_zone_id` - See Argument Reference above.
* `cluster_name` - See Argument Reference above.
* `vpc_id` - See Argument Reference above.
* `subnet_id` - See Argument Reference above.
* `cluster_version` - See Argument Reference above.
* `cluster_type` - See Argument Reference above.
* `volume_type` - See Argument Reference above.
* `volume_size` - See Argument Reference above.
* `node_public_cert_name` - See Argument Reference above.
* `safe_mode` - See Argument Reference above.
* `cluster_admin_secret` - See Argument Reference above.
* `log_collection` - See Argument Reference above.
* `component_list` - See Argument Reference below.
* `add_jobs` - See Argument Reference above.
* `order_id` - Order ID for creating clusters.
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
* `duration` - Cluster subscription duration.
* `vnc` - URI address for remote login of the elastic cloud server.
* `fee` - Cluster creation fee, which is automatically calculated.
* `deployment_id` - Deployment ID of a cluster.
* `cluster_state` - Cluster status Valid values include: existing history starting
    running terminated failed abnormal terminating rebooting shutdown frozen scaling-out
    scaling-in scaling-error.
* `tenant_id` - Project ID.
* `create_at` - Cluster creation time.
* `update_at` - Cluster update time.
* `error_info` - Error information.
* `charging_start_time` - Time when charging starts.
* `remark` - Remarks of a cluster.

The component_list attributes:
* `component_name` - (Required) Component name Currently, Hadoop, Spark, HBase,
    Hive, Hue, Loader, Flume, Kafka and Storm are supported. Loader and Flume are
    not supported by MRS 1.3.0.
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
