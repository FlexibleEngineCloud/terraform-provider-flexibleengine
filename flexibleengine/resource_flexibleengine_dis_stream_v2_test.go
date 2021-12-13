package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/dis/v2/streams"
)

func TestAccDisStreamV2_basic(t *testing.T) {
	streamName := fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))
	resourceName := "flexibleengine_dis_stream.test"

	resource.ParallelTest(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDisStreamV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDisStreamV2_basic(streamName, streams.StreamTypeCommon, 3),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDisStreamV2Exists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", streamName),
					resource.TestCheckResourceAttr(resourceName, "type", streams.StreamTypeCommon),
					resource.TestCheckResourceAttr(resourceName, "data_duration", "24"),
					resource.TestCheckResourceAttr(resourceName, "partition_count", "3"),
					resource.TestCheckResourceAttr(resourceName, "partitions.#", "3"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"type",
				},
			},
		},
	})
}

func testAccDisStreamV2_basic(streamName string, streamType string, partitionCount int) string {
	return fmt.Sprintf(`
resource flexibleengine_dis_stream "test" {
  name			= "%s"
  type			= "%s"
  partition_count	= "%d"
}`, streamName, streamType, partitionCount)
}

func testAccCheckDisStreamV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.Config.DisV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating dis client, err=%s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_dis_stream" {
			continue
		}
		getOpts := streams.GetOpts{}
		stream, err := streams.Get(client, rs.Primary.ID, getOpts)
		if err == nil {
			if stream.Status != "Terminated" {
				return fmt.Errorf("Stream still exists")
			}
		}
	}
	return nil
}

func testAccCheckDisStreamV2Exists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		client, err := config.Config.DisV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating dis client, err=%s", err)
		}

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Error checking flexibleengine_dis_stream.%s exist, err=not found this resource", resourceName)
		}
		getOpts := streams.GetOpts{}
		_, err = streams.Get(client, rs.Primary.ID, getOpts)
		if err != nil {
			return fmt.Errorf("Error checking flexibleengine_dis_stream.%s exist, err=%s", resourceName, err)
		}
		return nil
	}
}
