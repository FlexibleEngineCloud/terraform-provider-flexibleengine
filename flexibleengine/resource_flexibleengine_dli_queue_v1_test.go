package flexibleengine

import (
	"fmt"
	"strings"
	"testing"

	"github.com/huaweicloud/golangsdk/openstack/dli/v1/queues"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/golangsdk"
)

func TestAccDliQueueV1_basic(t *testing.T) {
	rName := fmt.Sprintf("tf_acc_test_dli_queue_%s", acctest.RandString(5))
	resourceName := "flexibleengine_dli_queue.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDliQueueV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDliQueueV1_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDliQueueV1Exists(resourceName),
				),
			},
		},
	})
}

func testAccDliQueueV1_basic(rName string) string {
	return fmt.Sprintf(`
resource flexibleengine_dli_queue "test" {
    name             =  "%s"
    cu_count         =   16
    resource_mode    =   0
    tags             = {
        k1    =    "1"
    }
}`, rName)
}

func testAccCheckDliQueueV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.Config.DliV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating sdk client, err=%s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_dli_queue" {
			continue
		}

		res, err := fetchDliQueueV1ByQueueNameOnTest(rs.Primary.ID, client)
		if err == nil && res != nil {
			return fmt.Errorf("flexibleengine_dli_queue still exists:%s,%+v,%+v", rs.Primary.ID, err, res)
		}
	}

	return nil
}

func testAccCheckDliQueueV1Exists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		client, err := config.Config.DliV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating sdk client, err=%s", err)
		}

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Error checking flexibleengine_dli_queue.queue exist, err=not found this resource")
		}
		_, err = fetchDliQueueV1ByQueueNameOnTest(rs.Primary.ID, client)
		if err != nil {
			if strings.Contains(err.Error(), "Error finding the resource by list api") {
				return fmt.Errorf("flexibleengine_dli_queue is not exist")
			}
			return fmt.Errorf("Error checking flexibleengine_dli_queue.queue exist, err=%s", err)
		}
		return nil
	}
}

func fetchDliQueueV1ByQueueNameOnTest(primaryID string,
	client *golangsdk.ServiceClient) (interface{}, error) {
	result := queues.Get(client, primaryID)
	return result.Body, result.Err
}
