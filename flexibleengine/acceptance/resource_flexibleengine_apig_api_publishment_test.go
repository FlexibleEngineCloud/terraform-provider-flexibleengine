package acceptance

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/apis"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/apig"
)

func getPublishmentResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.ApigV2Client(OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Flexibleengine APIG v2 client: %s", err)
	}
	return apig.GetVersionHistories(c, state.Primary.Attributes["instance_id"], state.Primary.Attributes["env_id"],
		state.Primary.Attributes["api_id"])
}

func TestAccApigApiPublishmentV2_basic(t *testing.T) {
	var histories []apis.ApiVersionInfo

	// The dedicated instance name only allow letters, digits and underscores (_).
	rName := acceptance.RandomAccResourceName()
	resourceName := "flexibleengine_apig_api_publishment.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&histories,
		getPublishmentResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccApigApiPublishment_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id",
						"flexibleengine_apig_instance.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "env_id",
						"flexibleengine_apig_environment.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "api_id",
						"flexibleengine_apig_api.test", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "env_name"),
					resource.TestCheckResourceAttrSet(resourceName, "publish_time"),
					resource.TestCheckResourceAttrSet(resourceName, "publish_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccApigApiPublishment_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_apig_environment" "test" {
  name        = "%s"
  instance_id = flexibleengine_apig_instance.test.id
  description = "Created by script"
}

resource "flexibleengine_apig_api_publishment" "test" {
  instance_id = flexibleengine_apig_instance.test.id
  env_id      = flexibleengine_apig_environment.test.id
  api_id      = flexibleengine_apig_api.test.id
}
`, testAccApigAPI_basic(rName), rName)
}
