---
subcategory: "Elastic Load Balance (ELB)"
description: ""
page_title: "flexibleengine_lb_listener_v2"
---

# flexibleengine_lb_listener_v2

Manages an **enhanced** lb listener resource within FlexibleEngine.

## Example Usage

```hcl
resource "flexibleengine_lb_loadbalancer_v2" "lb_1" {
  vip_subnet_id = flexibleengine_vpc_subnet_v1.example_subnet.ipv4_subnet_id

  tags = {
    key = "value"
  }
}

resource "flexibleengine_lb_listener_v2" "listener_1" {
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = flexibleengine_lb_loadbalancer_v2.lb_1.id

  tags = {
    key = "value"
  }
}
```

## Example Usage of TERMINATED_HTTPS protocol

```hcl
resource "flexibleengine_lb_loadbalancer_v2" "loadbalancer_1" {
  name          = "loadbalancer_cert"
  vip_subnet_id = flexibleengine_vpc_subnet_v1.example_subnet.ipv4_subnet_id
}

resource "flexibleengine_elb_certificate" "certificate_1" {
  name        = "cert"
  domain      = "www.elb.com"
  private_key = <<EOT
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAwZ5UJULAjWr7p6FVwGRQRjFN2s8tZ/6LC3X82fajpVsYqF1x
qEuUDndDXVD09E4u83MS6HO6a3bIVQDp6/klnYldiE6Vp8HH5BSKaCWKVg8lGWg1
UM9wZFnlryi14KgmpIFmcu9nA8yV/6MZAe6RSDmb3iyNBmiZ8aZhGw2pI1YwR+15
MVqFFGB+7ExkziROi7L8CFCyCezK2/oOOvQsH1dzQ8z1JXWdg8/9Zx7Ktvgwu5PQ
M3cJtSHX6iBPOkMU8Z8TugLlTqQXKZOEgwajwvQ5mf2DPkVgM08XAgaLJcLigwD5
13koAdtJd5v+9irw+5LAuO3JclqwTvwy7u/YwwIDAQABAoIBACU9S5fjD9/jTMXA
DRs08A+gGgZUxLn0xk+NAPX3LyB1tfdkCaFB8BccLzO6h3KZuwQOBPv6jkdvEDbx
Nwyw3eA/9GJsIvKiHc0rejdvyPymaw9I8MA7NbXHaJrY7KpqDQyk6sx+aUTcy5jg
iMXLWdwXYHhJ/1HVOo603oZyiS6HZeYU089NDUcX+1SJi3e5Ke0gPVXEqCq1O11/
rh24bMxnwZo4PKBWdcMBN5Zf/4ij9vrZE+fFzW7vGBO48A5lvZxWU2U5t/OZQRtN
1uLOHmMFa0FIF2aWbTVfwdUWAFsvAOkHj9VV8BXOUwKOUuEktdkfAlvrxmsFrO/H
yDeYYPkCgYEA/S55CBbR0sMXpSZ56uRn8JHApZJhgkgvYr+FqDlJq/e92nAzf01P
RoEBUajwrnf1ycevN/SDfbtWzq2XJGqhWdJmtpO16b7KBsC6BdRcH6dnOYh31jgA
vABMIP3wzI4zSVTyxRE8LDuboytF1mSCeV5tHYPQTZNwrplDnLQhywcCgYEAw8Yc
Uk/eiFr3hfH/ZohMfV5p82Qp7DNIGRzw8YtVG/3+vNXrAXW1VhugNhQY6L+zLtJC
aKn84ooup0m3YCg0hvINqJuvzfsuzQgtjTXyaE0cEwsjUusOmiuj09vVx/3U7siK
Hdjd2ICPCvQ6Q8tdi8jV320gMs05AtaBkZdsiWUCgYEAtLw4Kk4f+xTKDFsrLUNf
75wcqhWVBiwBp7yQ7UX4EYsJPKZcHMRTk0EEcAbpyaJZE3I44vjp5ReXIHNLMfPs
uvI34J4Rfot0LN3n7cFrAi2+wpNo+MOBwrNzpRmijGP2uKKrq4JiMjFbKV/6utGF
Up7VxfwS904JYpqGaZctiIECgYA1A6nZtF0riY6ry/uAdXpZHL8ONNqRZtWoT0kD
79otSVu5ISiRbaGcXsDExC52oKrSDAgFtbqQUiEOFg09UcXfoR6HwRkba2CiDwve
yHQLQI5Qrdxz8Mk0gIrNrSM4FAmcW9vi9z4kCbQyoC5C+4gqeUlJRpDIkQBWP2Y4
2ct/bQKBgHv8qCsQTZphOxc31BJPa2xVhuv18cEU3XLUrVfUZ/1f43JhLp7gynS2
ep++LKUi9D0VGXY8bqvfJjbECoCeu85vl8NpCXwe/LoVoIn+7KaVIZMwqoGMfgNl
nEqm7HWkNxHhf8A6En/IjleuddS1sf9e/x+TJN1Xhnt9W6pe7Fk1
-----END RSA PRIVATE KEY-----
EOT

  certificate = <<EOT
-----BEGIN CERTIFICATE-----
MIIDpTCCAo2gAwIBAgIJAKdmmOBYnFvoMA0GCSqGSIb3DQEBCwUAMGkxCzAJBgNV
BAYTAnh4MQswCQYDVQQIDAJ4eDELMAkGA1UEBwwCeHgxCzAJBgNVBAoMAnh4MQsw
CQYDVQQLDAJ4eDELMAkGA1UEAwwCeHgxGTAXBgkqhkiG9w0BCQEWCnh4QDE2My5j
b20wHhcNMTcxMjA0MDM0MjQ5WhcNMjAxMjAzMDM0MjQ5WjBpMQswCQYDVQQGEwJ4
eDELMAkGA1UECAwCeHgxCzAJBgNVBAcMAnh4MQswCQYDVQQKDAJ4eDELMAkGA1UE
CwwCeHgxCzAJBgNVBAMMAnh4MRkwFwYJKoZIhvcNAQkBFgp4eEAxNjMuY29tMIIB
IjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAwZ5UJULAjWr7p6FVwGRQRjFN
2s8tZ/6LC3X82fajpVsYqF1xqEuUDndDXVD09E4u83MS6HO6a3bIVQDp6/klnYld
iE6Vp8HH5BSKaCWKVg8lGWg1UM9wZFnlryi14KgmpIFmcu9nA8yV/6MZAe6RSDmb
3iyNBmiZ8aZhGw2pI1YwR+15MVqFFGB+7ExkziROi7L8CFCyCezK2/oOOvQsH1dz
Q8z1JXWdg8/9Zx7Ktvgwu5PQM3cJtSHX6iBPOkMU8Z8TugLlTqQXKZOEgwajwvQ5
mf2DPkVgM08XAgaLJcLigwD513koAdtJd5v+9irw+5LAuO3JclqwTvwy7u/YwwID
AQABo1AwTjAdBgNVHQ4EFgQUo5A2tIu+bcUfvGTD7wmEkhXKFjcwHwYDVR0jBBgw
FoAUo5A2tIu+bcUfvGTD7wmEkhXKFjcwDAYDVR0TBAUwAwEB/zANBgkqhkiG9w0B
AQsFAAOCAQEAWJ2rS6Mvlqk3GfEpboezx2J3X7l1z8Sxoqg6ntwB+rezvK3mc9H0
83qcVeUcoH+0A0lSHyFN4FvRQL6X1hEheHarYwJK4agb231vb5erasuGO463eYEG
r4SfTuOm7SyiV2xxbaBKrXJtpBp4WLL/s+LF+nklKjaOxkmxUX0sM4CTA7uFJypY
c8Tdr8lDDNqoUtMD8BrUCJi+7lmMXRcC3Qi3oZJW76ja+kZA5mKVFPd1ATih8TbA
i34R7EQDtFeiSvBdeKRsPp8c0KT8H1B4lXNkkCQs2WX5p4lm99+ZtLD4glw8x6Ic
i1YhgnQbn5E0hz55OLu5jvOkKQjPCW+9Aa==
-----END CERTIFICATE-----
EOT
}

resource "flexibleengine_lb_listener_v2" "listener_1" {
  name                      = "listener_cert"
  protocol                  = "TERMINATED_HTTPS"
  protocol_port             = 8080
  loadbalancer_id           = flexibleengine_lb_loadbalancer_v2.loadbalancer_1.id
  default_tls_container_ref = flexibleengine_elb_certificate.certificate_1.id
}
```

<!--markdownlint-disable MD033-->
## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the listener resource.
    If omitted, the `region` argument of the provider is used.
    Changing this creates a new listener.

* `loadbalancer_id` - (Required, String, ForceNew) The load balancer on which to provision this
    listener. Changing this creates a new listener.

* `protocol` - (Required, String, ForceNew) The protocol - can either be TCP, UDP, HTTP or TERMINATED_HTTPS.
    Changing this creates a new listener.

* `protocol_port` - (Required, Int, ForceNew) The port on which to listen for client traffic.
    Changing this creates a new listener.

* `default_pool_id` - (Optional, String, ForceNew) The ID of the default pool with which the
    listener is associated. Changing this creates a new listener.

* `name` - (Optional, String) Human-readable name for the listener. Does not have
    to be unique.

* `description` - (Optional, String) Human-readable description for the listener.

* `tags` - (Optional, Map) The key/value pairs to associate with the listener.

* `http2_enable` - (Optional, Bool) Specifies whether to use HTTP/2. The default value is false.
    This parameter is valid only when the protocol is set to *TERMINATED_HTTPS*.

* `transparent_client_ip_enable` - (Optional, Bool) Specifies whether to pass source IP addresses of the clients to
  backend servers.
  + For TCP and UDP listeners, the value can be true or false, and the default value is false.
  + For HTTP and HTTPS listeners, the value can only be true.

* `idle_timeout` - (Optional, Int) Specifies the idle timeout duration, in seconds.
  + For TCP listeners, the value ranges from 10 to 4000, and the default value is 300.
  + For HTTP and HTTPS listeners, the value ranges from 1 to 300, and the default value is 60.
  + For UDP listeners, this parameter does not take effect.

* `request_timeout` - (Optional, Int) Specifies the timeout duration for waiting for a request from a client,
  in seconds. This parameter is available only for HTTP and HTTPS listeners. The value ranges from 1 to 300,
  and the default value is 60.

* `response_timeout` - (Optional, Int) Specifies the timeout duration for waiting for a request from a backend
  server, in seconds. This parameter is available only for HTTP and HTTPS listeners. The value ranges from 1 to 300,
  and the default value is 60.

* `default_tls_container_ref` - (Optional, String) A reference to a Barbican Secrets
    container which stores TLS information. This is required if the protocol
    is `TERMINATED_HTTPS`. See
    [here](https://wiki.openstack.org/wiki/Network/LBaaS/docs/how-to-create-tls-loadbalancer)
    for more information.

* `sni_container_refs` - (Optional, List) A list of references to Barbican Secrets
    containers which store SNI information. See
    [here](https://wiki.openstack.org/wiki/Network/LBaaS/docs/how-to-create-tls-loadbalancer)
    for more information.

* `tls_ciphers_policy` - (Optional, String) Specifies the security policy used by the listener.
    This parameter is valid only when the load balancer protocol is set to TERMINATED_HTTPS.
    The value can be tls-1-0, tls-1-1, tls-1-2, or tls-1-2-strict, and the default value is tls-1-0.
    For details of cipher suites for each security policy, see the table below.

<table>
  <tr>
    <th>Security Policy</th>
    <th>TLS Version</th>
    <th>Cipher Suite</th>
  </tr >
  <tr >
    <td>tls-1-0</td>
    <td>TLSv1.2 TLSv1.1 TLSv1</td>
    <td rowspan="3">ECDHE-RSA-AES256-GCM-SHA384:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-ECDSA-AES128-GCM-SHA256:AES128-GCM-SHA256:AES256-GCM-SHA384:ECDHE-ECDSA-AES128-SHA256:ECDHE-RSA-AES128-SHA256:AES128-SHA256:AES256-SHA256:ECDHE-ECDSA-AES256-SHA384:ECDHE-RSA-AES256-SHA384:ECDHE-ECDSA-AES128-SHA:ECDHE-RSA-AES128-SHA:ECDHE-RSA-AES256-SHA:ECDHE-ECDSA-AES256-SHA:AES128-SHA:AES256-SHA</td>
  </tr>
  <tr>
    <td>tls-1-1</td>
    <td>TLSv1.2 TLSv1.1</td>
  </tr>
  <tr>
    <td>tls-1-2</td>
    <td>TLSv1.2</td>
  </tr>
  <tr>
    <td >tls-1-2-strict</td>
    <td >TLSv1.2</td>
    <td >ECDHE-RSA-AES256-GCM-SHA384:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-ECDSA-AES128-GCM-SHA256:AES128-GCM-SHA256:AES256-GCM-SHA384:ECDHE-ECDSA-AES128-SHA256:ECDHE-RSA-AES128-SHA256:AES128-SHA256:AES256-SHA256:ECDHE-ECDSA-AES256-SHA384:ECDHE-RSA-AES256-SHA384</td>
  </tr>
</table>

## Attribute Reference

The following attributes are exported:

* `id` - The unique ID for the listener.
* `protocol` - See Argument Reference above.
* `protocol_port` - See Argument Reference above.
* `name` - See Argument Reference above.
* `default_port_id` - See Argument Reference above.
* `description` - See Argument Reference above.
* `http2_enable` - See Argument Reference above.
* `default_tls_container_ref` - See Argument Reference above.
* `sni_container_refs` - See Argument Reference above.
* `tls_ciphers_policy` - See Argument Reference above.
* `tags` - See Argument Reference above.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.
