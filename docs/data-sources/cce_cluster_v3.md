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

* `name` -  (Optional)The Name of the cluster resource.
 
* `id` - (Optional) The ID of container cluster.

* `status` - (Optional) The state of the cluster.

* `cluster_type` - (Optional) Type of the cluster. Possible values: VirtualMachine, BareMetal or Windows

## Attributes Reference

All above argument parameters can be exported as attribute parameters along with attribute reference:

* `id` - The ID of the cluster.

* `name` - The name of the cluster in string format.

* `description` - Cluster description.

* `cluster_version` - The version of cluster in string format.

* `flavor_id` - The cluster specification in string format.

* `container_network_cidr` - The container network segment.

* `container_network_type` - The container network type: overlay_l2 , underlay_ipvlan or vpc-router.

* `vpc_id` - The ID of the VPC used to create the node.

* `subnet_id` - The ID of the subnet used to create the node.

* `security_group_id` - Security group ID of the cluster.

* `highway_subnet_id` - The ID of the high speed network used to create bare metal nodes.

* `internal_endpoint` - The internal network address.

* `external_endpoint` - The external network address.

* `external_apig_endpoint` - The endpoint of the cluster to be accessed through API Gateway.

* `billingMode` - Charging mode of the cluster.
