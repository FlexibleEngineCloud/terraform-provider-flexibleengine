package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccImagesImageV2DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageV2DataSource_cirros,
			},
			{
				Config: testAccImagesImageV2DataSource_basic,
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

func TestAccImagesImageV2DataSource_testQueries(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageV2DataSource_cirros,
			},
			{
				Config: testAccImagesImageV2DataSource_queryTag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID("data.flexibleengine_images_image_v2.image_1"),
				),
			},
			{
				Config: testAccImagesImageV2DataSource_querySizeMin,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID("data.flexibleengine_images_image_v2.image_1"),
				),
			},
			{
				Config: testAccImagesImageV2DataSource_querySizeMax,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID("data.flexibleengine_images_image_v2.image_1"),
				),
			},
			{
				Config: testAccImagesImageV2DataSource_cirros,
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
const testAccImagesImageV2DataSource_cirros = `
resource "flexibleengine_images_image_v2" "image_1" {
  name             = "cirrOS-tf"
  container_format = "bare"
  disk_format      = "qcow2"
  image_source_url = "http://download.cirros-cloud.net/0.3.5/cirros-0.3.5-x86_64-disk.img"
  tags             = ["cirros-tf"]
}
`

var testAccImagesImageV2DataSource_basic = fmt.Sprintf(`
%s

data "flexibleengine_images_image_v2" "image_1" {
  most_recent = true
  name        = flexibleengine_images_image_v2.image_1.name
}
`, testAccImagesImageV2DataSource_cirros)

var testAccImagesImageV2DataSource_queryTag = fmt.Sprintf(`
%s

data "flexibleengine_images_image_v2" "image_1" {
  most_recent = true
  visibility  = "private"
  tag         = "cirros-tf"
}
`, testAccImagesImageV2DataSource_cirros)

var testAccImagesImageV2DataSource_querySizeMin = fmt.Sprintf(`
%s

data "flexibleengine_images_image_v2" "image_1" {
  most_recent = true
  visibility  = "private"
  size_min    = "13000000"
}
`, testAccImagesImageV2DataSource_cirros)

var testAccImagesImageV2DataSource_querySizeMax = fmt.Sprintf(`
%s

data "flexibleengine_images_image_v2" "image_1" {
  most_recent = true
  visibility  = "private"
  size_max    = "23000000"
}
`, testAccImagesImageV2DataSource_cirros)
