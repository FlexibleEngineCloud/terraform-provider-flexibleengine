package acceptance

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/dli/v2/batches"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getSparkJobResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.DliV2Client(OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Flexibleengine DLI v2 client: %s", err)
	}
	return batches.Get(c, state.Primary.ID)
}

func TestAccDliSparkJobV2_basic(t *testing.T) {
	var job batches.CreateResp

	rName := acceptance.RandomAccResourceName()
	dashName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "flexibleengine_dli_spark_job.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&job,
		getSparkJobResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckDliSparkJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDliSparkJob_basic(rName, dashName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "queue_name",
						"flexibleengine_dli_queue.test", "name"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
		},
	})
}

func testAccCheckDliSparkJobDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	client, err := config.DliV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating Dli v2 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_dli_spark_job" {
			continue
		}

		resp, err := batches.GetState(client, rs.Primary.ID)
		// If the status of the spark job is "dead" or "success", it means that the life cycle of the job has ended.
		if err == nil && resp != nil && (resp.State != batches.StateDead && resp.State != batches.StateSuccess) {
			return fmt.Errorf("Spark job (%s) still exists.", rs.Primary.ID)
		}
	}

	return nil
}

func testAccDliSparkJob_basic(name, dashName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_dli_queue" "test" {
  name       = "%s"
  cu_count   = 16
  queue_type = "general"
}

%s

resource "flexibleengine_dli_spark_job" "test" {
  queue_name = flexibleengine_dli_queue.test.name
  name       = "%s"
  app_name   = "${flexibleengine_dli_package.test.group_name}/${flexibleengine_dli_package.test.object_name}"
  
 depends_on = [
    flexibleengine_obs_bucket.test,
    flexibleengine_obs_bucket_object.test,
  ]
}
`, name, testAccDliPackage_basic(dashName), name)
}
