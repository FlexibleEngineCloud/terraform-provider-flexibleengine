package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/compute/v2/extensions/volumeattach"
	"github.com/huaweicloud/golangsdk/openstack/compute/v2/servers"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/security/groups"
	"github.com/huaweicloud/golangsdk/pagination"
)

func TestAccComputeV2Instance_basic(t *testing.T) {
	var instance servers.Server
	resourceName := "flexibleengine_compute_instance_v2.instance_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(resourceName, &instance),
					testAccCheckComputeV2InstanceMetadata(&instance, "foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", OS_AVAILABILITY_ZONE),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"stop_before_destroy",
				},
			},
		},
	})
}

func TestAccComputeV2Instance_secgroupMulti(t *testing.T) {
	var instance_1 servers.Server
	var secgroup_1 groups.SecGroup

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_secgroupMulti,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SecGroupExists(
						"flexibleengine_networking_secgroup_v2.secgroup_1", &secgroup_1),
					testAccCheckComputeV2InstanceExists(
						"flexibleengine_compute_instance_v2.instance_1", &instance_1),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_secgroupMultiUpdate(t *testing.T) {
	var instance_1 servers.Server
	var secgroup_1, secgroup_2 groups.SecGroup

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_secgroupMultiUpdate_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SecGroupExists(
						"flexibleengine_networking_secgroup_v2.secgroup_1", &secgroup_1),
					testAccCheckNetworkingV2SecGroupExists(
						"flexibleengine_networking_secgroup_v2.secgroup_2", &secgroup_2),
					testAccCheckComputeV2InstanceExists(
						"flexibleengine_compute_instance_v2.instance_1", &instance_1),
				),
			},
			{
				Config: testAccComputeV2Instance_secgroupMultiUpdate_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SecGroupExists(
						"flexibleengine_networking_secgroup_v2.secgroup_1", &secgroup_1),
					testAccCheckNetworkingV2SecGroupExists(
						"flexibleengine_networking_secgroup_v2.secgroup_2", &secgroup_2),
					testAccCheckComputeV2InstanceExists(
						"flexibleengine_compute_instance_v2.instance_1", &instance_1),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_bootFromVolumeImage(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_bootFromVolumeImage,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("flexibleengine_compute_instance_v2.instance_1", &instance),
					testAccCheckComputeV2InstanceBootVolumeAttachment(&instance),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_bootFromVolumeVolume(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_bootFromVolumeVolume,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("flexibleengine_compute_instance_v2.instance_1", &instance),
					testAccCheckComputeV2InstanceBootVolumeAttachment(&instance),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_bootFromVolumeForceNew(t *testing.T) {
	var instance1_1 servers.Server
	var instance1_2 servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_bootFromVolumeForceNew_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(
						"flexibleengine_compute_instance_v2.instance_1", &instance1_1),
				),
			},
			{
				Config: testAccComputeV2Instance_bootFromVolumeForceNew_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(
						"flexibleengine_compute_instance_v2.instance_1", &instance1_2),
					testAccCheckComputeV2InstanceInstanceIDsDoNotMatch(&instance1_1, &instance1_2),
				),
			},
		},
	})
}

// TODO: verify the personality really exists on the instance.
func TestAccComputeV2Instance_personality(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_personality,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("flexibleengine_compute_instance_v2.instance_1", &instance),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_changeFixedIP(t *testing.T) {
	var instance1_1 servers.Server
	var instance1_2 servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_changeFixedIP_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(
						"flexibleengine_compute_instance_v2.instance_1", &instance1_1),
				),
			},
			{
				Config: testAccComputeV2Instance_changeFixedIP_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(
						"flexibleengine_compute_instance_v2.instance_1", &instance1_2),
					testAccCheckComputeV2InstanceInstanceIDsDoNotMatch(&instance1_1, &instance1_2),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_stopBeforeDestroy(t *testing.T) {
	var instance servers.Server
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_stopBeforeDestroy,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("flexibleengine_compute_instance_v2.instance_1", &instance),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_metadataRemove(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_metadataRemove_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("flexibleengine_compute_instance_v2.instance_1", &instance),
					testAccCheckComputeV2InstanceMetadata(&instance, "foo", "bar"),
					testAccCheckComputeV2InstanceMetadata(&instance, "abc", "def"),
					resource.TestCheckResourceAttr(
						"flexibleengine_compute_instance_v2.instance_1", "all_metadata.foo", "bar"),
					resource.TestCheckResourceAttr(
						"flexibleengine_compute_instance_v2.instance_1", "all_metadata.abc", "def"),
				),
			},
			{
				Config: testAccComputeV2Instance_metadataRemove_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("flexibleengine_compute_instance_v2.instance_1", &instance),
					testAccCheckComputeV2InstanceMetadata(&instance, "foo", "bar"),
					testAccCheckComputeV2InstanceMetadata(&instance, "ghi", "jkl"),
					testAccCheckComputeV2InstanceNoMetadataKey(&instance, "abc"),
					resource.TestCheckResourceAttr(
						"flexibleengine_compute_instance_v2.instance_1", "all_metadata.foo", "bar"),
					resource.TestCheckResourceAttr(
						"flexibleengine_compute_instance_v2.instance_1", "all_metadata.ghi", "jkl"),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_timeout(t *testing.T) {
	var instance servers.Server
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("flexibleengine_compute_instance_v2.instance_1", &instance),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_auto_recovery(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("flexibleengine_compute_instance_v2.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"flexibleengine_compute_instance_v2.instance_1", "auto_recovery", "false"),
				),
			},
			{
				Config: testAccComputeV2Instance_auto_recovery,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("flexibleengine_compute_instance_v2.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"flexibleengine_compute_instance_v2.instance_1", "auto_recovery", "true"),
				),
			},
		},
	})
}

func testAccCheckComputeV2InstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	computeClient, err := config.computeV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_compute_instance_v2" {
			continue
		}

		server, err := servers.Get(computeClient, rs.Primary.ID).Extract()
		if err == nil {
			if server.Status != "SOFT_DELETED" {
				return fmt.Errorf("Instance still exists")
			}
		}
	}

	return nil
}

func testAccCheckComputeV2InstanceExists(n string, instance *servers.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		computeClient, err := config.computeV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine compute client: %s", err)
		}

		found, err := servers.Get(computeClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Instance not found")
		}

		*instance = *found

		return nil
	}
}

func testAccCheckComputeV2InstanceDoesNotExist(n string, instance *servers.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		computeClient, err := config.computeV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine compute client: %s", err)
		}

		_, err = servers.Get(computeClient, instance.ID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil
			}
			return err
		}

		return fmt.Errorf("Instance still exists")
	}
}

func testAccCheckComputeV2InstanceMetadata(
	instance *servers.Server, k string, v string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instance.Metadata == nil {
			return fmt.Errorf("No metadata")
		}

		for key, value := range instance.Metadata {
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

func testAccCheckComputeV2InstanceNoMetadataKey(
	instance *servers.Server, k string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instance.Metadata == nil {
			return nil
		}

		for key := range instance.Metadata {
			if k == key {
				return fmt.Errorf("Metadata found: %s", k)
			}
		}

		return nil
	}
}

func testAccCheckComputeV2InstanceBootVolumeAttachment(
	instance *servers.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var attachments []volumeattach.VolumeAttachment

		config := testAccProvider.Meta().(*Config)
		computeClient, err := config.computeV2Client(OS_REGION_NAME)
		if err != nil {
			return err
		}

		err = volumeattach.List(computeClient, instance.ID).EachPage(
			func(page pagination.Page) (bool, error) {

				actual, err := volumeattach.ExtractVolumeAttachments(page)
				if err != nil {
					return false, fmt.Errorf("Unable to lookup attachment: %s", err)
				}

				attachments = actual
				return true, nil
			})

		if len(attachments) == 1 {
			return nil
		}

		return fmt.Errorf("No attached volume found.")
	}
}

func testAccCheckComputeV2InstanceInstanceIDsDoNotMatch(
	instance1, instance2 *servers.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instance1.ID == instance2.ID {
			return fmt.Errorf("Instance was not recreated.")
		}

		return nil
	}
}

var testAccComputeV2Instance_basic = fmt.Sprintf(`
resource "flexibleengine_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  availability_zone = "%s"
  metadata = {
    foo = "bar"
  }
  network {
    uuid = "%s"
  }
  tags = {
    key1 = "value1"
    key2 = "value2"
  }
}
`, OS_AVAILABILITY_ZONE, OS_NETWORK_ID)

var testAccComputeV2Instance_secgroupMulti = fmt.Sprintf(`
resource "flexibleengine_networking_secgroup_v2" "secgroup_1" {
  name = "flexibleengine_acctest_secgroup_1"
  description = "a security group"
}

resource "flexibleengine_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default", "${flexibleengine_networking_secgroup_v2.secgroup_1.name}"]
  network {
    uuid = "%s"
  }
}
`, OS_NETWORK_ID)

var testAccComputeV2Instance_secgroupMultiUpdate_1 = fmt.Sprintf(`
resource "flexibleengine_networking_secgroup_v2" "secgroup_1" {
  name = "flexibleengine_acctest_secgroup_1"
  description = "a security group"
}

resource "flexibleengine_networking_secgroup_v2" "secgroup_2" {
  name = "flexibleengine_acctest_secgroup_2"
  description = "another security group"
}

resource "flexibleengine_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
}
`, OS_NETWORK_ID)

var testAccComputeV2Instance_secgroupMultiUpdate_2 = fmt.Sprintf(`
resource "flexibleengine_networking_secgroup_v2" "secgroup_1" {
  name = "flexibleengine_acctest_secgroup_1"
  description = "a security group"
}

resource "flexibleengine_networking_secgroup_v2" "secgroup_2" {
  name = "flexibleengine_acctest_secgroup_2"
  description = "another security group"
}

resource "flexibleengine_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default", "${flexibleengine_networking_secgroup_v2.secgroup_1.name}", "${flexibleengine_networking_secgroup_v2.secgroup_2.name}"]
  network {
    uuid = "%s"
  }
}
`, OS_NETWORK_ID)

var testAccComputeV2Instance_bootFromVolumeImage = fmt.Sprintf(`
resource "flexibleengine_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  availability_zone = "%s"
  network {
    uuid = "%s"
  }
  block_device {
    uuid = "%s"
    source_type = "image"
    volume_size = 50
    boot_index = 0
    destination_type = "volume"
    delete_on_termination = true
  }
}
`, OS_AVAILABILITY_ZONE, OS_NETWORK_ID, OS_IMAGE_ID)

var testAccComputeV2Instance_bootFromVolumeVolume = fmt.Sprintf(`
resource "flexibleengine_blockstorage_volume_v2" "vol_1" {
  name = "vol_1"
  size = 50
  image_id = "%s"
}

resource "flexibleengine_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
  block_device {
    uuid = "${flexibleengine_blockstorage_volume_v2.vol_1.id}"
    source_type = "volume"
    boot_index = 0
    destination_type = "volume"
    delete_on_termination = true
  }
}
`, OS_IMAGE_ID, OS_NETWORK_ID)

var testAccComputeV2Instance_bootFromVolumeForceNew_1 = fmt.Sprintf(`
resource "flexibleengine_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
  block_device {
    uuid = "%s"
    source_type = "image"
    volume_size = 50
    boot_index = 0
    destination_type = "volume"
    delete_on_termination = true
  }
}
`, OS_NETWORK_ID, OS_IMAGE_ID)

var testAccComputeV2Instance_bootFromVolumeForceNew_2 = fmt.Sprintf(`
resource "flexibleengine_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
  block_device {
    uuid = "%s"
    source_type = "image"
    volume_size = 51
    boot_index = 0
    destination_type = "volume"
    delete_on_termination = true
  }
}
`, OS_NETWORK_ID, OS_IMAGE_ID)

var testAccComputeV2Instance_blockDeviceNewVolume = fmt.Sprintf(`
resource "flexibleengine_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
  block_device {
    uuid = "%s"
    source_type = "image"
    destination_type = "local"
    boot_index = 0
    delete_on_termination = true
  }
  block_device {
    source_type = "blank"
    destination_type = "volume"
    volume_size = 1
    boot_index = 1
    delete_on_termination = true
  }
}
`, OS_NETWORK_ID, OS_IMAGE_ID)

var testAccComputeV2Instance_blockDeviceExistingVolume = fmt.Sprintf(`
resource "flexibleengine_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  size = 1
}

resource "flexibleengine_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
  block_device {
    uuid = "%s"
    source_type = "image"
    destination_type = "local"
    boot_index = 0
    delete_on_termination = true
  }
  block_device {
    uuid = "${flexibleengine_blockstorage_volume_v2.volume_1.id}"
    source_type = "volume"
    destination_type = "volume"
    boot_index = 1
    delete_on_termination = true
  }
}
`, OS_NETWORK_ID, OS_IMAGE_ID)

var testAccComputeV2Instance_personality = fmt.Sprintf(`
resource "flexibleengine_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
  personality {
    file = "/tmp/foobar.txt"
    content = "happy"
  }
  personality {
    file = "/tmp/barfoo.txt"
    content = "angry"
  }
}
`, OS_NETWORK_ID)

var testAccComputeV2Instance_multiEphemeral = fmt.Sprintf(`
resource "flexibleengine_compute_instance_v2" "instance_1" {
  name = "terraform-test"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
  block_device {
    boot_index = 0
    delete_on_termination = true
    destination_type = "local"
    source_type = "image"
    uuid = "%s"
  }
  block_device {
    boot_index = -1
    delete_on_termination = true
    destination_type = "local"
    source_type = "blank"
    volume_size = 1
  }
  block_device {
    boot_index = -1
    delete_on_termination = true
    destination_type = "local"
    source_type = "blank"
    volume_size = 1
  }
}
`, OS_NETWORK_ID, OS_IMAGE_ID)

var testAccComputeV2Instance_accessIPv4 = fmt.Sprintf(`
resource "flexibleengine_networking_network_v2" "network_1" {
  name = "network_1"
}

resource "flexibleengine_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  network_id = "${flexibleengine_networking_network_v2.network_1.id}"
  cidr = "192.168.1.0/24"
  ip_version = 4
  enable_dhcp = true
  no_gateway = true
}

resource "flexibleengine_compute_instance_v2" "instance_1" {
  depends_on = ["flexibleengine_networking_subnet_v2.subnet_1"]

  name = "instance_1"
  security_groups = ["default"]

  network {
    uuid = "%s"
  }

  network {
    uuid = "${flexibleengine_networking_network_v2.network_1.id}"
    fixed_ip_v4 = "192.168.1.100"
    access_network = true
  }
}
`, OS_NETWORK_ID)

var testAccComputeV2Instance_changeFixedIP_1 = fmt.Sprintf(`
resource "flexibleengine_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
    fixed_ip_v4 = "172.16.10.24"
  }
}
`, OS_NETWORK_ID)

var testAccComputeV2Instance_changeFixedIP_2 = fmt.Sprintf(`
resource "flexibleengine_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
    fixed_ip_v4 = "172.16.10.25"
  }
}
`, OS_NETWORK_ID)

var testAccComputeV2Instance_stopBeforeDestroy = fmt.Sprintf(`
resource "flexibleengine_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
  stop_before_destroy = true
}
`, OS_NETWORK_ID)

var testAccComputeV2Instance_metadataRemove_1 = fmt.Sprintf(`
resource "flexibleengine_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
  metadata = {
    foo = "bar"
    abc = "def"
  }
}
`, OS_NETWORK_ID)

var testAccComputeV2Instance_metadataRemove_2 = fmt.Sprintf(`
resource "flexibleengine_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
  metadata = {
    foo = "bar"
    ghi = "jkl"
  }
}
`, OS_NETWORK_ID)

/*
var testAccComputeV2Instance_forceDelete = fmt.Sprintf(`
resource "flexibleengine_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
  force_delete = true
}
`, OS_NETWORK_ID)
*/

var testAccComputeV2Instance_timeout = fmt.Sprintf(`
resource "flexibleengine_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }

  timeouts {
    create = "10m"
  }
}
`, OS_NETWORK_ID)

var testAccComputeV2Instance_networkNameToID = fmt.Sprintf(`
resource "flexibleengine_networking_network_v2" "network_1" {
  name = "network_1"
}

resource "flexibleengine_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  network_id = "${flexibleengine_networking_network_v2.network_1.id}"
  cidr = "192.168.1.0/24"
  ip_version = 4
  enable_dhcp = true
  no_gateway = true
}

resource "flexibleengine_compute_instance_v2" "instance_1" {
  depends_on = ["flexibleengine_networking_subnet_v2.subnet_1"]

  name = "instance_1"
  security_groups = ["default"]

  network {
    uuid = "%s"
  }

  network {
    name = "${flexibleengine_networking_network_v2.network_1.name}"
  }

}
`, OS_NETWORK_ID)

var testAccComputeV2Instance_crazyNICs = fmt.Sprintf(`
resource "flexibleengine_networking_network_v2" "network_1" {
  name = "network_1"
}

resource "flexibleengine_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  network_id = "${flexibleengine_networking_network_v2.network_1.id}"
  cidr = "192.168.1.0/24"
  ip_version = 4
  enable_dhcp = true
  no_gateway = true
}

resource "flexibleengine_networking_network_v2" "network_2" {
  name = "network_2"
}

resource "flexibleengine_networking_subnet_v2" "subnet_2" {
  name = "subnet_2"
  network_id = "${flexibleengine_networking_network_v2.network_2.id}"
  cidr = "192.168.2.0/24"
  ip_version = 4
  enable_dhcp = true
  no_gateway = true
}

resource "flexibleengine_networking_port_v2" "port_1" {
  name = "port_1"
  network_id = "${flexibleengine_networking_network_v2.network_1.id}"
  admin_state_up = "true"

  fixed_ip {
    subnet_id = "${flexibleengine_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.1.103"
  }
}

resource "flexibleengine_networking_port_v2" "port_2" {
  name = "port_2"
  network_id = "${flexibleengine_networking_network_v2.network_2.id}"
  admin_state_up = "true"

  fixed_ip {
    subnet_id = "${flexibleengine_networking_subnet_v2.subnet_2.id}"
    ip_address = "192.168.2.103"
  }
}

resource "flexibleengine_networking_port_v2" "port_3" {
  name = "port_3"
  network_id = "${flexibleengine_networking_network_v2.network_1.id}"
  admin_state_up = "true"

  fixed_ip {
    subnet_id = "${flexibleengine_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.1.104"
  }
}

resource "flexibleengine_networking_port_v2" "port_4" {
  name = "port_4"
  network_id = "${flexibleengine_networking_network_v2.network_2.id}"
  admin_state_up = "true"

  fixed_ip {
    subnet_id = "${flexibleengine_networking_subnet_v2.subnet_2.id}"
    ip_address = "192.168.2.104"
  }
}

resource "flexibleengine_compute_instance_v2" "instance_1" {
  depends_on = [
    "flexibleengine_networking_subnet_v2.subnet_1",
    "flexibleengine_networking_subnet_v2.subnet_2",
    "flexibleengine_networking_port_v2.port_1",
    "flexibleengine_networking_port_v2.port_2",
  ]

  name = "instance_1"
  security_groups = ["default"]

  network {
    uuid = "%s"
  }

  network {
    uuid = "${flexibleengine_networking_network_v2.network_1.id}"
    fixed_ip_v4 = "192.168.1.100"
  }

  network {
    uuid = "${flexibleengine_networking_network_v2.network_2.id}"
    fixed_ip_v4 = "192.168.2.100"
  }

  network {
    uuid = "${flexibleengine_networking_network_v2.network_1.id}"
    fixed_ip_v4 = "192.168.1.101"
  }

  network {
    uuid = "${flexibleengine_networking_network_v2.network_2.id}"
    fixed_ip_v4 = "192.168.2.101"
  }

  network {
    port = "${flexibleengine_networking_port_v2.port_1.id}"
  }

  network {
    port = "${flexibleengine_networking_port_v2.port_2.id}"
  }

  network {
    port = "${flexibleengine_networking_port_v2.port_3.id}"
  }

  network {
    port = "${flexibleengine_networking_port_v2.port_4.id}"
  }
}
`, OS_NETWORK_ID)

var testAccComputeV2Instance_auto_recovery = fmt.Sprintf(`
resource "flexibleengine_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  availability_zone = "%s"
  metadata = {
    foo = "bar"
  }
  network {
    uuid = "%s"
  }
  auto_recovery = true
}
`, OS_AVAILABILITY_ZONE, OS_NETWORK_ID)
