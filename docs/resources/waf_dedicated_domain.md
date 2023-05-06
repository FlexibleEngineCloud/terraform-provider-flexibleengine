---
subcategory: "Web Application Firewall (WAF)"
description: ""
page_title: "flexibleengine_waf_dedicated_domain"
---

# flexibleengine_waf_dedicated_domain

Manages a dedicated mode domain resource within Flexibleengine.

-> **NOTE:** All WAF resources depend on WAF instances, and the WAF instances need to be purchased before they can be
used. The dedicated mode domain name resource can be used in Dedicated Mode and ELB Mode.

## Example Usage

```hcl
resource "flexibleengine_vpc_v1" "example_vpc" {
  name = "example-vpc"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_waf_dedicated_certificate" "certificate_1" {
  name        = "cert_1"
  certificate = <<EOT
-----BEGIN CERTIFICATE-----
MIIFazCCA1OgAwIBAgIUN3w1KX8/T/HWVxZIOdHXPhUOnsAwDQYJKoZIhvcNAQEL
BQAwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM
...
dKvZbPEsygYRIjwyhHHUh/YXH8KDI/uu6u6AxDckQ3rP1BkkKXr5NPBGjVgM3ZI=
-----END CERTIFICATE-----
EOT
  private_key = <<EOT
-----BEGIN PRIVATE KEY-----
MIIJQQIBADANBgkqhkiG9w0BAQEFAASCCSswggknAgEAAoICAQC+9uwFVenCdPD9
5LWSWMuy4riZW718wxBpYV5Y9N8nM7N0qZLLdpImZrzBbaBldTI+AZGI3Nupuurw
...
s9urs/Kk/tbQhsEvu0X8FyGwo0zH6rG8apTFTlac+v4mJ4vlpxSvT5+FW2lgLISE
+4sM7kp0qO3/p+45HykwBY5iHq3H
-----END PRIVATE KEY-----
EOT
}


resource "flexibleengine_waf_dedicated_domain" "domain_1" {
  domain         = "www.example.com"
  certificate_id = flexibleengine_waf_dedicated_certificate.certificate_1.id

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "192.168.1.100"
    port            = 8080
    type            = "ipv4"
    vpc_id          = flexibleengine_vpc_v1.example_vpc.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the dedicated mode domain resource. If omitted,
  the provider-level region will be used. Changing this setting will push a new domain.

* `domain` - (Required, String, ForceNew) Specifies the domain name to be protected. For example, `www.example.com` or
  `*.example.com`. Changing this creates a new domain.

* `server` - (Required, List, ForceNew) The server configuration list of the domain. A maximum of 80 can be configured.
  The object structure is documented below.

* `certificate_id` - (Optional, String) Specifies the certificate ID. This parameter is mandatory when `client_protocol`
  is set to HTTPS.

* `policy_id` - (Optional, String, ForceNew) Specifies the policy ID associated with the domain. If not specified,
  a new policy will be created automatically. Changing this creates a new domain.

* `proxy` - (Optional, Bool) Specifies whether a proxy is configured. Default value is `false`.

  -> **NOTE:** WAF forwards only HTTP/S traffic. So WAF cannot serve your non-HTTP/S traffic, such as UDP, SMTP, FTP,
  and basically all other non-HTTP/S traffic. If a proxy such as public network ELB (or Nginx) has been used, set
  proxy `true` to ensure that the WAF security policy takes effect for the real source IP address.

* `keep_policy` - (Optional, Bool) Specifies whether to retain the policy when deleting a domain name.
  Defaults to `true`.

* `protect_status` - (Optional, Int) The protection status of domain, `0`: suspended, `1`: enabled.
  Default value is `1`.

The `server` block supports:

* `client_protocol` - (Required, String, ForceNew) Protocol type of the client. The options include `HTTP` and `HTTPS`.
  Changing this creates a new service.

* `server_protocol` - (Required, String, ForceNew) Protocol used by WAF to forward client requests to the server. The
  options include `HTTP` and `HTTPS`. Changing this creates a new service.

* `vpc_id` - (Required, String, ForceNew) The id of the vpc used by the server. Changing this creates a service.

* `type` - (Required, String, ForceNew) Server network type, IPv4 or IPv6. Valid values are: `ipv4` and `ipv6`. Changing
  this creates a new service.

* `address` - (Required, String, ForceNew) IP address or domain name of the web server that the client accesses. For
  example, `192.168.1.1` or `www.example.com`. Changing this creates a new service.

* `port` - (Required, Int, ForceNew) Port number used by the web server. The value ranges from 0 to 65535. Changing this
  creates a new service.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the domain.

* `certificate_name` - The name of the certificate used by the domain name.

* `access_status` - Whether a domain name is connected to WAF. Valid values are:
  + `0` - The domain name is not connected to WAF,
  + `1` - The domain name is connected to WAF.

* `protocol` - The protocol type of the client. The options are `HTTP` and `HTTPS`.

* `tls` - The TLS configuration of domain.

* `cihper` - The cipher suite of domain.

* `compliance_certification` - The compliance certifications of the domain, values are:
  + `pci_dss` - The status of domain PCI DSS, `true`: enabled, `false`: disabled.
  + `pci_3ds` - The status of domain PCI 3DS, `true`: enabled, `false`: disabled.

* `alarm_page` - The alarm page of domain. Valid values are:
  + `template_name` - The template of alarm page, values are: `default`, `custom` and `redirection`.
  + `redirect_url` - The redirection URL when `template_name` is set to `redirection`.

* `traffic_identifier` - The traffic identifier of domain. Valid values are:
  + `ip_tag` - The IP tag of traffic identifier.
  + `session_tag` - The session tag of traffic identifier.
  + `user_tag` - The user tag of traffic identifier.

## Import

Dedicated mode domain can be imported using the `id`, e.g.

```sh
terraform import flexibleengine_waf_dedicated_domain.domain_1 69e9a86becb4424298cc6bdeacbf69d5
```
