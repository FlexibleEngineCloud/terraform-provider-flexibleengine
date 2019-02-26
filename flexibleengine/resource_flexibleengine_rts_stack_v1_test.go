package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/huaweicloud/golangsdk/openstack/rts/v1/stacks"
)

func TestAccRTSStackV1_basic(t *testing.T) {
	var stacks stacks.RetrievedStack
	var stackName = fmt.Sprintf("terra_test_%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRTSStackV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRTSStackV1_basic(stackName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRTSStackV1Exists("flexibleengine_rts_stack_v1.stack_1", &stacks),
					resource.TestCheckResourceAttr(
						"flexibleengine_rts_stack_v1.stack_1", "name", stackName),
					resource.TestCheckResourceAttr(
						"flexibleengine_rts_stack_v1.stack_1", "status", "CREATE_COMPLETE"),
					resource.TestCheckResourceAttr(
						"flexibleengine_rts_stack_v1.stack_1", "disable_rollback", "true"),
					resource.TestCheckResourceAttr(
						"flexibleengine_rts_stack_v1.stack_1", "timeout_mins", "60"),
				),
			},
			{
				Config: testAccRTSStackV1_update(stackName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRTSStackV1Exists("flexibleengine_rts_stack_v1.stack_1", &stacks),
					resource.TestCheckResourceAttr(
						"flexibleengine_rts_stack_v1.stack_1", "disable_rollback", "false"),
					resource.TestCheckResourceAttr(
						"flexibleengine_rts_stack_v1.stack_1", "timeout_mins", "50"),
					resource.TestCheckResourceAttr(
						"flexibleengine_rts_stack_v1.stack_1", "status", "UPDATE_COMPLETE"),
				),
			},
		},
	})
}

func TestAccRTSStackV1_timeout(t *testing.T) {
	var stacks stacks.RetrievedStack
	var stackName = fmt.Sprintf("terra_test_%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRTSStackV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRTSStackV1_timeout(stackName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRTSStackV1Exists("flexibleengine_rts_stack_v1.stack_1", &stacks),
				),
			},
		},
	})
}

func testAccCheckRTSStackV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	orchestrationClient, err := config.orchestrationV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating RTS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_rts_stack_v1" {
			continue
		}

		stack, err := stacks.Get(orchestrationClient, "terraform_provider_stack").Extract()

		if err == nil {
			if stack.Status != "DELETE_COMPLETE" {
				return fmt.Errorf("Stack still exists")
			}
		}
	}

	return nil
}

func testAccCheckRTSStackV1Exists(n string, stack *stacks.RetrievedStack) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		orchestrationClient, err := config.orchestrationV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating RTS Client : %s", err)
		}

		found, err := stacks.Get(orchestrationClient, "terraform_provider_stack").Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("stack not found")
		}

		*stack = *found

		return nil
	}
}

func testAccRTSStackV1_basic(stackName string) string {
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
  		"default": "keysclick",
        "description": "Name of existing key pair for the instance to be created."
      }
    },
    "resources": {
      "random": {
        "type": "OS::Heat::RandomString",
        "properties": {
          "length": 6
        }
      }
    }
  }
JSON

}
`, stackName)
}

func testAccRTSStackV1_update(stackName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_rts_stack_v1" "stack_1" {
  name = "%s"
  disable_rollback= false
  timeout_mins=50
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
  		"default": "keysclick",
        "description": "Name of existing key pair for the instance to be created."
      }
    },
    "resources": {
      "random": {
        "type": "OS::Heat::RandomString",
        "properties": {
          "length": 6
        }
      }
    }
  }
JSON

}
`, stackName)
}

func testAccRTSStackV1_timeout(stackName string) string {
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
  		"default": "keysclick",
        "description": "Name of existing key pair for the instance to be created."
      }
    },
    "resources": {
      "random": {
        "type": "OS::Heat::RandomString",
        "properties": {
          "length": 6
        }
      }
    }
  }
JSON

  timeouts {
    create = "10m"
    delete = "10m"
  }
}
`, stackName)
}
