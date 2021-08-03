package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccRTSStackResourcesV1_dataSource(t *testing.T) {
	var stackName = fmt.Sprintf("terra-test-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRTSStackResourcesV1Config(stackName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRTSStackResourcesV1DataSourceID("data.flexibleengine_rts_stack_resource_v1.resource_1"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_rts_stack_resource_v1.resource_1", "resource_name", "random"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_rts_stack_resource_v1.resource_1", "resource_type", "OS::Heat::RandomString"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_rts_stack_resource_v1.resource_1", "resource_status", "CREATE_COMPLETE"),
				),
			},
		},
	})
}

func testAccCheckRTSStackResourcesV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find stack resource data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("stack resource data source ID not set")
		}

		return nil
	}
}

func testAccRTSStackResourcesV1Config(stackName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_rts_stack_v1" "stack_1" {
  name = "%s"
  disable_rollback= true
  timeout_mins=60
  template_body = <<JSON
          {
			"outputs": {
              "str1": {
                 "description": "The description of the nat server.",
                 "value": {
                   "get_resource": "random"
                 }
	          }
            },
            "heat_template_version": "2013-05-23",
            "description": "A HOT template that create a single server and boot from volume.",
            "parameters": {
              "key_name": {
                "type": "string",
                "description": "Name of existing key pair for the instance to be created.",
                "default": "KeyPair-click2cloud"
	          }
	        },
            "resources": {
               "random": {
                  "type": "OS::Heat::RandomString",
                  "properties": {
                  "length": "6"
                  }
	          }
	       }
}
JSON
}

data "flexibleengine_rts_stack_resource_v1" "resource_1" {
  stack_name = "${flexibleengine_rts_stack_v1.stack_1.name}"
  resource_name = "random"
}
`, stackName)
}
