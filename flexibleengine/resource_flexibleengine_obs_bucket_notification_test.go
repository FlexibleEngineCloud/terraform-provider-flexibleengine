package flexibleengine

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccObsBucket_notification(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "flexibleengine_obs_bucket.bucket"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckS3(t)
			testAccPreCheckOBSNotification(t)
		},
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckObsBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObsBucketConfigWithNotification(rInt, OS_OBS_URN_SMN),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketExists(resourceName),
					resource.TestCheckResourceAttr(
						resourceName, "topic_configurations.0.name", "001"),
					resource.TestCheckResourceAttr(
						resourceName, "topic_configurations.0.topic_urn", OS_OBS_URN_SMN),
					resource.TestCheckResourceAttr(
						resourceName, "topic_configurations.0.events.0", "ObjectCreated:*"),
					resource.TestCheckResourceAttr(
						resourceName, "topic_configurations.0.prefix", "tf"),
					resource.TestCheckResourceAttr(
						resourceName, "topic_configurations.0.suffix", ".jpg"),
					resource.TestCheckResourceAttr(
						resourceName, "topic_configurations.1.name", "002"),
					resource.TestCheckResourceAttr(
						resourceName, "topic_configurations.1.topic_urn", OS_OBS_URN_SMN),
					resource.TestCheckResourceAttr(
						resourceName, "topic_configurations.1.events.0", "ObjectCreated:Put"),
					resource.TestCheckResourceAttr(
						resourceName, "topic_configurations.1.prefix", "demo"),
					resource.TestCheckResourceAttr(
						resourceName, "topic_configurations.1.suffix", ".png"),
				),
			},
			{
				Config: testAccObsBucketConfigWithNotificationUpdate(rInt, OS_OBS_URN_SMN),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketExists(resourceName),
					resource.TestCheckResourceAttr(
						resourceName, "topic_configurations.0.name", "003"),
					resource.TestCheckResourceAttr(
						resourceName, "topic_configurations.0.topic_urn", OS_OBS_URN_SMN),
					resource.TestCheckResourceAttr(
						resourceName, "topic_configurations.0.events.0", "ObjectRemoved:*"),
					resource.TestCheckResourceAttr(
						resourceName, "topic_configurations.0.prefix", "tf_update"),
					resource.TestCheckResourceAttr(
						resourceName, "topic_configurations.0.suffix", ".png"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccObsBucketNotificationImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccObsBucketNotificationImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmtp.Errorf("resource not found")
		}
		bucket := rs.Primary.ID
		name := rs.Primary.Attributes["name"]
		return fmt.Sprintf("%s/%s", bucket, name), nil
	}
}

func testAccObsBucketConfigWithNotification(randInt int, urnSmn string) string {
	return fmt.Sprintf(`
resource "flexibleengine_obs_bucket" "bucket" {
    bucket = "tf-test-bucket-%d"
    acl = "public-read"
}

resource "flexibleengine_obs_bucket_notification" notification {
	bucket = flexibleengine_obs_bucket.bucket.bucket
	topic_configurations {
		name      = "001"
		events    = ["ObjectCreated:*"]
		prefix    = "tf"
		suffix    = ".jpg"
		topic_urn = "%s"
	}

	topic_configurations {
		name      = "002"
		events    = ["ObjectCreated:Put"]
		prefix    = "demo"
		suffix    = ".png"
		topic_urn = "%s"
	}
}

`, randInt, urnSmn, urnSmn)
}

func testAccObsBucketConfigWithNotificationUpdate(randInt int, urnSmn string) string {
	return fmt.Sprintf(`
resource "flexibleengine_obs_bucket" "bucket" {
    bucket = "tf-test-bucket-%d"
    acl = "public-read"
}

resource "flexibleengine_obs_bucket_notification" notification {
	bucket = flexibleengine_obs_bucket.bucket.bucket
	topic_configurations {
		name      = "003"
		events    = ["ObjectRemoved:*"]
		prefix    = "tf_update"
		suffix    = ".png"
		topic_urn = "%s"
	}
}

`, randInt, urnSmn)
}
