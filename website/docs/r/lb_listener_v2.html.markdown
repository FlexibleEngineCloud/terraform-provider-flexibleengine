---
layout: "flexibleengine"
page_title: "FlexibleEngine: flexibleengine_lb_listener_v2"
sidebar_current: "docs-flexibleengine-resource-lb-listener-v2"
description: |-
  Manages a V2 listener resource within FlexibleEngine.
---

# flexibleengine\_lb\_listener\_v2

Manages a V2 listener resource within FlexibleEngine.

## Example Usage

```hcl
resource "flexibleengine_lb_listener_v2" "listener_1" {
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = "d9415786-5f1a-428b-b35f-2f1523e146d2"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Networking client.
    A Networking client is needed to create an . If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    Listener.

* `protocol` - (Required) The protocol - can either be TCP, UDP, HTTP or TERMINATED_HTTPS.
    Changing this creates a new Listener.

* `protocol_port` - (Required) The port on which to listen for client traffic.
    Changing this creates a new Listener.

* `tenant_id` - (Optional) Required for admins. The UUID of the tenant who owns
    the Listener.  Only administrative users can specify a tenant UUID
    other than their own. Changing this creates a new Listener.

* `loadbalancer_id` - (Required) The load balancer on which to provision this
    Listener. Changing this creates a new Listener.

* `name` - (Optional) Human-readable name for the Listener. Does not have
    to be unique.

* `default_pool_id` - (Optional) The ID of the default pool with which the
    Listener is associated. Changing this creates a new Listener.

* `description` - (Optional) Human-readable description for the Listener.

* `default_tls_container_ref` - (Optional) A reference to a Barbican Secrets
    container which stores TLS information. This is required if the protocol
    is `TERMINATED_HTTPS`. See
    [here](https://wiki.openstack.org/wiki/Network/LBaaS/docs/how-to-create-tls-loadbalancer)
    for more information.

* `sni_container_refs` - (Optional) A list of references to Barbican Secrets
    containers which store SNI information. See
    [here](https://wiki.openstack.org/wiki/Network/LBaaS/docs/how-to-create-tls-loadbalancer)
    for more information.

* `admin_state_up` - (Optional) The administrative state of the Listener.
    A valid value is true (UP) or false (DOWN).

## Attributes Reference

The following attributes are exported:

* `id` - The unique ID for the Listener.
* `protocol` - See Argument Reference above.
* `protocol_port` - See Argument Reference above.
* `tenant_id` - See Argument Reference above.
* `name` - See Argument Reference above.
* `default_port_id` - See Argument Reference above.
* `description` - See Argument Reference above.
* `connection_limit` - See Argument Reference above.
* `default_tls_container_ref` - See Argument Reference above.
* `sni_container_refs` - See Argument Reference above.
* `admin_state_up` - See Argument Reference above.
