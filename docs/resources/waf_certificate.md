---
subcategory: "Web Application Firewall (WAF)"
---

# flexibleengine_waf_certificate

Manages a WAF certificate resource within FlexibleEngine.

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
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies the certificate name. The maximum length is 256 characters.
  Only digits, letters, underscores(`_`), and hyphens(`-`) are allowed.

* `certificate` - (Required, String, ForceNew) Specifies the certificate content. Changing this creates a new certificate.

* `private_key` - (Required, String, ForceNew) Specifies the private key. Changing this creates a new certificate.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The certificate ID in UUID format.

* `expiration` - Indicates the time when the certificate expires.

## Import

Certificates can be imported using the `id`, e.g.

```sh
terraform import flexibleengine_waf_certificate.cert_1 9251a0ed5aa640b68a35cf2eb6a3b733
```

Note that the imported state is not identical to your resource definition, due to security reason.
The missing attributes include `certificate`, and `private_key`. You can ignore changes as below.
```
resource "flexibleengine_waf_certificate" "cert_1" {
    ...

  lifecycle {
    ignore_changes = [
      certificate, private_key,
    ]
  }
}
```
