---
subcategory: "Web Application Firewall (WAF)"
description: ""
page_title: "flexibleengine_waf_domain"
---

# flexibleengine_waf_domain

Manages a WAF domain resource within FlexibleEngine.

## Example Usage

```hcl
resource "flexibleengine_waf_certificate" "certificate_1" {
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

resource "flexibleengine_waf_domain" "domain_1" {
  domain          = "www.example.com"
  certificate_id  = flexibleengine_waf_certificate.certificate_1.id
  proxy           = true
  sip_header_name = "default"
  sip_header_list = ["X-Forwarded-For"]

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "90.84.181.77"
    port            = "8080"
  }
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String, ForceNew) Specifies the domain name to be protected. For example, www.example.com or *.example.com.
  Changing this creates a new domain.

* `server` - (Required, List) Specifies an array of origin web servers. The object structure is documented below.

* `certificate_id` - (Optional, String) Specifies the certificate ID.
  This parameter is mandatory when `client_protocol` is set to HTTPS.

* `policy_id` - (Optional, String, ForceNew) Specifies the policy ID associated with the domain.
  If not specified, a new policy will be created automatically. Changing this create a new domain.

* `keep_proxy` - (Optional, Bool) Specifies whether to retain the policy when deleting a domain name. Defaults to true.

* `proxy` - (Optional, Bool) Specifies whether a proxy is configured.

* `sip_header_name` - (Optional, String) Specifies the type of the source IP header.
  This parameter is required only when proxy is set to true. The options are as follows:
  *default*, *cloudflare*, *akamai*, and *custom*.

* `sip_header_list` - (Optional, List) Specifies an array of HTTP request header for identifying the real source IP address.
  This parameter is required only when proxy is set to true.
  + If `sip_header_name` is *default*, the value is ["X-Forwarded-For"].
  + If `sip_header_name` is *cloudflare*, the value is ["CF-Connecting-IP", "X-Forwarded-For"].
  + If `sip_header_name` is *akamai*, the value is ["True-Client-IP"].
  + If `sip_header_name` is *custom*, you can customize a value.

The `server` block supports:

* `client_protocol` - (Required, String) Protocol type of the client. The options are *HTTP* and *HTTPS*.

* `server_protocol` - (Required, String) Protocol used by WAF to forward client requests to the server.
  The options are *HTTP* and *HTTPS*.

* `address` - (Required, String) IP address or domain name of the web server that the client accesses.
  For example, 192.168.1.1 or www.a.com.

* `port` - (Required, Int) Port number used by the web server. The value ranges from 0 to 65535, for example, 8080.

## Attributes Reference

The following attributes are exported:

* `id` -  ID of the domain.

* `cname` - The CNAME value.

* `txt_code` - The TXT record. This attribute is returned only when proxy is set to true.

* `sub_domain` - The subdomain name. This attribute is returned only when proxy is set to true.

* `protect_status` - The WAF mode. -1: bypassed, 0: disabled, 1: enabled.

* `access_status` - Whether a domain name is connected to WAF.
  + 0: The domain name is not connected to WAF;
  + 1: The domain name is connected to WAF.

* `protocol` - The protocol type of the client. The options are HTTP, HTTPS, and HTTP&HTTPS.

## Import

Domains can be imported using the `id`, e.g.

```sh
terraform import flexibleengine_waf_domain.dom_1 7117d38e-4c8f-4624-a505-bd96b97d024c
```
