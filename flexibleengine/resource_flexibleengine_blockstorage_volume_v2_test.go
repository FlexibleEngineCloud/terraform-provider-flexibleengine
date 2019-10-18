package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/blockstorage/v2/volumes"
)

func TestAccBlockStorageV2Volume_basic(t *testing.T) {
	var volume volumes.Volume

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBlockStorageV2VolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageV2Volume_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageV2VolumeExists("flexibleengine_blockstorage_volume_v2.volume_1", &volume),
					testAccCheckBlockStorageV2VolumeMetadata(&volume, "foo", "bar"),
					resource.TestCheckResourceAttr(
						"flexibleengine_blockstorage_volume_v2.volume_1", "name", "volume_1"),
					resource.TestCheckResourceAttr(
						"flexibleengine_blockstorage_volume_v2.volume_1", "size", "1"),
					resource.TestCheckResourceAttr(
						"flexibleengine_blockstorage_volume_v2.volume_1", "multiattach", "true"),
				),
			},
			{
				Config: testAccBlockStorageV2Volume_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageV2VolumeExists("flexibleengine_blockstorage_volume_v2.volume_1", &volume),
					testAccCheckBlockStorageV2VolumeMetadata(&volume, "foo", "bar"),
					resource.TestCheckResourceAttr(
						"flexibleengine_blockstorage_volume_v2.volume_1", "name", "volume_1-updated"),
					resource.TestCheckResourceAttr(
						"flexibleengine_blockstorage_volume_v2.volume_1", "size", "2"),
				),
			},
		},
	})
}

func TestAccBlockStorageV2Volume_online_resize(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBlockStorageV2VolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageV2Volume_online_resize,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"flexibleengine_blockstorage_volume_v2.volume_1", "size", "1"),
				),
			},
			{
				Config: testAccBlockStorageV2Volume_online_resize_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"flexibleengine_blockstorage_volume_v2.volume_1", "size", "2"),
				),
			},
		},
	})
}

func TestAccBlockStorageV2Volume_image(t *testing.T) {
	var volume volumes.Volume

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBlockStorageV2VolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageV2Volume_image,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageV2VolumeExists("flexibleengine_blockstorage_volume_v2.volume_1", &volume),
					resource.TestCheckResourceAttr(
						"flexibleengine_blockstorage_volume_v2.volume_1", "name", "volume_1"),
				),
			},
		},
	})
}

func TestAccBlockStorageV2Volume_timeout(t *testing.T) {
	var volume volumes.Volume

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBlockStorageV2VolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageV2Volume_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageV2VolumeExists("flexibleengine_blockstorage_volume_v2.volume_1", &volume),
				),
			},
		},
	})
}

func testAccCheckBlockStorageV2VolumeDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	blockStorageClient, err := config.blockStorageV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine block storage client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_blockstorage_volume_v2" {
			continue
		}

		_, err := volumes.Get(blockStorageClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Volume still exists")
		}
	}

	return nil
}

func testAccCheckBlockStorageV2VolumeExists(n string, volume *volumes.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		blockStorageClient, err := config.blockStorageV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine block storage client: %s", err)
		}

		found, err := volumes.Get(blockStorageClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Volume not found")
		}

		*volume = *found

		return nil
	}
}

func testAccCheckBlockStorageV2VolumeDoesNotExist(t *testing.T, n string, volume *volumes.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		blockStorageClient, err := config.blockStorageV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine block storage client: %s", err)
		}

		_, err = volumes.Get(blockStorageClient, volume.ID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil
			}
			return err
		}

		return fmt.Errorf("Volume still exists")
	}
}

func testAccCheckBlockStorageV2VolumeMetadata(
	volume *volumes.Volume, k string, v string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if volume.Metadata == nil {
			return fmt.Errorf("No metadata")
		}

		for key, value := range volume.Metadata {
			if k != key {
				continue
			}

			if v == value {
				return nil
			}

			return fmt.Errorf("Bad value for %s: %s", k, value)
		}

		return fmt.Errorf("Metadata not found: %s", k)
	}
}

const testAccBlockStorageV2Volume_basic = `
resource "flexibleengine_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  description = "first test volume"
  metadata = {
    foo = "bar"
  }
  size = 1
  multiattach = true
}
`

const testAccBlockStorageV2Volume_update = `
resource "flexibleengine_blockstorage_volume_v2" "volume_1" {
  name = "volume_1-updated"
  description = "first test volume"
  metadata = {
    foo = "bar"
  }
  size = 2
}
`

var testAccBlockStorageV2Volume_online_resize = fmt.Sprintf(`
resource "flexibleengine_compute_instance_v2" "basic" {
  name            = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
}

resource "flexibleengine_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  description = "test volume"
  size = 1
}

resource "flexibleengine_compute_volume_attach_v2" "va_1" {
  instance_id = "${flexibleengine_compute_instance_v2.basic.id}"
  volume_id   = "${flexibleengine_blockstorage_volume_v2.volume_1.id}"
}
`, OS_NETWORK_ID)

var testAccBlockStorageV2Volume_online_resize_update = fmt.Sprintf(`
resource "flexibleengine_compute_instance_v2" "basic" {
  name            = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
}

resource "flexibleengine_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  description = "test volume"
  size = 2
}

resource "flexibleengine_compute_volume_attach_v2" "va_1" {
  instance_id = "${flexibleengine_compute_instance_v2.basic.id}"
  volume_id   = "${flexibleengine_blockstorage_volume_v2.volume_1.id}"
}
`, OS_NETWORK_ID)

var testAccBlockStorageV2Volume_image = fmt.Sprintf(`
resource "flexibleengine_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  size = 50
  image_id = "%s"
}
`, OS_IMAGE_ID)

const testAccBlockStorageV2Volume_timeout = `
resource "flexibleengine_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  description = "first test volume"
  size = 1

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`
