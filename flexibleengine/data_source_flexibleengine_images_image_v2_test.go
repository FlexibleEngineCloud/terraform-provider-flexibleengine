package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// PASS
func TestAccOrangeCloudImagesV2ImageDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOrangeCloudImagesV2ImageDataSource_cirros,
			},
			resource.TestStep{
				Config: testAccOrangeCloudImagesV2ImageDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID("data.flexibleengine_images_image_v2.image_1"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_images_image_v2.image_1", "name", "CirrOS-tf"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_images_image_v2.image_1", "container_format", "bare"),
					/*resource.TestCheckResourceAttr(
					"data.flexibleengine_images_image_v2.image_1", "disk_format", "qcow2"), */
					/*resource.TestCheckResourceAttr(
					"data.flexibleengine_images_image_v2.image_1", "min_disk_gb", "0"), */
					resource.TestCheckResourceAttr(
						"data.flexibleengine_images_image_v2.image_1", "min_ram_mb", "0"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_images_image_v2.image_1", "protected", "false"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_images_image_v2.image_1", "visibility", "private"),
				),
			},
		},
	})
}

// PASS
func TestAccOrangeCloudImagesV2ImageDataSource_testQueries(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOrangeCloudImagesV2ImageDataSource_cirros,
			},
			resource.TestStep{
				Config: testAccOrangeCloudImagesV2ImageDataSource_queryTag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID("data.flexibleengine_images_image_v2.image_1"),
				),
			},
			resource.TestStep{
				Config: testAccOrangeCloudImagesV2ImageDataSource_querySizeMin,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID("data.flexibleengine_images_image_v2.image_1"),
				),
			},
			resource.TestStep{
				Config: testAccOrangeCloudImagesV2ImageDataSource_querySizeMax,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID("data.flexibleengine_images_image_v2.image_1"),
				),
			},
			resource.TestStep{
				Config: testAccOrangeCloudImagesV2ImageDataSource_cirros,
			},
		},
	})
}

func testAccCheckImagesV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find image data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Image data source ID not set")
		}

		return nil
	}
}

// Standard CirrOS image
const testAccOrangeCloudImagesV2ImageDataSource_cirros = `
resource "flexibleengine_images_image_v2" "image_1" {
	name = "CirrOS-tf"
	container_format = "bare"
	disk_format = "qcow2"
	image_source_url = "http://download.cirros-cloud.net/0.3.5/cirros-0.3.5-x86_64-disk.img"
	tags = ["cirros-tf"]
}
`

var testAccOrangeCloudImagesV2ImageDataSource_basic = fmt.Sprintf(`
%s

data "flexibleengine_images_image_v2" "image_1" {
	most_recent = true
	name = "${flexibleengine_images_image_v2.image_1.name}"
}
`, testAccOrangeCloudImagesV2ImageDataSource_cirros)

var testAccOrangeCloudImagesV2ImageDataSource_queryTag = fmt.Sprintf(`
%s

data "flexibleengine_images_image_v2" "image_1" {
	most_recent = true
	visibility = "private"
	tag = "cirros-tf"
}
`, testAccOrangeCloudImagesV2ImageDataSource_cirros)

var testAccOrangeCloudImagesV2ImageDataSource_querySizeMin = fmt.Sprintf(`
%s

data "flexibleengine_images_image_v2" "image_1" {
	most_recent = true
	visibility = "private"
	size_min = "13000000"
}
`, testAccOrangeCloudImagesV2ImageDataSource_cirros)

var testAccOrangeCloudImagesV2ImageDataSource_querySizeMax = fmt.Sprintf(`
%s

data "flexibleengine_images_image_v2" "image_1" {
	most_recent = true
	visibility = "private"
	size_max = "23000000"
}
`, testAccOrangeCloudImagesV2ImageDataSource_cirros)
