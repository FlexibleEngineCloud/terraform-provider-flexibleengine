# FlexibleEngine Provider

The FlexibleEngine provider is used to interact with the
many resources supported by FlexibleEngine. The provider needs to be configured
with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the FlexibleEngine Provider
provider "flexibleengine" {
  domain_name = "admin"
  user_name   = "admin"
  password    = "pwd"
  region      = "eu-west-0"
}

# Create a web server
resource "flexibleengine_compute_instance_v2" "test-server" {
  # ...
}
```

## Authentication

### User name + Password

```hcl
provider "flexibleengine" {
  domain_name = var.domain_name
  user_name   = var.user_name
  password    = var.password
  region      = "eu-west-0"
}
```

### AKSK

```hcl
provider "flexibleengine" {
  access_key  = var.access_key
  secret_key  = var.secret_key
  domain_name = var.domain_name
  region      = "eu-west-0"
}
```

### Token

```hcl
provider "flexibleengine" {
  token       = var.token
  domain_name = var.domain_name
  tenant_name = var.tenant_name
  region      = "eu-west-0"
}
```

-> If token, aksk and password are set simultaneously, then it will authenticate in the order of Token, Password and AKSK.

### Federated

```hcl
provider "flexibleengine" {
  token          = var.token
  security_token = var.security_token
  access_key     = var.access_key
  secret_key     = var.secret_key
  domain_name    = var.domain_name
  tenant_name    = var.tenant_name
  region         = "eu-west-0"
}
```

## Configuration Reference

The following arguments are supported:

* `region` - (Required) The region of the FlexibleEngine cloud to use. It must be provided,
  but it can also be sourced from the `OS_REGION_NAME` environment variables.

* `domain_id` - (Optional; Required if not using `domain_name`) The ID of the Domain to scope to.
  If omitted, the following environment variables are checked (in this order):
  `OS_USER_DOMAIN_ID`, `OS_PROJECT_DOMAIN_ID`, `OS_DOMAIN_ID`.

* `domain_name` - (Optional; Required if not using `domain_id`) The Name of the Domain to scope to.
  If omitted, the following environment variables are checked (in this order):
  `OS_USER_DOMAIN_NAME`, `OS_PROJECT_DOMAIN_NAME`, `OS_DOMAIN_NAME`,
  `DEFAULT_DOMAIN`.

* `access_key` - (Optional) The access key of the FlexibleEngine cloud to use.
  If omitted, the `OS_ACCESS_KEY` environment variable is used.

* `secret_key` - (Optional) The secret key of the FlexibleEngine cloud to use.
  If omitted, the `OS_SECRET_KEY` environment variable is used.

* `user_name` - (Optional) The User name to login with. If omitted, the
  `OS_USER_NAME` environment variable is used.

* `user_id` - (Optional) The User ID to login with. If omitted, the
  `OS_USER_ID` environment variable is used.

* `password` - (Optional) The Password to login with. If omitted, the
  `OS_PASSWORD` environment variable is used.

* `tenant_id` - (Optional) The ID of the Project to login with.
  If omitted, the `OS_TENANT_ID` or `OS_PROJECT_ID` environment variables are used.

* `tenant_name` - (Optional) The Name of the Project to login with.
  If omitted, the `OS_TENANT_NAME`, `OS_PROJECT_NAME` environment variable or `region` is used.

* `token` - (Optional) A token is an expiring, temporary means of access issued via the
  IAM service. By specifying a token, you do not have to specify a username/password
  combination, since the token was already created by a username/password out of
  band of Terraform. If omitted, the `OS_AUTH_TOKEN` environment variable is used.

* `security_token` - (Optional) Security token to use for OBS federated authentication.

* `auth_url` - (Optional) The Identity authentication URL.
   If omitted, the `OS_AUTH_URL` environment variable is used.
   The default value is `https://iam.{{region}}.prod-cloud-ocb.orange-business.com/v3`.

* `max_retries` - (Optional) This is the maximum number of times an API
  call is retried, in the case where requests are being throttled or
  experiencing transient failures. The delay between the subsequent API
  calls increases exponentially. The default value is `5`.
  If omitted, the `OS_MAX_RETRIES` environment variable is used.

* `insecure` - (Optional) Trust self-signed SSL certificates. If omitted, the
  `OS_INSECURE` environment variable is used.

* `cacert_file` - (Optional) Specify a custom CA certificate when communicating
  over SSL. You can specify either a path to the file or the contents of the
  certificate. If omitted, the `OS_CACERT` environment variable is used.

* `cert` - (Optional) Specify client certificate file for SSL client
  authentication. You can specify either a path to the file or the contents of
  the certificate. If omitted the `OS_CERT` environment variable is used.

* `key` - (Optional) Specify client private key file for SSL client
  authentication. You can specify either a path to the file or the contents of
  the key. If omitted the `OS_KEY` environment variable is used.

## Logging

This provider has the ability to log all HTTP requests and responses between
Terraform and the FlexibleEngine cloud which is useful for troubleshooting and
debugging.

To enable these logs, set the `TF_LOG=DEBUG` environment variable:

```shell
TF_LOG=DEBUG terraform apply
```

If you submit these logs with a bug report, please ensure any sensitive
information has been scrubbed first!

## Testing and Development

In order to run the Acceptance Tests for development, the following environment
variables must also be set:

* `OS_REGION_NAME` - The region in which to create the server instance.

* `OS_ACCESS_KEY` - The access key of the FlexibleEngine cloud to use.

* `OS_SECRET_KEY` - The secret key of the FlexibleEngine cloud to use.

You should be able to use any FlexibleEngine environment to develop on as long as the
above environment variables are set.
