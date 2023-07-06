package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/imageservice/v2/images"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccImagesImageV2_basic(t *testing.T) {
	var image images.Image
	resourceName := "flexibleengine_images_image_v2.image_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeprecated(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckImagesImageV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageV2Exists(resourceName, &image),
					resource.TestCheckResourceAttr(resourceName, "name", "rancheros-test"),
					resource.TestCheckResourceAttr(resourceName, "container_format", "bare"),
					resource.TestCheckResourceAttr(resourceName, "schema", "/v2/schemas/image"),
				),
			},
			{
				Config: testAccImagesImageV2_update_name,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "rancheros-openstack"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"local_file_path",
					"image_cache_path",
					"image_source_url",
				},
			},
		},
	})
}

func TestAccImagesImageV2_tags(t *testing.T) {
	var image images.Image

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeprecated(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckImagesImageV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageV2_tags_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageV2Exists("flexibleengine_images_image_v2.image_1", &image),
					testAccCheckImagesImageV2HasTag("flexibleengine_images_image_v2.image_1", "foo"),
					testAccCheckImagesImageV2HasTag("flexibleengine_images_image_v2.image_1", "bar"),
					testAccCheckImagesImageV2TagCount("flexibleengine_images_image_v2.image_1", 2),
				),
			},
			{
				Config: testAccImagesImageV2_tags_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageV2Exists("flexibleengine_images_image_v2.image_1", &image),
					testAccCheckImagesImageV2HasTag("flexibleengine_images_image_v2.image_1", "foo"),
					testAccCheckImagesImageV2HasTag("flexibleengine_images_image_v2.image_1", "bar"),
					testAccCheckImagesImageV2HasTag("flexibleengine_images_image_v2.image_1", "baz"),
					testAccCheckImagesImageV2TagCount("flexibleengine_images_image_v2.image_1", 3),
				),
			},
			{
				Config: testAccImagesImageV2_tags_3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageV2Exists("flexibleengine_images_image_v2.image_1", &image),
					testAccCheckImagesImageV2HasTag("flexibleengine_images_image_v2.image_1", "foo"),
					testAccCheckImagesImageV2HasTag("flexibleengine_images_image_v2.image_1", "baz"),
					testAccCheckImagesImageV2TagCount("flexibleengine_images_image_v2.image_1", 2),
				),
			},
		},
	})
}

func TestAccImagesImageV2_visibility(t *testing.T) {
	var image images.Image

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckDeprecated(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckImagesImageV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageV2_visibility_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageV2Exists("flexibleengine_images_image_v2.image_1", &image),
					resource.TestCheckResourceAttr(
						"flexibleengine_images_image_v2.image_1", "visibility", "private"),
				),
			},
			{
				Config: testAccImagesImageV2_visibility_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageV2Exists("flexibleengine_images_image_v2.image_1", &image),
					resource.TestCheckResourceAttr(
						"flexibleengine_images_image_v2.image_1", "visibility", "public"),
				),
			},
		},
	})
}

func testAccCheckImagesImageV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	imageClient, err := config.ImageV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine image client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_images_image_v2" {
			continue
		}

		_, err := images.Get(imageClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Image still exists")
		}
	}

	return nil
}

func testAccCheckImagesImageV2Exists(n string, image *images.Image) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		imageClient, err := config.ImageV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine image client: %s", err)
		}

		found, err := images.Get(imageClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Image not found")
		}

		*image = *found

		return nil
	}
}

func testAccCheckImagesImageV2HasTag(n, tag string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		imageClient, err := config.ImageV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine image client: %s", err)
		}

		found, err := images.Get(imageClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Image not found")
		}

		for _, v := range found.Tags {
			if tag == v {
				return nil
			}
		}

		return fmt.Errorf("Tag not found: %s", tag)
	}
}

func testAccCheckImagesImageV2TagCount(n string, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		imageClient, err := config.ImageV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine image client: %s", err)
		}

		found, err := images.Get(imageClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Image not found")
		}

		if len(found.Tags) != expected {
			return fmt.Errorf("Expecting %d tags, found %d", expected, len(found.Tags))
		}

		return nil
	}
}

var testAccImagesImageV2_basic = `
resource "flexibleengine_images_image_v2" "image_1" {
  name             = "rancheros-test"
  image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
  container_format = "bare"
  disk_format      = "qcow2"
}`

var testAccImagesImageV2_update_name = `
resource "flexibleengine_images_image_v2" "image_1" {
  name             = "rancheros-openstack"
  image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
  container_format = "bare"
  disk_format      = "qcow2"
}`

var testAccImagesImageV2_tags_1 = `
resource "flexibleengine_images_image_v2" "image_1" {
  name             = "cirrOS-tf"
  container_format = "bare"
  disk_format      = "qcow2"
  image_source_url = "http://download.cirros-cloud.net/0.3.5/cirros-0.3.5-x86_64-disk.img"
  tags             = ["foo","bar"]
}`

var testAccImagesImageV2_tags_2 = `
resource "flexibleengine_images_image_v2" "image_1" {
  name             = "cirrOS-tf"
  container_format = "bare"
  disk_format      = "qcow2"
  image_source_url = "http://download.cirros-cloud.net/0.3.5/cirros-0.3.5-x86_64-disk.img"
  tags             = ["foo","bar","baz"]
}`

var testAccImagesImageV2_tags_3 = `
resource "flexibleengine_images_image_v2" "image_1" {
  name             = "cirrOS-tf"
  container_format = "bare"
  disk_format      = "qcow2"
  image_source_url = "http://download.cirros-cloud.net/0.3.5/cirros-0.3.5-x86_64-disk.img"
  tags             = ["foo","baz"]
}`

var testAccImagesImageV2_visibility_1 = `
resource "flexibleengine_images_image_v2" "image_1" {
  name             = "cirrOS-tf"
  container_format = "bare"
  disk_format      = "qcow2"
  image_source_url = "http://download.cirros-cloud.net/0.3.5/cirros-0.3.5-x86_64-disk.img"
  visibility       = "private"
}`

var testAccImagesImageV2_visibility_2 = `
resource "flexibleengine_images_image_v2" "image_1" {
  name             = "cirrOS-tf"
  container_format = "bare"
  disk_format      = "qcow2"
  image_source_url = "http://download.cirros-cloud.net/0.3.5/cirros-0.3.5-x86_64-disk.img"
  visibility       = "public"
}`
