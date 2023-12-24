package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/acls"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getAclPolicyAssociateFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.ApigV2Client(OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}
	opt := acls.ListBindOpts{
		InstanceId: state.Primary.Attributes["instance_id"],
		PolicyId:   state.Primary.Attributes["policy_id"],
	}
	resp, err := acls.ListBind(c, opt)
	if len(resp) < 1 {
		return nil, golangsdk.ErrDefault404{}
	}
	return resp, err
}

func TestAccAclPolicyAssociate_basic(t *testing.T) {
	var (
		apiDetails []acls.AclBindApiInfo

		name  = acceptance.RandomAccResourceName()
		rName = "flexibleengine_apig_acl_policy_associate.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&apiDetails,
		getAclPolicyAssociateFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAclPolicyAssociate_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"flexibleengine_apig_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "policy_id",
						"flexibleengine_apig_acl_policy.test", "id"),
					resource.TestCheckResourceAttr(rName, "publish_ids.#", "1"),
				),
			},
			{
				Config: testAccAclPolicyAssociate_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"flexibleengine_apig_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "policy_id",
						"flexibleengine_apig_acl_policy.test", "id"),
					resource.TestCheckResourceAttr(rName, "publish_ids.#", "1"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAclPolicyAssociateImportStateFunc(rName),
			},
		},
	})
}

func testAccAclPolicyAssociateImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		if rs.Primary.Attributes["instance_id"] == "" || rs.Primary.Attributes["policy_id"] == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<policy_id>', but got '%s/%s'",
				rs.Primary.Attributes["instance_id"], rs.Primary.Attributes["policy_id"])
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.Attributes["policy_id"]), nil
	}
}

func testAccAclPolicyAssociate_base(name string) string {
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
  count = 2

  name        = "%[2]s_${count.index}"
  instance_id = flexibleengine_apig_instance.test.id
}

resource "flexibleengine_apig_api_publishment" "test" {
  count = 2

  instance_id = flexibleengine_apig_instance.test.id
  api_id      = flexibleengine_apig_api.test.id
  env_id      = flexibleengine_apig_environment.test[count.index].id
}

resource "flexibleengine_apig_acl_policy" "test" {
  instance_id = flexibleengine_apig_instance.test.id
  name        = "%[2]s"
  type        = "PERMIT"
  entity_type = "IP"
  value       = "10.201.33.4,10.30.2.15"
}
`, testBaseComputeResources(name), name)
}

func testAccAclPolicyAssociate_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "flexibleengine_apig_acl_policy_associate" "test" {
  instance_id = flexibleengine_apig_instance.test.id
  policy_id   = flexibleengine_apig_acl_policy.test.id

  publish_ids = [
    flexibleengine_apig_api_publishment.test[0].publish_id
  ]
}
`, testAccAclPolicyAssociate_base(name))
}

func testAccAclPolicyAssociate_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "flexibleengine_apig_acl_policy_associate" "test" {
  instance_id = flexibleengine_apig_instance.test.id
  policy_id   = flexibleengine_apig_acl_policy.test.id

  publish_ids = [
    flexibleengine_apig_api_publishment.test[1].publish_id
  ]
}
`, testAccAclPolicyAssociate_base(name))
}
