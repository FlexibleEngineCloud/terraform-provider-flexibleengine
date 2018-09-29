---
layout: "flexibleengine"
page_title: "FlexibleEngine: flexibleengine_dws_cluster_v1"
sidebar_current: "docs-flexibleengine-resource-dws-cluster-v1"
description: |-
  Manages a DWS cluster resource within FlexibleEngine.
---

# flexibleengine\_dws\_cluster\_v1

Manages a DWS cluster resource within FlexibleEngine

## Example Usage

```hcl
resource "flexibleengine_dws_cluster_v1" "cluster" {
  node_type = "dws.d1.xlarge"
  number_of_node = 3
  subnet_id = "{{ subnet_id }}"
  vpc_id = "{{ vpc_id }}"
  security_group_id = "{{ security_group_id }}"
  availability_zone = "{{ availability_zone }}"
  name = "terraform_dws_cluster_test"
  user_name = "test_cluster_admin"
  user_pwd = "cluster123@!"

  timeouts {
    create = "30m"
    delete = "30m"
  }
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional) AZ in a cluster

* `name` - (Required) Cluster name, which must be unique and contains 4 to 64
    characters, which consist of letters, digits, hyphens (-), or underscores
    (_) only and must start with a letter.

* `node_type` - (Required) Node type.

* `number_of_node` - (Required) Number of nodes in a cluster. The value ranges
    from 3 to 32.

* `port` - (Optional) Service port of a cluster (8000 to 10000). The default
    value is 8000.

* `public_ip` - (Optional) Public IP address. If the value is not specified,
    public connection is not used by default.

* `security_group_id` - (Required) ID of a security group. The ID is used for
    configuring cluster network.

* `subnet_id` - (Required) Subnet ID, which is used for configuring cluster
    network.

* `user_name` - (Required) Administrator username for logging in to a data
    warehouse cluster The administrator username must:

    Consist of lowercase letters, digits, or underscores.

    Start with a lowercase letter or an underscore.

    Contain 1 to 63 characters.

    Cannot be a keyword of the DWS database.

* `user_pwd` - (Required) Administrator password for logging in to a data
    warehouse cluster

    A password must conform to the following rules:

    Contains 8 to 32 characters.

    Cannot be the same as the username or the username written in reverse
    order.

    Contains three types of the following:

    Lowercase letters

    Uppercase letters

    Digits

    Special characters ~!@#%^&*()-_=+|[{}];:,<.>/?

* `vpc_id` - (Required) VPC ID, which is used for configuring cluster network.

The `public_ip` block supports:

* `eip_id` - (Optional) EIP ID

* `public_bind_type` - (Optional) Binding type of an EIP. The value can be
    either of the following:

    auto_assign

    not_use

    bind_existing

    The default value is not_use.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `name` - See Argument Reference above.
* `number_of_node` - See Argument Reference above.
* `availability_zone` - See Argument Reference above.
* `subnet_id` - See Argument Reference above.
* `user_name` - See Argument Reference above.
* `security_group_id` - See Argument Reference above.
* `public_ip` - See Argument Reference above.
* `node_type` - See Argument Reference above.
* `vpc_id` - See Argument Reference above.
* `port` - See Argument Reference above.

* `created` - Cluster creation time. The format is
    ISO8601:YYYY-MM-DDThh:mm:ssZ.

* `endpoints` - View the private network connection information about the
    cluster.

* `id` - Cluster ID

* `public_endpoints` - Public network connection information about the cluster.
    If the value is not specified, the public network connection information is
    not used by default.

* `status` - Cluster status, which can be one of the following:

    CREATING

    AVAILABLE

    UNAVAILABLE

    CREATION FAILED

* `sub_status` - Sub-status of clusters in the AVAILABLE state. The value can
    be one of the following:

    NORMAL

    READONLY

    REDISTRIBUTING

    REDISTRIBUTION-FAILURE

    UNBALANCED

    UNBALANCED | READONLY

    DEGRADED

    DEGRADED | READONLY

    DEGRADED | UNBALANCED

    UNBALANCED | REDISTRIBUTING

    UNBALANCED | REDISTRIBUTION-FAILURE

    READONLY | REDISTRIBUTION-FAILURE

    UNBALANCED | READONLY | REDISTRIBUTION-FAILURE

    DEGRADED | REDISTRIBUTION-FAILURE

    DEGRADED | UNBALANCED | REDISTRIBUTION-FAILURE

    DEGRADED | UNBALANCED | READONLY | REDISTRIBUTION-FAILURE

    DEGRADED | UNBALANCED | READONLY

* `task_status` - Cluster management task. The value can be one of the
    following:

    RESTORING

    SNAPSHOTTING

    GROWING

    REBOOTING

    SETTING_CONFIGURATION

    CONFIGURING_EXT_DATASOURCE

    DELETING_EXT_DATASOURCE

    REBOOT_FAILURE

    RESIZE_FAILURE

* `updated` - Last modification time of a cluster. The format is
    ISO8601:YYYY-MM-DDThh:mm:ssZ.

* `version` - Data warehouse version

The `endpoints` block supports:

* `connect_info` - Private network connection information

* `jdbc_url` - JDBC URL. The following is the default format:

    jdbc:postgresql://< connect_info>/<YOUR_DATABASE_NAME>

The `public_endpoints` block supports:

* `public_connect_info` - Public network connection information

* `jdbc_url` - JDBC URL. The following is the default format:

    jdbc:postgresql://< public_connect_info>/<YOUR_DATABASE_NAME>
