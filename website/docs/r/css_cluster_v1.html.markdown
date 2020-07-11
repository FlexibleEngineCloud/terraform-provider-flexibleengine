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

The `node_config` block supports:

* `availability_zone` -
  (Optional)
  Availability zone (AZ).  Changing this parameter will create a new resource.

* `flavor` -
  (Required)
  Instance flavor name. Value range of flavor ess.spec-1u8g: 40 GB
  to 640 GB Value range of flavor ess.spec-2u16g: 40 GB to 1280 GB
  Value range of flavor ess.spec-4u32g: 40 GB to 2560 GB Value
  range of flavor ess.spec-8u64g: 80 GB to 5120 GB Value range of
  flavor ess.spec-16u128g: 160 GB to 10240 GB.
  Changing this parameter will create a new resource.

* `network_info` -
  (Required)
  Network information. Structure is documented below. Changing this parameter will create a new resource.

* `volume` -
  (Required)
  Information about the volume. Structure is documented below. Changing this parameter will create a new resource.

The `network_info` block supports:

* `vpc_id` -
  (Required)
  VPC ID, which is used for configuring cluster network.  Changing this parameter will create a new resource.

* `subnet_id` -
  (Required)
  Subnet ID. All instances in a cluster must have the same
  subnets and security groups.  Changing this parameter will create a new resource.

* `security_group_id` -
  (Required)
  Security group ID. All instances in a cluster must have the
  same subnets and security groups.  Changing this parameter will create a new resource.

The `volume` block supports:

* `size` -
  (Required)
  Volume size, which must be a multiple of 4 and 10.  Changing this parameter will create a new resource.

* `volume_type` -
  (Required)
  COMMON: Common I/O. The SATA disk is used. HIGH: High I/O.
  The SAS disk is used. ULTRAHIGH: Ultra-high I/O. The
  solid-state drive (SSD) is used.  Changing this parameter will create a new resource.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `endpoint` -
  Indicates the IP address and port number.

* `created` -
  Time when a cluster is created. The format is ISO8601:
  CCYY-MM-DDThh:mm:ss.

* `nodes` -
  List of node objects. Structure is documented below.

The `nodes` block contains:

* `id` -
  Instance ID.

* `name` -
  Instance name.

* `type` -
  Supported type: ess (indicating the Elasticsearch node).

## Timeouts

This resource provides the following timeouts configuration options:

- `create` - Default is 30 minute.
- `update` - Default is 30 minute.
