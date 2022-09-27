package acceptance

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/environments"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/apig"
)

func TestAccApigEnvironmentV2_basic(t *testing.T) {
	var (
		// Only letters, digits and underscores (_) are allowed in the environment name and dedicated instance name.
		rName        = fmt.Sprintf("tf_acc_test_%s", acctest.RandString(5))
		resourceName = "flexibleengine_apig_environment.test"
		env          environments.Environment
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckApigEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApigEnvironment_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigEnvironmentExists(resourceName, &env),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by script"),
				),
			},
			{
				Config: testAccApigEnvironment_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigEnvironmentExists(resourceName, &env),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"_update"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated by script"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccApigSubResNameImportStateFunc(resourceName),
			},
		},
	})
}

func testAccCheckApigEnvironmentDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	client, err := config.ApigV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating Flexibleengine APIG v2 client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_apig_environment" {
			continue
		}
		_, err := apig.GetEnvironmentFormServer(client, rs.Primary.Attributes["instance_id"], rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("APIG v2 API environment (%s) is still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckApigEnvironmentExists(name string, env *environments.Environment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Resource %s not found", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No APIG V2 API group Id")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		client, err := config.ApigV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating Flexibleengine APIG v2 client: %s", err)
		}
		found, err := apig.GetEnvironmentFormServer(client, rs.Primary.Attributes["instance_id"], rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error getting APIG v2 API environment (%s): %s", rs.Primary.ID, err)
		}
		*env = *found
		return nil
	}
}

func testAccApigEnvironment_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_apig_environment" "test" {
  name        = "%s"
  instance_id = flexibleengine_apig_instance.test.id
  description = "Created by script"
}
`, testAccApigApplication_base(rName), rName)
}

func testAccApigEnvironment_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_apig_environment" "test" {
  name        = "%s_update"
  instance_id = flexibleengine_apig_instance.test.id
  description = "Updated by script"
}
`, testAccApigApplication_base(rName), rName)
}
