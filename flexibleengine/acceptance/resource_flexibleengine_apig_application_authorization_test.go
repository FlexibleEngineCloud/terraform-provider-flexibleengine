package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/appauths"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getAppAuthFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.ApigV2Client(OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}

	opts := appauths.ListOpts{
		InstanceId: state.Primary.Attributes["instance_id"],
		AppId:      state.Primary.Attributes["application_id"],
	}
	resp, err := appauths.ListAuthorized(client, opts)
	if err != nil {
		return nil, err
	}
	if len(resp) < 1 {
		return nil, golangsdk.ErrDefault404{}
	}
	return resp, nil
}

func TestAccAppAuth_basic(t *testing.T) {
	var (
		authApis []appauths.ApiAuthInfo

		rName      = "flexibleengine_apig_application_authorization.test"
		rc         = acceptance.InitResourceCheck(rName, &authApis, getAppAuthFunc)
		baseConfig = testAccAppAuth_base()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAppAuth_basic_step1(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				Config: testAccAppAuth_basic_step2(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAppAuthImportIdFunc(rName),
			},
		},
	})
}

func testAccAppAuthImportIdFunc(rsName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rsName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rsName, rs)
		}

		instanceId := rs.Primary.Attributes["instance_id"]
		resourceId := rs.Primary.ID
		if instanceId == "" || resourceId == "" {
			return "", fmt.Errorf("missing some attributes, want '<instance_id>/<id>' (the format of resource ID is "+
				"'<env_id>/<application_id>'), but got '%s/%s'", instanceId, resourceId)
		}
		return fmt.Sprintf("%s/%s", instanceId, resourceId), nil
	}
}

func testAccAppAuth_base() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
%[1]s

resource "flexibleengine_compute_instance_v2" "test" {
  name              = "%[2]s"
  image_id          = data.flexibleengine_images_image_v2.test.id
  flavor_id         = data.flexibleengine_compute_flavors_v2.test.flavors[0]
  availability_zone = data.flexibleengine_availability_zones.test.names[0]

  network {
    uuid = flexibleengine_vpc_subnet_v1.test.id
  }
}

resource "flexibleengine_apig_instance" "test" {
  name                  = "%[2]s"
  edition               = "BASIC"
  vpc_id                = flexibleengine_vpc_v1.test.id
  subnet_id             = flexibleengine_vpc_subnet_v1.test.id
  security_group_id     = flexibleengine_networking_secgroup_v2.test.id
  enterprise_project_id = "0"

  availability_zones = try(slice(data.flexibleengine_availability_zones.test.names, 0, 1), null)
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
  count = 3

  instance_id             = flexibleengine_apig_instance.test.id
  group_id                = flexibleengine_apig_group.test.id
  name                    = "%[2]s_${count.index}"
  type                    = "Public"
  request_protocol        = "HTTP"
  request_method          = "GET"
  request_path            = "/user_info/${count.index}"
  security_authentication = "APP"
  matching                = "Exact"

  web {
    path             = "/getUserAge/${count.index}"
    vpc_channel_id   = flexibleengine_apig_vpc_channel.test.id
    request_method   = "GET"
    request_protocol = "HTTP"
    timeout          = 30000
  }
}

resource "flexibleengine_apig_environment" "test" {
  instance_id = flexibleengine_apig_instance.test.id
  name        = "%[2]s"
}

resource "flexibleengine_apig_api_publishment" "test" {
  count = 3

  instance_id = flexibleengine_apig_instance.test.id
  api_id      = flexibleengine_apig_api.test[count.index].id
  env_id      = flexibleengine_apig_environment.test.id
}

resource "flexibleengine_apig_application" "test" {
  instance_id = flexibleengine_apig_instance.test.id// flexibleengine_apig_instance.test.id
  name        = "%[2]s"
}
`, testBaseComputeResources(name), name)
}

func testAccAppAuth_basic_step1(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "flexibleengine_apig_application_authorization" "test" {
  depends_on = [flexibleengine_apig_api_publishment.test]

  instance_id    = flexibleengine_apig_instance.test.id
  application_id = flexibleengine_apig_application.test.id
  env_id         = flexibleengine_apig_environment.test.id
  api_ids        = slice(flexibleengine_apig_api.test[*].id, 0, 2)
}
`, baseConfig)
}

func testAccAppAuth_basic_step2(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "flexibleengine_apig_application_authorization" "test" {
  depends_on = [flexibleengine_apig_api_publishment.test]

  instance_id    = flexibleengine_apig_instance.test.id
  application_id = flexibleengine_apig_application.test.id
  env_id         = flexibleengine_apig_environment.test.id
  api_ids        = slice(flexibleengine_apig_api.test[*].id, 1, 3)
}
`, baseConfig)
}
