---
layout: "flexibleengine"
page_title: "FlexibleEngine: flexibleengine_css_cluster_v1"
sidebar_current: "docs-flexibleengine-resource-css-cluster-v1"
description: |-
 CSS cluster management
---

# flexibleengine\_css\_cluster\_v1

CSS cluster management

## Example Usage

### create a cluster

```hcl
resource "flexibleengine_networking_secgroup_v2" "secgroup" {
  name        = "terraform_test_security_group"
  description = "terraform security group acceptance test"
}

resource "flexibleengine_css_cluster_v1" "cluster" {
  name           = "terraform_test_cluster"
  engine_version = "7.1.1"
  node_number    = 1

  node_config {
    availability_zone = "{{ availability_zone }}"
    flavor            = "ess.spec-4u16g"

    network_info {
      vpc_id            = "{{ vpc_id }}"
      subnet_id         = "{{ network_id }}"
      security_group_id = flexibleengine_networking_secgroup_v2.secgroup.id
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

* `name` -
  (Required)
  Cluster name. It contains 4 to 32 characters. Only letters, digits,
  hyphens (-), and underscores (_) are allowed. The value must start
  with a letter. Changing this parameter will create a new resource.

* `engine_type` -
  (Optional)
  Engine type. The default value is "elasticsearch". Currently, the value
  can only be "elasticsearch". Changing this parameter will create a new resource.

* `engine_version` -
  (Required)
  Engine version. Versions 6.5.4 and 7.1.1 are supported. Changing this parameter will create a new resource.

* `node_number` -
  (Optional)
  Number of cluster instances. The value range is 1 to 32. Defaults to 1.

* `node_config` -
  (Required)
  Node configuration. Structure is documented below. Changing this parameter will create a new resource.

* `backup_strategy` - (Optional) Specifies the advanced backup policy. Structure is documented below.

* `tags` - (Optional) The key/value pairs to associate with the cluster.

The `node_config` block supports:

* `availability_zone` - (Optional)
  Availability zone(s). You can set multiple vailability zones, and use commas (,) to separate one from another.
  Cluster instances will be evenly distributed to each AZ. The `node_number` should be greater than or equal to
  the number of available zones. Changing this parameter will create a new resource.

* `flavor` - (Required)
  Instance flavor name. For example: value range of flavor ess.spec-2u8g:
  40 GB to 800 GB, value range of flavor ess.spec-4u16g: 40 GB to 1600 GB,
  value range of flavor ess.spec-8u32g: 80 GB to 3200 GB, value range of
  flavor ess.spec-16u64g: 100 GB to 6400 GB, value range of
  flavor ess.spec-32u128g: 100 GB to 10240 GB.
  Changing this parameter will create a new resource.

* `network_info` - (Required)
  Network information. Structure is documented below. Changing this parameter will create a new resource.

* `volume` - (Required)
  Information about the volume. Structure is documented below. Changing this parameter will create a new resource.

The `network_info` block supports:

* `vpc_id` - (Required)
  VPC ID, which is used for configuring cluster network. Changing this parameter will create a new resource.

* `subnet_id` -(Required)
  Subnet ID. All instances in a cluster must have the same subnet which should be configured with a **DNS address**.
  Changing this parameter will create a new resource.

* `security_group_id` - (Required)
  Security group ID. All instances in a cluster must have the same security group.
  Changing this parameter will create a new resource.

The `volume` block supports:

* `size` - (Required)
  Specifies volume size in GB, which must be a multiple of 10.

* `volume_type` - (Required)
  Specifies the volume type. Changing this parameter will create a new resource. Supported value:
  - "COMMON": The SATA disk is used;
  - "HIGH": The SAS disk is used;
  - "ULTRAHIGH": The solid-state drive (SSD) is used.

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

-> **NOTE:** `backup_strategy` requires the authority of *OBS Bucket* and *IAM Agency*.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `endpoint` -
  Indicates the IP address and port number.

* `created` -
  Time when a cluster is created. The format is ISO8601: CCYY-MM-DDThh:mm:ss.

* `nodes` -
  List of node objects. Structure is documented below.

The `nodes` block contains:

* `id` - Instance ID.

* `name` - Instance name.

* `type` - Supported type: ess (indicating the Elasticsearch node).

## Timeouts

This resource provides the following timeouts configuration options:

- `create` - Default is 60 minute.
- `update` - Default is 60 minute.
