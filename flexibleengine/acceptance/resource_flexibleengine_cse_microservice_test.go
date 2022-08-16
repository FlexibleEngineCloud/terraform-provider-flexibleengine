package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/cse/dedicated/v4/services"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cse"
)

func getMicroserviceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	token, err := cse.GetAuthorizationToken(state.Primary.Attributes["connect_address"],
		state.Primary.Attributes["admin_user"], state.Primary.Attributes["admin_pass"])
	if err != nil {
		return nil, err
	}

	client := common.NewCustomClient(true, state.Primary.Attributes["connect_address"], "v4", "default")
	return services.Get(client, state.Primary.ID, token)
}

func TestAccMicroservice_basic(t *testing.T) {
	var (
		service      services.Service
		randName     = acceptance.RandomAccResourceNameWithDash()
		resourceName = "flexibleengine_cse_microservice.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&service,
		getMicroserviceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccMicroservice_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "connect_address",
						"flexibleengine_cse_microservice_engine.test", "service_registry_addresses.0.public"),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "app_name", randName),
					resource.TestCheckResourceAttr(resourceName, "environment", "development"),
					resource.TestCheckResourceAttr(resourceName, "version", "1.0.1"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform test"),
					resource.TestCheckResourceAttr(resourceName, "level", "BACK"),
					resource.TestCheckResourceAttr(resourceName, "status", "UP"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccMicroserviceImportStateIdFunc(),
			},
		},
	})
}

func testAccMicroserviceImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var connAddr, username, password, microserviceId string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "flexibleengine_cse_microservice" {
				connAddr = rs.Primary.Attributes["connect_address"]
				username = rs.Primary.Attributes["admin_user"]
				password = rs.Primary.Attributes["admin_pass"]
				microserviceId = rs.Primary.ID
			}
		}
		if connAddr != "" && microserviceId != "" {
			if username != "" && password != "" {
				return fmt.Sprintf("%s/%s/%s/%s", connAddr, microserviceId, username, password), nil
			}
			return fmt.Sprintf("%s/%s", connAddr, microserviceId), nil
		}
		return "", fmt.Errorf("resource not found: %s/%s", connAddr, microserviceId)
	}
}

func testAccMicroservice_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "flexibleengine_cse_microservice" "test" {
  connect_address = flexibleengine_cse_microservice_engine.test.service_registry_addresses.0.public

  name        = "%[2]s"
  app_name    = "%[2]s"
  environment = "development"
  version     = "1.0.1"
  description = "Created by terraform test"
  level       = "BACK"

  admin_user = "root"
  admin_pass = "AccTest!123"
}
`, testAccMicroserviceEngine_basic(rName), rName)
}
