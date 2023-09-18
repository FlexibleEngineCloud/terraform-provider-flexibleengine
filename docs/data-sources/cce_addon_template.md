---
subcategory: "Cloud Container Engine (CCE)"
---

# flexibleengine_cce_addon_template

Use this data source to get an available FlexibleEngine CCE add-on template.

## Example Usage

```hcl
variable "cluster_id" {}
variable "addon_name" {}
variable "addon_version" {}

data "flexibleengine_cce_addon_template" "test" {
  cluster_id = var.cluster_id
  name       = var.addon_name
  version    = var.addon_version
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region in which to query the data source. If omitted, the provider-level region
  will be used.

* `cluster_id` - (Required, String) Specifies the ID of CCE cluster.

* `name` - (Required, String) Specifies the add-on name. The supported addons are as follows:

  + **autoscaler**: AutoScaler is a component that automatically adjusts the size of a Kubernetes cluster so that all pods
    have a place to run and there are no unneeded nodes. Latest version: 1.19.6.

  + **coredns**: CoreDNS is a DNS server that chains plugins and provides Kubernetes DNS Services. Latest version: 1.17.7.

  + **everest**: Everest is a cloud native container storage system based on CSI, used to support cloud storage services
    for Kubernetes. Latest version: 1.2.9.

  + **metrics-server**: Metrics Server is a cluster-level resource usage data aggregator. Latest version: 1.1.2.

  + **gpu-beta**: A device plugin for nvidia.com/gpu resource on nvidia driver. Latest version: 1.2.2.

* `version` - (Required, String) Specifies the add-on version.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource id of the addon template.

* `description` - The description of the add-on.

* `spec` - The detail configuration of the add-on template.

* `stable` - Whether the add-on template is a stable version.

* `support_version` - The cluster information.
The [support_version](#cce_support_version) object structure is documented below.

<a name="cce_support_version"></a>
The `support_version` block supports:

* `virtual_machine` - The cluster (Virtual Machine) version that the add-on template supported.

* `bare_metal` - The cluster (Bare Metal) version that the add-on template supported.
