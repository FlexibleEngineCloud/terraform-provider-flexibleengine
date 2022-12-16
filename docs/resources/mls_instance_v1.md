---
subcategory: "Machine Learning Service (MLS)"
description: ""
page_title: "flexibleengine_mls_instance_v1"
---

# flexibleengine_mls_instance_v1

Manages mls instance resource within FlexibleEngine

## Example Usage:  Creating a MLS instance

```hcl

resource "flexibleengine_mrs_cluster_v1" "cluster1" {
  cluster_name = "mrs-cluster-acc"
  region = "eu-west-0"
  billing_type = 12
  master_node_num = 2
  core_node_num = 3
  master_node_size = "s1.4xlarge.linux.mrs"
  core_node_size = "s1.xlarge.linux.mrs"
  available_zone_id = "eu-west-0a"
  vpc_id = "c1095fe7-03df-4205-ad2d-6f4c181d436e"
  subnet_id = "b65f8d25-c533-47e2-8601-cfaa265a3e3e"
  cluster_version = "MRS 1.3.0"
  volume_type = "SATA"
  volume_size = 100
  safe_mode = 0
  cluster_type = 0
  node_public_cert_name = "KeyPair-ci"
  cluster_admin_secret = ""
  component_list {
      component_name = "Hadoop"
  }
  component_list {
      component_name = "Spark"
  }
  component_list {
      component_name = "Hive"
  }
}

resource "flexibleengine_mls_instance_v1" "instance" {
  name = "terraform-mls-instance"
  version = "1.2.0"
  flavor = "mls.c2.2xlarge.common"
  network {
    vpc_id = "c1095fe7-03df-4205-ad2d-6f4c181d436e"
    subnet_id = "b65f8d25-c533-47e2-8601-cfaa265a3e3e"
    available_zone = "eu-west-0a"
    public_ip {
      bind_type = "not_use"
    }
  }
  mrs_cluster {
    id = flexibleengine_mrs_cluster_v1.cluster1.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to create the MLS instance. If
    omitted, the `region` argument of the provider is used. Changing this
    creates a new instance.

* `name` - (Required) Specifies the MLS instance name. The DB instance name of
    the same type is unique in the same tenant. Changing this creates a new instance.

* `version` - (Required) Specifies MLS Software version, only `1.2.0` is supported now.
  Changing this creates a new instance.

* `network` - (Required) Specifies the instance network information. The structure
  is described below. Changing this creates a new instance.

* `agency` - (Optional) Specifies the agency name. This parameter is mandatory only
  when you bind an instance to an elastic IP address (EIP). An instance must be
  bound to an EIP to grant MLS rights to abtain a tenant's token. NOTE: The tenant
  must create an agency on the Identity and Access Management (IAM) interface in
  advance. Changing this creates a new instance.

* `flavor` - (Required) Specifies the instance flavor, only `mls.c2.2xlarge.common`
  is supported now. Changing this creates a new instance.

* `mrs_cluster` - (Required) Specifies the MRS cluster information which the instance
  is associated. The structure is described below. NOTE: The current MRS instance
  requires an MRS cluster whose version is 1.3.0 and that is configured with the
  Spark component. MRS clusters whose version is not 1.3.0 or that are not configured
  with the Spark component cannot be selected. Changing this creates a new instance.

The `network` block supports:

* `vpc_id` - (Required) Specifies the ID of the virtual private cloud (VPC) where the
  instance resides. Changing this creates a new instance.

* `subnet_id` - (Required) Specifies the ID of the subnet where the instance resides.
  Changing this creates a new instance.

* `security_group` - (Optional) Specifies the ID of the security group of the instance.
  Changing this creates a new instance.

* `available_zone` - (Required) Specifies the AZ of the instance.
  Changing this creates a new instance.

* `public_ip` - (Required) Specifies the IP address of the instance. The structure is
  described below. Changing this creates a new instance.

The `public_ip` block supports:

* `bind_type` - (Required) Specifies the bind type. Possible values: `auto_assign` and
  `not_use`. Changing this creates a new instance.

The `mrs_cluster` block supports:

* `id` - (Required) Specifies the ID of the MRS cluster. Changing this creates a new instance.

* `user_name` - (Optional) Specifies the MRS cluster username. This parameter is mandatory
  only when the MRS cluster is in the security mode. Changing this creates a new instance.

* `user_password` - (Optional) Specifies the password of the MRS cluster user. The password
  and username work in a pair. Changing this creates a new instance.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `version` - See Argument Reference above.
* `agency` - See Argument Reference above.
* `flavor` - See Argument Reference above.
* `network/vpc_id` - See Argument Reference above.
* `network/subnet_id` - See Argument Reference above.
* `network/security_group` - See Argument Reference above.
* `network/available_zone` - See Argument Reference above.
* `network/public_ip/bind_type` - See Argument Reference above.
* `network/public_ip/eip_id` - Indicates the EIP ID, This is returned only when bind_type is
  set to auto_assign.
* `mrs_cluster` - See Argument Reference above.
* `status` - Indicates the MLS instance status.
* `inner_endpoint` - Indicates the URL for accessing the instance. Only machines in the same
  VPC and subnet as the instance can access the URL.
* `public_endpoint` - Indicates the URL for accessing the instance. The URL can be accessed
  from the Internet. The URL is created only after the instance is bound to an EIP.
* `created` - Indicates the creation time in the following format: yyyy-mm-dd Thh:mm:ssZ.
* `updated` - Indicates the update time in the following format: yyyy-mm-dd Thh:mm:ssZ.
