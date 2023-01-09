---
subcategory: "Data Warehouse Service (DWS)"
description: ""
page_title: "flexibleengine_dws_cluster_v1"
---

# flexibleengine_dws_cluster_v1

Manages a DWS cluster resource within FlexibleEngine.

## Example Usage

```hcl
data "flexibleengine_dws_flavors" "flavor" {
  availability_zone = "eu-west-0a"
  vcpus             = 8
}

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

resource "flexibleengine_dws_cluster_v1" "cluster" {
  name              = "dws_cluster_test"
  node_type         = data.flexibleengine_dws_flavors.test.flavors[0].flavor_id
  number_of_node    = 3
  user_name         = "cluster_admin"
  user_pwd          = "Cluster123@!"
  vpc_id            = flexibleengine_vpc_v1.example_vpc.id
  subnet_id         = flexibleengine_vpc_subnet_v1.example_subnet.id
  security_group_id = flexibleengine_networking_secgroup_v2.example_secgroup.id
  availability_zone = "eu-west-0a"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Cluster name, which must be unique and contains 4 to 64
    characters, which consist of letters, digits, hyphens (-), or underscores
    (_) only and must start with a letter.

* `node_type` - (Required) Node type.

* `number_of_node` - (Required) Number of nodes in a cluster. The value ranges
    from 3 to 32.

* `user_name` - (Required) Administrator username for logging in to a data
    warehouse cluster The administrator username must:
    - Consist of lowercase letters, digits, or underscores.
    - Start with a lowercase letter or an underscore.
    - Contain 1 to 63 characters.
    - Cannot be a keyword of the DWS database.

* `user_pwd` - (Required) Administrator password for logging in to a data
    warehouse cluster. A password must conform to the following rules:
    - Contains 8 to 32 characters.
    - Cannot be the same as the username or the username written in reverse order.
    - Contains three types of lowercase letters, uppercase letters, digits and
      special characters ~!@#%^&*()-_=+|[{}];:,<.>/?

* `vpc_id` - (Required) VPC ID, which is used for configuring cluster network.

* `subnet_id` - (Required) The ID of the VPC Subnet, which is used for configuring cluster network.

* `security_group_id` - (Required) ID of a security group. The ID is used for
    configuring cluster network.

* `port` - (Optional) Service port of a cluster (8000 to 10000). The default value is 8000.

* `availability_zone` - (Optional) AZ in a cluster.

* `public_ip` - (Optional) Public IP address. The object structure is documented below.

The `public_ip` block supports:

* `public_bind_type` - (Optional) Binding type of an EIP. The value can be
    either of the following: *auto_assign*, *not_use* and *bind_existing*.
    The default value is *not_use*.

* `eip_id` - (Optional) EIP ID

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Cluster ID

* `endpoints` - The private network connection information about the cluster.
    The object structure is documented below.

* `public_endpoints` - The public network connection information about the cluster.
    The object structure is documented below.

* `private_ip` - List of private network IP address.

* `status` - Cluster status, which can be one of the following: *CREATING*, *AVAILABLE*, *UNAVAILABLE* and *CREATION FAILED*.

* `sub_status` - Sub-status of clusters in the AVAILABLE state.

* `task_status` - Cluster management task.

* `version` - Data warehouse version

* `created` - Cluster creation time. The format is ISO8601:YYYY-MM-DDThh:mm:ssZ.

* `updated` - Last modification time of a cluster. The format is ISO8601:YYYY-MM-DDThh:mm:ssZ.

The `endpoints` block supports:

* `connect_info` - Private network connection information

* `jdbc_url` - JDBC URL. The following is the default format:
    jdbc:postgresql://< connect_info>/<YOUR_DATABASE_NAME>

The `public_endpoints` block supports:

* `public_connect_info` - Public network connection information

* `jdbc_url` - JDBC URL. The following is the default format:
    jdbc:postgresql://< public_connect_info>/<YOUR_DATABASE_NAME>

## Import

DWS cluster can be imported using the `id`, e.g.

```shell
terraform import flexibleengine_dws_cluster_v1.cluster 1a2b3c4d-5e6f-7g8h-9i0j-1k2l3m4n5o6p
```
