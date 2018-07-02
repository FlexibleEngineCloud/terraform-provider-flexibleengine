---
layout: "flexibleengine"
page_title: "FlexibleEngine: flexibleengine_vpc_route_v2"
sidebar_current: "docs-flexibleengine-datasource-vpc-route-v2"
description: |-
  Provides details about a specific VPC Route.
---

# Data Source: flexibleengine_vpc_route_v2

`flexibleengine_vpc_route_v2` provides details about a specific VPC route.

## Example Usage

 ```hcl
 variable "route_id" { }

data "flexibleengine_vpc_route_v2" "vpc_route" {
  id = "${var.route_id}"
}

resource "flexibleengine_vpc_subnet_v1" "subnet_v1" {
  name = "test-subnet"
  cidr = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id = "${data.flexibleengine_vpc_route_v2.vpc_route.vpc_id}"
}

 ```

## Argument Reference

The arguments of this data source act as filters for querying the available
routes in the current tenant. The given filters must match exactly one
route whose data will be exported as attributes.

* `id` (Optional) - The id of the specific route to retrieve.

* `vpc_id` (Optional) - The id of the VPC that the desired route belongs to.

* `destination` (Optional) - The route destination address (CIDR).

* `tenant_id` (Optional) - Only the administrator can specify the tenant ID of other tenants.

* `type` (Optional) - Route type for filtering.

## Attribute Reference

All of the argument attributes are also exported as
result attributes.

* `nexthop` - The next hop of the route. If the route type is peering, it will provide VPC peering connection ID.