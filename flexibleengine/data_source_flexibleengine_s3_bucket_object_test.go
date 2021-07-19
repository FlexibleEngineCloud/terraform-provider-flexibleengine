package flexibleengine

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceS3BucketObject_basic(t *testing.T) {
	rInt := acctest.RandInt()
	resourceOnlyConf, conf := testAccDataSourceS3ObjectConfig_basic(rInt)

	var rObj s3.GetObjectOutput
	var dsObj s3.GetObjectOutput

	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: resourceOnlyConf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketObjectExists("flexibleengine_s3_bucket_object.object", &rObj),
				),
			},
			{
				Config: conf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3ObjectDataSourceExists("data.flexibleengine_s3_bucket_object.obj", &dsObj),
					resource.TestCheckResourceAttr("data.flexibleengine_s3_bucket_object.obj", "content_length", "11"),
					resource.TestCheckResourceAttr("data.flexibleengine_s3_bucket_object.obj", "content_type", "binary/octet-stream"),
					resource.TestCheckResourceAttr("data.flexibleengine_s3_bucket_object.obj", "etag", "b10a8db164e0754105b7a99be72e3fe5"),
					resource.TestMatchResourceAttr("data.flexibleengine_s3_bucket_object.obj", "last_modified",
						regexp.MustCompile("^[a-zA-Z]{3}, [0-9]+ [a-zA-Z]+ [0-9]{4} [0-9:]+ [A-Z]+$")),
					resource.TestCheckNoResourceAttr("data.flexibleengine_s3_bucket_object.obj", "body"),
				),
			},
		},
	})
}

func TestAccDataSourceS3BucketObject_readableBody(t *testing.T) {
	rInt := acctest.RandInt()
	resourceOnlyConf, conf := testAccDataSourceS3ObjectConfig_readableBody(rInt)

	var rObj s3.GetObjectOutput
	var dsObj s3.GetObjectOutput

	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: resourceOnlyConf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketObjectExists("flexibleengine_s3_bucket_object.object", &rObj),
				),
			},
			{
				Config: conf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3ObjectDataSourceExists("data.flexibleengine_s3_bucket_object.obj", &dsObj),
					resource.TestCheckResourceAttr("data.flexibleengine_s3_bucket_object.obj", "content_length", "3"),
					resource.TestCheckResourceAttr("data.flexibleengine_s3_bucket_object.obj", "content_type", "text/plain"),
					resource.TestCheckResourceAttr("data.flexibleengine_s3_bucket_object.obj", "etag", "a6105c0a611b41b08f1209506350279e"),
					resource.TestMatchResourceAttr("data.flexibleengine_s3_bucket_object.obj", "last_modified",
						regexp.MustCompile("^[a-zA-Z]{3}, [0-9]+ [a-zA-Z]+ [0-9]{4} [0-9:]+ [A-Z]+$")),
					resource.TestCheckResourceAttr("data.flexibleengine_s3_bucket_object.obj", "body", "yes"),
				),
			},
		},
	})
}

func TestAccDataSourceS3BucketObject_allParams(t *testing.T) {
	rInt := acctest.RandInt()
	resourceOnlyConf, conf := testAccDataSourceS3ObjectConfig_allParams(rInt)

	var rObj s3.GetObjectOutput
	var dsObj s3.GetObjectOutput

	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: resourceOnlyConf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketObjectExists("flexibleengine_s3_bucket_object.object", &rObj),
				),
			},
			{
				Config: conf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3ObjectDataSourceExists("data.flexibleengine_s3_bucket_object.obj", &dsObj),
					resource.TestCheckResourceAttr("data.flexibleengine_s3_bucket_object.obj", "content_length", "21"),
					resource.TestCheckResourceAttr("data.flexibleengine_s3_bucket_object.obj", "content_type", "application/unknown"),
					resource.TestCheckResourceAttr("data.flexibleengine_s3_bucket_object.obj", "etag", "723f7a6ac0c57b445790914668f98640"),
					resource.TestMatchResourceAttr("data.flexibleengine_s3_bucket_object.obj", "last_modified",
						regexp.MustCompile("^[a-zA-Z]{3}, [0-9]+ [a-zA-Z]+ [0-9]{4} [0-9:]+ [A-Z]+$")),
					resource.TestCheckNoResourceAttr("data.flexibleengine_s3_bucket_object.obj", "body"),
					resource.TestCheckResourceAttr("data.flexibleengine_s3_bucket_object.obj", "cache_control", "no-cache"),
					resource.TestCheckResourceAttr("data.flexibleengine_s3_bucket_object.obj", "content_disposition", "attachment"),
					resource.TestCheckResourceAttr("data.flexibleengine_s3_bucket_object.obj", "content_encoding", "identity"),
					resource.TestCheckResourceAttr("data.flexibleengine_s3_bucket_object.obj", "content_language", "en-GB"),
					// Encryption is off
					resource.TestCheckResourceAttr("data.flexibleengine_s3_bucket_object.obj", "server_side_encryption", ""),
					resource.TestCheckResourceAttr("data.flexibleengine_s3_bucket_object.obj", "sse_kms_key_id", ""),
					resource.TestCheckResourceAttr("data.flexibleengine_s3_bucket_object.obj", "expiration", ""),
					// Currently unsupported in flexibleengine_s3_bucket_object resource
					resource.TestCheckResourceAttr("data.flexibleengine_s3_bucket_object.obj", "expires", ""),
					resource.TestCheckResourceAttr("data.flexibleengine_s3_bucket_object.obj", "website_redirect_location", ""),
					resource.TestCheckResourceAttr("data.flexibleengine_s3_bucket_object.obj", "metadata.%", "0"),
				),
			},
		},
	})
}

func testAccCheckS3ObjectDataSourceExists(n string, obj *s3.GetObjectOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find S3 object data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("S3 object data source ID not set")
		}

		config := testAccProvider.Meta().(*Config)
		s3conn, err := config.computeS3conn(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine s3 client: %s", err)
		}
		out, err := s3conn.GetObject(
			&s3.GetObjectInput{
				Bucket: aws.String(rs.Primary.Attributes["bucket"]),
				Key:    aws.String(rs.Primary.Attributes["key"]),
			})
		if err != nil {
			return fmt.Errorf("Failed getting S3 Object from %s: %s",
				rs.Primary.Attributes["bucket"]+"/"+rs.Primary.Attributes["key"], err)
		}

		*obj = *out

		return nil
	}
}

func testAccDataSourceS3ObjectConfig_basic(randInt int) (string, string) {
	resources := fmt.Sprintf(`
resource "flexibleengine_s3_bucket" "object_bucket" {
	bucket = "tf-object-test-bucket-%d"
}
resource "flexibleengine_s3_bucket_object" "object" {
	bucket = "${flexibleengine_s3_bucket.object_bucket.bucket}"
	key = "tf-testing-obj-%d"
	content = "Hello World"
}
`, randInt, randInt)

	both := fmt.Sprintf(`%s
data "flexibleengine_s3_bucket_object" "obj" {
	bucket = "tf-object-test-bucket-%d"
	key = "tf-testing-obj-%d"
}`, resources, randInt, randInt)

	return resources, both
}

func testAccDataSourceS3ObjectConfig_readableBody(randInt int) (string, string) {
	resources := fmt.Sprintf(`
resource "flexibleengine_s3_bucket" "object_bucket" {
	bucket = "tf-object-test-bucket-%d"
}
resource "flexibleengine_s3_bucket_object" "object" {
	bucket = "${flexibleengine_s3_bucket.object_bucket.bucket}"
	key = "tf-testing-obj-%d-readable"
	content = "yes"
	content_type = "text/plain"
}
`, randInt, randInt)

	both := fmt.Sprintf(`%s
data "flexibleengine_s3_bucket_object" "obj" {
	bucket = "tf-object-test-bucket-%d"
	key = "tf-testing-obj-%d-readable"
}`, resources, randInt, randInt)

	return resources, both
}

func testAccDataSourceS3ObjectConfig_allParams(randInt int) (string, string) {
	resources := fmt.Sprintf(`
resource "flexibleengine_s3_bucket" "object_bucket" {
	bucket = "tf-object-test-bucket-%d"
	versioning {
		enabled = true
	}
}

resource "flexibleengine_s3_bucket_object" "object" {
	bucket = "${flexibleengine_s3_bucket.object_bucket.bucket}"
	key = "tf-testing-obj-%d-all-params"
	content = <<CONTENT
{"msg": "Hi there!"}
CONTENT
	content_type = "application/unknown"
	cache_control = "no-cache"
	content_disposition = "attachment"
	content_encoding = "identity"
	content_language = "en-GB"
	#tags {
	#	Key1 = "Value 1"
	#}
}
`, randInt, randInt)

	both := fmt.Sprintf(`%s
data "flexibleengine_s3_bucket_object" "obj" {
	bucket = "tf-object-test-bucket-%d"
	key = "tf-testing-obj-%d-all-params"
}`, resources, randInt, randInt)

	return resources, both
}
