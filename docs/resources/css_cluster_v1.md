---
subcategory: "Cloud Search Service (CSS)"
description: ""
page_title: "flexibleengine_css_cluster_v1"
---

# flexibleengine_css_cluster_v1

CSS cluster management

## Example Usage

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

resource "flexibleengine_networking_secgroup_v2" "example_secgroup" {
  name        = "example-secgroup"
  description = "My neutron security group"
}

resource "flexibleengine_css_cluster_v1" "cluster" {
  name           = "terraform_test_cluster"
  engine_version = "7.9.3"
  node_number    = 1

  node_config {
    availability_zone = "eu-west-0a"
    flavor            = "ess.spec-4u16g"

    network_info {
      vpc_id            = flexibleengine_vpc_v1.example_vpc.id
      subnet_id         = flexibleengine_vpc_subnet_v1.example_subnet.id
      security_group_id = flexibleengine_networking_secgroup_v2.example_secgroup.id
    }
    volume {
      volume_type = "COMMON"
      size        = 40
    }
  }

  tags = {
    foo = "bar"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the cluster name. It contains 4 to 32 characters. Only letters, digits,
  hyphens (-), and underscores (_) are allowed. The value must start with a letter.
  Changing this parameter will create a new resource.

* `engine_version` - (Required) Specifies the engine version. For example, `7.6.2` and `7.9.3`.
   For details, see CSS [Supported Cluster Versions](https://docs.prod-cloud-ocb.orange-business.com/api/css/css_03_0056.html).
   Changing this parameter will create a new resource.

* `node_config` - (Required) Specifies the node configuration. [Structure](#css_node_config_object) is documented below.
  Changing this parameter will create a new resource.

* `node_number` - (Optional) Specifies the number of cluster instances. The value range is 1 to 32. Defaults to 1.

* `engine_type` - (Optional) Specifies the engine type. The default value is `elasticsearch`. Currently, the value
  can only be "elasticsearch". Changing this parameter will create a new resource.

* `security_mode` - (Optional) Whether to enable communication encryption and security authentication.
  Available values include *true* and *false*. security_mode is disabled by default.
  Changing this parameter will create a new resource.

* `password` - (Optional) Specifies the password of the cluster administrator admin in security mode.
  This parameter is mandatory only when security_mode is set to true. Changing this parameter will create a new resource.
  The administrator password must meet the following requirements:
  - The password can contain 8 to 32 characters.
  - The password must contain at least 3 of the following character types: uppercase letters, lowercase letters,
    digits, and special characters (~!@#$%^&*()-_=+\\|[{}];:,<.>/?).

* `backup_strategy` - (Optional) Specifies the advanced backup policy.
  [Structure](#css_backup_strategy_object) is documented below.

  -> **NOTE:** `backup_strategy` requires the authority of *OBS Bucket* and *IAM Agency*.

* `tags` - (Optional) Specifies the key/value pairs to associate with the cluster.

<a name="css_node_config_object"></a>
The `node_config` block supports:

* `flavor` - (Required) Specifies the instance flavor name. For example: value range of flavor `ess.spec-2u8g`:
  40 GB to 800 GB; value range of flavor `ess.spec-4u16g`: 40 GB to 1600 GB; value range of flavor `ess.spec-8u32g`:
  80 GB to 3200 GB; value range of flavor `ess.spec-16u64g`: 100 GB to 6400 GB; value range of flavor `ess.spec-32u128g`:
  100 GB to 10240 GB. Changing this parameter will create a new resource.

* `network_info` - (Required) Specifies the network information. [Structure](#css_network_info_object) is documented below.
  Changing this parameter will create a new resource.

* `volume` - (Required) Specifies the information about the volume. [Structure]($css_volume_object) is documented below.
  Changing this parameter will create a new resource.

* `availability_zone` - (Optional) Specifies the availability zone(s). You can set multiple vailability zones,
  and use commas (,) to separate one from another. Cluster instances will be evenly distributed to each AZ.
  The `node_number` should be greater than or equal to the number of available zones.
  Changing this parameter will create a new resource.

<a name="css_network_info_object"></a>
The `network_info` block supports:

* `vpc_id` - (Required) Specifies the VPC ID, which is used for configuring cluster network.
  Changing this parameter will create a new resource.

* `subnet_id` -(Required) Specifies the ID of the VPC Subnet. All instances in a cluster must have the same
  subnet which should be configured with a **DNS address**. Changing this parameter will create a new resource.

* `security_group_id` - (Required) Specifies the security group ID. All instances in a cluster must have the same
  security group. Changing this parameter will create a new resource.

<a name="css_volume_object"></a>
The `volume` block supports:

* `size` - (Required) Specifies the volume size in GB, which must be a multiple of 10.

* `volume_type` - (Required) Specifies the volume type. Changing this parameter will create a new resource. Supported value:
  - **COMMON**: The SATA disk is used;
  - **HIGH**: The SAS disk is used;
  - **ULTRAHIGH**: The solid-state drive (SSD) is used.

<a name="css_backup_strategy_object"></a>
The `backup_strategy` block supports:

* `start_time` - (Required) Specifies the time when a snapshot is automatically
  created everyday. Snapshots can only be created on the hour. The time format is
  the time followed by the time zone, specifically, **HH:mm z**. In the format, HH:mm
  refers to the hour time and z refers to the time zone. For example, "00:00 GMT+01:00"
  and "01:00 GMT+03:00".

* `keep_days` - (Optional) Specifies the number of days to retain the generated
   snapshots. Snapshots are reserved for seven days by default.

* `prefix` - (Optional) Specifies the prefix of the snapshot that is automatically
  created. The default value is "snapshot".

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `endpoint` - Indicates the IP address and port number.

* `created` - Time when a cluster is created. The format is ISO8601: CCYY-MM-DDThh:mm:ss.

* `nodes` - List of node objects. [Structure](#css_al_nodes_object) is documented below.

<a name="css_al_nodes_object"></a>
The `nodes` block contains:

* `id` - Instance ID.

* `name` - Instance name.

* `type` - Supported type: ess (indicating the Elasticsearch node).

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minute.
* `update` - Default is 60 minute.
