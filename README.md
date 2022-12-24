Terraform FlexibleEngine Provider
============================

<!-- markdownlint-disable-next-line MD034 -->
- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<!-- markdownlint-disable-next-line MD033 -->
<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Quick Start
-----------

When using the FlexibleEngineCloud Provider with Terraform 0.13 and later, the
recommended approach is to declare Provider versions in the root module Terraform
configuration, using a `required_providers` block as per the following example.
For previous versions, please continue to pin the version within the provider block.

1. Add [FlexibleEngineCloud/flexibleengine](https://registry.terraform.io/providers/FlexibleEngineCloud/flexibleengine/latest/docs)
  to your `required_providers`.

    ```hcl
    # provider.tf
    terraform {
      required_version = ">= 0.13"

      required_providers {
        flexibleengine = {
          source = "FlexibleEngineCloud/flexibleengine"
          version = ">= 1.30.0"
        }
      }
    }
    ```

2. Run `terraform init -upgrade` to download/upgrade the provider.

3. Add the provider and [Authenticate](https://registry.terraform.io/providers/FlexibleEngineCloud/flexibleengine/latest/docs#authentication).

    + **AK/SK Authenticate**

    ```hcl
    # provider.tf

    # Configure the FlexibleEngine Provider with AK/SK
    provider "flexibleengine" {
      access_key  = "access key"
      secret_key  = "secret key"
      domain_name = "domain name"
      region      = "eu-west-0"
    }
    ```

    + **Username/Password Authenticate**

    ```hcl
    # provider.tf

    # Configure the FlexibleEngine Provider with Username/Password 
    provider "flexibleengine" {
      user_name   = "user name"
      password    = "password"
      domain_name = "domain name"
      region      = "eu-west-0"
    }
    ```

4. Create your first resource.

    ```hcl
    # main.tf

    # Create an Elastic Cloud Server resource
    resource "flexibleengine_compute_instance_v2" "test-server" {
      name        = "test-server"
      image_name  = "OBS Ubuntu 18.04"
      flavor_name = "t2.micro"
      key_pair    = "kp_ecs"
      security_groups = ["default"]
      network {
        uuid = "55534eaa-533a-419d-9b40-ec427ea7195a"
      }
    }
    ```

Developing the Provider
------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org)
installed on your machine (version 1.18+ is *required*). You'll also need to
correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as
adding `$GOPATH/bin` to your `$PATH`.

Building the Provider
-----------

1. Clone repository to *$GOPATH/src/github.com/FlexibleEngineCloud/terraform-provider-flexibleengine*
  with `go get` or `git clone`.

    ```sh
    go get github.com/FlexibleEngineCloud/terraform-provider-flexibleengine
    ```

    ```sh
    cd $GOPATH/src/github.com/FlexibleEngineCloud/terraform-provider-flexibleengine
    git clone git@github.com:FlexibleEngineCloud/terraform-provider-flexibleengine.git
    ```

2. Enter the provider directory and build the provider, run `make build`. This will build the provider and
  put the provider binary in the `$GOPATH/bin` directory.

    ```sh
    cd $GOPATH/src/github.com/FlexibleEngineCloud/terraform-provider-flexibleengine
    make build
    ```

3. In order to test the provider, you can simply run `make test`.

    ```sh
    make test
    ```

Acceptance Testing
-----------

Before making a Pull Request or a release, the resources and data sources shoule be
tested with acceptance tests.

The following environment variables are required before running the acceptance testing:

```sh
export OS_ACCESS_KEY=xxx
export OS_SECRET_KEY=xxx
export OS_REGION_NAME=xxx
export OS_IMAGE_ID=xxx
export OS_FLAVOR_ID=xxx
export OS_NETWORK_ID=xxx
```

Then we can run the acceptance tests with `make testacc`.

```sh
make testacc TEST='./flexibleengine' TESTARGS='-run TestAccXXXX'
```

**Note:** Acceptance tests create real resources, and often cost money to run.

[Debugging Providers](https://www.terraform.io/docs/extend/debugging.html)
-----------

Add the `TF_LOG` and `TF_LOG_PATH` environment variables to the system, and then you can view detailed logs.
For example, in a Linux operating system, run the following commands:

```sh
export TF_LOG=TRACE
export TF_LOG_PATH="./terraform.log"
```

License
-----------

Terraform-Provider-FlexibleEngine is under the Mozilla Public License 2.0. See the [LICENSE](LICENSE) file for details.
