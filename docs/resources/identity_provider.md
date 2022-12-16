---
subcategory: "Identity and Access Management (IAM)"
description: ""
page_title: "flexibleengine_identity_provider"
---

# flexibleengine_identity_provider

Manages the identity providers within FlexibleEngine IAM service.

-> **NOTE:** You can create up to 10 identity providers.

## Example Usage

### Create a SAML protocol provider

```hcl
resource "flexibleengine_identity_provider" "provider_1" {
  name     = "saml_idp_demo"
  protocol = "saml"
}
```

### Create a OpenID Connect protocol provider

```hcl
resource "flexibleengine_identity_provider" "provider_2" {
  name     = "oidc_idp_demo"
  protocol = "oidc"
  
  openid_connect_config {
    access_type            = "program_console"
    provider_url           = "https://accounts.example.com"
    client_id              = "your_client_id"
    authorization_endpoint = "https://accounts.example.com/o/oauth2/v2/auth"
    scopes                 = ["openid"]
    signing_key            = jsonencode(
    {
      keys = [
        {
          alg = "RS256"
          e   = "AQAB"
          kid = "..."
          kty = "RSA"
          n   = "..."
          use = "sig"
        },
      ]
    }
    )
  }
}
```

<!--markdownlint-disable MD033-->
## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew) Specifies the name of the identity provider to be registered.
  The maximum length is 64 characters. Only letters, digits, underscores (_), and hyphens (-) are allowed.
  The name is unique, it is recommended to include domain name information.
  Changing this creates a new resource.

* `protocol` - (Required, String, ForceNew) Specifies the protocol of the identity provider.
  Valid values are *saml* and *oidc*.
  Changing this creates a new resource.

* `enabled` - (Optional, Bool) Specifies the status for the identity provider. Defaults to true.

* `description` - (Optional, String) Specifies the description of the identity provider.

* `metadata` - (Optional, String) Specifies the metadata of the IDP(Identity Provider) server.
  To obtain the metadata file of your enterprise IDP, contact the enterprise administrator.
  This field is used to import a metadata file to IAM to implement federated identity authentication.
  This field is required only if the protocol is set to *saml*.
  The maximum length is 30,000 characters and it stores in the state with SHA1 algorithm.

  -> **NOTE:**
    The metadata file specifies API addresses and certificate information in compliance with the SAML 2.0 standard.
    It is usually stored in a file. In the TF script, you can import the metafile through the **file** function,
    for example:
    <br/>`metadata = file("/usr/local/data/files/metadata.txt")`

* `openid_connect_config` - (Optional, List) Specifies the description of the identity provider.
  This field is required only if the protocol is set to *oidc*.

The `openid_connect_config` block supports:

* `access_type` - (Required, String) Specifies the access type of the identity provider.
  Available options are:
  + `program`: programmatic access only.
  + `program_console`: programmatic access and management console access.

* `provider_url` - (Required, String) Specifies the URL of the identity provider.
  This field corresponds to the iss field in the ID token.

* `client_id` - (Required, String) Specifies the ID of a client registered with the OpenID Connect identity provider.

* `signing_key` - (Required, String) Public key used to sign the ID token of the OpenID Connect identity provider.
  This field is required only if the protocol is set to *oidc*.

* `authorization_endpoint` - (Optional, String) Specifies the authorization endpoint of the OpenID Connect identity
  provider. This field is required only if the access type is set to `program_console`.

* `scopes` - (Optional, List) Specifies the scopes of authorization requests. It is an array of one or more scopes.
  Valid values are *openid*, *email*, *profile* and other values defined by you.
  This field is required only if the access type is set to `program_console`.

* `response_type` - (Optional, String) Response type. Valid values is *id_token*, default value is *id_token*.
  This field is required only if the access type is set to `program_console`.

* `response_mode` - (Optional, String) Response mode.
  Valid values is *form_post* and *fragment*, default value is *form_post*.
  This field is required only if the access type is set to `program_console`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which equals to the name.

* `login_link` - The login link of the identity provider.

* `sso_type` - The single sign-on type of the identity provider.

* `conversion_rules` - The identity conversion rules of the identity provider.
  The [object](#conversion_rules) structure is documented below

<a name="conversion_rules"></a>
The `conversion_rules` block supports:

* `local` - The federated user information on the cloud platform.

* `remote` - The description of the identity provider.

The `local` block supports:

* `username` - The name of a federated user on the cloud platform.

* `group` - The user group to which the federated user belongs on the cloud platform.

The `remote` block supports:

* `attribute` - The attribute in the IDP assertion.

* `condition` - The condition of conversion rule.

* `value` - The rule is matched only if the specified strings appear in the attribute type.

## Import

Identity provider can be imported using the `name`, e.g.

```
$ terraform import flexibleengine_identity_provider.provider_1 example_com_provider_saml
```
