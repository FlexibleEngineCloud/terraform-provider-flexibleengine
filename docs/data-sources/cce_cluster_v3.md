---
subcategory: "Cloud Container Engine (CCE)"
---

# flexibleengine_cce_cluster_v3

Provides details about a specified CCE cluster.

## Example Usage

 ```hcl
variable "cluster_name" {}

data "flexibleengine_cce_cluster_v3" "cluster" {
  name   = var.cluster_name
  status = "Available"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region in which to query the data source. If omitted, the provider-level region
  will be used.

* `name` - (Optional, String)The Name of the cluster resource.

* `id` - (Optional, String) The ID of container cluster.

* `status` - (Optional, String) The state of the cluster.

* `cluster_type` - (Optional, String) Type of the cluster. Possible values: VirtualMachine, BareMetal or Windows

* `vpc_id` - (Optional, String) The ID of the VPC used to create the node.

## Attribute Reference

All above argument parameters can be exported as attribute parameters along with attribute reference:

* `id` - The ID of the cluster.

* `name` - The name of the cluster in string format.

* `description` - Cluster description.

* `cluster_version` - The version of cluster in string format.

* `flavor_id` - The cluster specification in string format.

* `container_network_cidr` - The container network segment.

* `container_network_type` - The container network type: overlay_l2 , underlay_ipvlan or vpc-router.

* `service_network_cidr` - The service network segment.

* `custom_san` -  Custom san list for certificate. (array of string)

* `vpc_id` - The ID of the VPC used to create the node.

* `subnet_id` - The ID of the VPC Subnet used to create the node.

* `security_group_id` - Security group ID of the cluster.

* `highway_subnet_id` - The ID of the high speed network used to create bare metal nodes.

* `internal_endpoint` - The internal network address.

* `external_endpoint` - The external network address.

* `external_apig_endpoint` - The endpoint of the cluster to be accessed through API Gateway.

* `billingMode` - Charging mode of the cluster.

* `authentication_mode` - Authentication mode of the cluster, possible values are x509 and rbac.

* `masters` - Advanced configuration of master nodes.
  The [masters](#cce_masters) object structure is documented below.

<a name="cce_masters"></a>
The `masters` block supports:

* `availability_zone` - The availability zone (AZ) of the master node.
