package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/signs"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getSignatureAssociateFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.ApigV2Client(OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}
	opts := signs.ListBindOpts{
		InstanceId:  state.Primary.Attributes["instance_id"],
		SignatureId: state.Primary.Attributes["signature_id"],
		Limit:       500,
	}
	resp, err := signs.ListBind(c, opts)
	if err != nil {
		return nil, err
	}
	if len(resp) < 1 {
		return nil, golangsdk.ErrDefault404{}
	}
	return resp, nil
}

func TestAccSignatureAssociate_basic(t *testing.T) {
	var (
		apiDetails []signs.SignBindApiInfo

		name   = acceptance.RandomAccResourceName()
		rName1 = "flexibleengine_apig_signature_associate.basic_bind"
		rName2 = "flexibleengine_apig_signature_associate.hmac_bind"
		rName3 = "flexibleengine_apig_signature_associate.aes_bind"

		rc1 = acceptance.InitResourceCheck(rName1, &apiDetails, getSignatureAssociateFunc)
		rc2 = acceptance.InitResourceCheck(rName2, &apiDetails, getSignatureAssociateFunc)
		rc3 = acceptance.InitResourceCheck(rName3, &apiDetails, getSignatureAssociateFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc1.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSignatureAssociate_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc1.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName1, "instance_id",
						"flexibleengine_apig_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName1, "signature_id",
						"flexibleengine_apig_signature.basic", "id"),
					resource.TestCheckResourceAttr(rName1, "publish_ids.#", "2"),
					rc2.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName2, "instance_id",
						"flexibleengine_apig_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName2, "signature_id",
						"flexibleengine_apig_signature.hmac", "id"),
					resource.TestCheckResourceAttr(rName2, "publish_ids.#", "2"),
					rc3.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName3, "instance_id",
						"flexibleengine_apig_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName3, "signature_id",
						"flexibleengine_apig_signature.aes", "id"),
					resource.TestCheckResourceAttr(rName3, "publish_ids.#", "2"),
				),
			},
			{
				Config: testAccSignatureAssociate_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc1.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName1, "publish_ids.#", "2"),
					rc2.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName2, "publish_ids.#", "2"),
					rc3.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName3, "publish_ids.#", "2"),
				),
			},
			{
				ResourceName:      rName1,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSignatureAssociateImportStateFunc(rName1),
			},
			{
				ResourceName:      rName2,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSignatureAssociateImportStateFunc(rName2),
			},
			{
				ResourceName:      rName3,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSignatureAssociateImportStateFunc(rName3),
			},
		},
	})
}

func testAccSignatureAssociateImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		if rs.Primary.Attributes["instance_id"] == "" || rs.Primary.Attributes["signature_id"] == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<signature_id>', but got '%s/%s'",
				rs.Primary.Attributes["instance_id"], rs.Primary.Attributes["signature_id"])
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.Attributes["signature_id"]), nil
	}
}

func testAccSignatureAssociate_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "flexibleengine_apig_instance" "test" {
  name                  = "%[2]s"
  edition               = "BASIC"
  vpc_id                = flexibleengine_vpc_v1.test.id
  subnet_id             = flexibleengine_vpc_subnet_v1.test.id
  security_group_id     = flexibleengine_networking_secgroup_v2.test.id
  enterprise_project_id = "0"

  availability_zones = try(slice(data.flexibleengine_availability_zones.test.names, 0, 1), null)
}

resource "flexibleengine_compute_instance_v2" "test" {
  name              = "%[2]s"
  image_id          = data.flexibleengine_images_image_v2.test.id
  flavor_id         = data.flexibleengine_compute_flavors_v2.test.flavors[0]
  availability_zone = data.flexibleengine_availability_zones.test.names[0]

  network {
    uuid = flexibleengine_vpc_subnet_v1.test.id
  }
}

resource "flexibleengine_apig_group" "test" {
  name        = "%[2]s"
  instance_id = flexibleengine_apig_instance.test.id
}

resource "flexibleengine_apig_vpc_channel" "test" {
  name        = "%[2]s"
  instance_id = flexibleengine_apig_instance.test.id
  port        = 80
  algorithm   = "WRR"
  protocol    = "HTTP"
  path        = "/"
  http_code   = "201"

  members {
    id = flexibleengine_compute_instance_v2.test.id
  }
}

resource "flexibleengine_apig_api" "test" {
  instance_id             = flexibleengine_apig_instance.test.id
  group_id                = flexibleengine_apig_group.test.id
  name                    = "%[2]s"
  type                    = "Public"
  request_protocol        = "HTTP"
  request_method          = "GET"
  request_path            = "/user_info/{user_age}"
  security_authentication = "APP"
  matching                = "Exact"
  success_response        = "Success response"
  failure_response        = "Failed response"
  description             = "Created by script"

  request_params {
    name     = "user_age"
    type     = "NUMBER"
    location = "PATH"
    required = true
    maximum  = 200
    minimum  = 0
  }
  
  backend_params {
    type     = "REQUEST"
    name     = "userAge"
    location = "PATH"
    value    = "user_age"
  }

  web {
    path             = "/getUserAge/{userAge}"
    vpc_channel_id   = flexibleengine_apig_vpc_channel.test.id
    request_method   = "GET"
    request_protocol = "HTTP"
    timeout          = 30000
  }

  web_policy {
    name             = "%[2]s_policy1"
    request_protocol = "HTTP"
    request_method   = "GET"
    effective_mode   = "ANY"
    path             = "/getUserAge/{userAge}"
    timeout          = 30000
    vpc_channel_id   = flexibleengine_apig_vpc_channel.test.id

    backend_params {
      type     = "REQUEST"
      name     = "userAge"
      location = "PATH"
      value    = "user_age"
    }

    conditions {
      source     = "param"
      param_name = "user_age"
      type       = "Equal"
      value      = "28"
    }
  }
}

resource "flexibleengine_apig_environment" "test" {
  count = 6

  name        = "%[2]s_${count.index}"
  instance_id = flexibleengine_apig_instance.test.id
}

resource "flexibleengine_apig_api_publishment" "test" {
  count = 6

  instance_id = flexibleengine_apig_instance.test.id
  api_id      = flexibleengine_apig_api.test.id
  env_id      = flexibleengine_apig_environment.test[count.index].id
}

resource "flexibleengine_apig_signature" "basic" {
  instance_id = flexibleengine_apig_instance.test.id
  name        = "%[2]s_basic"
  type        = "basic"
}

resource "flexibleengine_apig_signature" "hmac" {
  instance_id = flexibleengine_apig_instance.test.id
  name        = "%[2]s_hmac"
  type        = "hmac"
}

resource "flexibleengine_apig_signature" "aes" {
  instance_id = flexibleengine_apig_instance.test.id
  name        = "%[2]s_aes"
  type        = "aes"
  algorithm   = "aes-128-cfb"
}
`, testBaseComputeResources(name), name)
}

func testAccSignatureAssociate_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "flexibleengine_apig_signature_associate" "basic_bind" {
  instance_id  = flexibleengine_apig_instance.test.id
  signature_id = flexibleengine_apig_signature.basic.id
  publish_ids  = slice(flexibleengine_apig_api_publishment.test[*].publish_id, 0, 2)
}

resource "flexibleengine_apig_signature_associate" "hmac_bind" {
  instance_id  = flexibleengine_apig_instance.test.id
  signature_id = flexibleengine_apig_signature.hmac.id
  publish_ids  = slice(flexibleengine_apig_api_publishment.test[*].publish_id, 2, 4)
}

resource "flexibleengine_apig_signature_associate" "aes_bind" {
  instance_id  = flexibleengine_apig_instance.test.id
  signature_id = flexibleengine_apig_signature.aes.id
  publish_ids  = slice(flexibleengine_apig_api_publishment.test[*].publish_id, 4, 6)
}
`, testAccSignatureAssociate_base(name))
}

func testAccSignatureAssociate_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "flexibleengine_apig_signature_associate" "basic_bind" {
  instance_id  = flexibleengine_apig_instance.test.id
  signature_id = flexibleengine_apig_signature.basic.id
  publish_ids  = slice(flexibleengine_apig_api_publishment.test[*].publish_id, 1, 3)
}

resource "flexibleengine_apig_signature_associate" "hmac_bind" {
  instance_id  = flexibleengine_apig_instance.test.id
  signature_id = flexibleengine_apig_signature.hmac.id
  publish_ids  = slice(flexibleengine_apig_api_publishment.test[*].publish_id, 3, 5)
}

resource "flexibleengine_apig_signature_associate" "aes_bind" {
  instance_id  = flexibleengine_apig_instance.test.id
  signature_id = flexibleengine_apig_signature.aes.id
  publish_ids  = setunion(slice(flexibleengine_apig_api_publishment.test[*].publish_id, 0, 1),
    slice(flexibleengine_apig_api_publishment.test[*].publish_id, 5, 6))
}
`, testAccSignatureAssociate_base(name))
}
