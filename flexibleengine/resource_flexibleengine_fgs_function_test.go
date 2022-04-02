package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/fgs/v2/function"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getResourceObj(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.FgsV2Client(OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating FunctionGraph client: %s", err)
	}
	return function.GetMetadata(c, state.Primary.ID).Extract()
}

func TestAccFgsV2Function_basic(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	resourceName := "flexibleengine_fgs_function.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFgsV2Function_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "fuction test"),
					resource.TestCheckResourceAttrSet(resourceName, "urn"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
				),
			},
			{
				Config: testAccFgsV2Function_update(randName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "fuction test update"),
					resource.TestCheckResourceAttrSet(resourceName, "urn"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"app",
					"package",
					"func_code",
				},
			},
		},
	})
}

func TestAccFgsV2Function_text(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	resourceName := "flexibleengine_fgs_function.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFgsV2Function_text(randName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "fuction test"),
					resource.TestCheckResourceAttrSet(resourceName, "urn"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"app",
					"package",
					"func_code",
				},
			},
		},
	})
}

func TestAccFgsV2Function_agency(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	resourceName := "flexibleengine_fgs_function.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFgsV2Function_agency(randName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "agency", randName),
					resource.TestCheckResourceAttr(resourceName, "func_mounts.0.mount_type", "sfs"),
					resource.TestCheckResourceAttr(resourceName, "func_mounts.0.status", "active"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"xrole",
					"agency",
					"func_code",
				},
			},
		},
	})
}

func testAccFgsV2Function_basic(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_fgs_function" "test" {
  name        = "%s"
  app         = "default"
  description = "fuction test"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = "aW1wb3J0IGpzb24KZGVmIGhhbmRsZXIgKGV2ZW50LCBjb250ZXh0KToKICAgIG91dHB1dCA9ICdIZWxsbyBtZXNzYWdlOiAnICsganNvbi5kdW1wcyhldmVudCkKICAgIHJldHVybiBvdXRwdXQ="
}
`, rName)
}

func testAccFgsV2Function_update(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_fgs_function" "test" {
  name        = "%s"
  app         = "default"
  description = "fuction test update"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = "aW1wb3J0IGpzb24KZGVmIGhhbmRsZXIgKGV2ZW50LCBjb250ZXh0KToKICAgIG91dHB1dCA9ICdIZWxsbyBtZXNzYWdlOiAnICsganNvbi5kdW1wcyhldmVudCkKICAgIHJldHVybiBvdXRwdXQ="
}
`, rName)
}

func testAccFgsV2Function_text(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_fgs_function" "test" {
  name        = "%s"
  app         = "default"
  description = "fuction test"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"

  func_code = <<EOF
# -*- coding:utf-8 -*-
import json
def handler (event, context):
    return {
        "statusCode": 200,
        "isBase64Encoded": False,
        "body": json.dumps(event),
        "headers": {
            "Content-Type": "application/json"
        }
    }
EOF
}
`, rName)
}

func testAccFgsV2Function_agency(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_vpc_v1" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "test" {
  name       = "%s"
  cidr       = "192.168.1.0/24"
  gateway_ip = "192.168.1.1"
  vpc_id     = flexibleengine_vpc_v1.test.id
}

resource "flexibleengine_sfs_file_system_v2" "test" {
  share_proto = "NFS"
  size        = 10
  name        = "%s"
  description = "test sfs for fgs"
}

resource "flexibleengine_identity_agency_v3" "test" {
  name                   = "%s"
  description            = "test agency for fgs"
  delegated_service_name = "op_svc_cff"

  project_role {
    project = "%s"
    roles = [
      "VPC Administrator",
      "SFS Administrator",
    ]
  }
}

resource "flexibleengine_fgs_function" "test" {
  name        = "%s"
  package     = "default"
  description = "fuction test"
  handler     = "test.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = "aW1wb3J0IGpzb24KZGVmIGhhbmRsZXIgKGV2ZW50LCBjb250ZXh0KToKICAgIG91dHB1dCA9ICdIZWxsbyBtZXNzYWdlOiAnICsganNvbi5kdW1wcyhldmVudCkKICAgIHJldHVybiBvdXRwdXQ="
  agency      = flexibleengine_identity_agency_v3.test.name
  vpc_id      = flexibleengine_vpc_v1.test.id
  network_id  = flexibleengine_vpc_subnet_v1.test.id

  func_mounts {
    mount_type       = "sfs"
    mount_resource   = flexibleengine_sfs_file_system_v2.test.id
    mount_share_path = flexibleengine_sfs_file_system_v2.test.export_location
    local_mount_path = "/mnt"
  }
}
`, rName, rName, rName, rName, OS_REGION_NAME, rName)
}
