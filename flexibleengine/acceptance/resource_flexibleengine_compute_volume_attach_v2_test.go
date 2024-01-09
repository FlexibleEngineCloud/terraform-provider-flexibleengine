package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/ecs/v1/block_devices"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getVolumeAttachResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.ComputeV1Client(OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating compute v1 client: %s", err)
	}

	instanceId := state.Primary.Attributes["instance_id"]
	volumeId := state.Primary.Attributes["volume_id"]
	found, err := block_devices.Get(c, instanceId, volumeId).Extract()
	if err != nil {
		return nil, err
	}

	if found.ServerId != instanceId || found.VolumeId != volumeId {
		return nil, fmt.Errorf("volume attach not found %s", state.Primary.ID)
	}

	return found, nil
}

func TestAccComputeVolumeAttach_basic(t *testing.T) {
	var va block_devices.VolumeAttachment
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "flexibleengine_compute_volume_attach_v2.va_1"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&va,
		getVolumeAttachResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeVolumeAttach_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id",
						"flexibleengine_compute_instance_v2.instance_1", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "volume_id", "flexibleengine_blockstorage_volume_v2.test", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "pci_address"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeVolumeAttach_device(t *testing.T) {
	var va block_devices.VolumeAttachment
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "flexibleengine_compute_volume_attach_v2.va_1"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&va,
		getVolumeAttachResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeVolumeAttach_device(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id",
						"flexibleengine_compute_instance_v2.instance_1", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "volume_id", "flexibleengine_blockstorage_volume_v2.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "device", "/dev/vdb"),
					resource.TestCheckResourceAttrSet(resourceName, "pci_address"),
				),
			},
		},
	})
}

func TestAccComputeVolumeAttach_multiple(t *testing.T) {
	var va block_devices.VolumeAttachment
	rName := acceptance.RandomAccResourceNameWithDash()
	rc := acceptance.InitResourceCheck(
		"flexibleengine_compute_volume_attach_v2.test",
		&va,
		getVolumeAttachResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeVolumeAttach_multiple(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckMultiResourcesExists(2),
					resource.TestCheckResourceAttrPair("flexibleengine_compute_volume_attach_v2.test.0", "instance_id",
						"flexibleengine_compute_instance_v2.test.0", "id"),
					resource.TestCheckResourceAttrPair("flexibleengine_compute_volume_attach_v2.test.0", "volume_id",
						"flexibleengine_blockstorage_volume_v2.test", "id"),
					resource.TestCheckResourceAttrPair("flexibleengine_compute_volume_attach_v2.test.1", "instance_id",
						"flexibleengine_compute_instance_v2.test.1", "id"),
					resource.TestCheckResourceAttrPair("flexibleengine_compute_volume_attach_v2.test.1", "volume_id",
						"flexibleengine_blockstorage_volume_v2.test", "id"),
				),
			},
		},
	})
}

const testAccCompute_data = `
data "flexibleengine_availability_zones" "test" {}

data "flexibleengine_compute_flavors_v2" "test" {
  availability_zone = data.flexibleengine_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core          = 2
  memory_size       = 4
}

data "flexibleengine_vpc_subnet_v1" "test" {
  name = "subnet-default"
}

data "flexibleengine_images_image" "test" {
  name        = "OBS Ubuntu 18.04"
  most_recent = true
}

data "flexibleengine_networking_secgroup_v2" "test" {
  name = "default"
}
`

func testAccComputeVolumeAttach_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_blockstorage_volume_v2" "test" {
  name              = "%s"
  availability_zone = data.flexibleengine_availability_zones.test.names[0]
  volume_type       = "SAS"
  size              = 10
}

resource "flexibleengine_compute_instance_v2" "instance_1" {
  name               = "%s"
  image_id           = data.flexibleengine_images_image.test.id
  flavor_id          = data.flexibleengine_compute_flavors_v2.test.flavors[0]
  security_groups    = [data.flexibleengine_networking_secgroup_v2.test.name]
  availability_zone  = data.flexibleengine_availability_zones.test.names[0]

  network {
    uuid = data.flexibleengine_vpc_subnet_v1.test.id
  }
}

resource "flexibleengine_compute_volume_attach_v2" "va_1" {
  instance_id = flexibleengine_compute_instance_v2.instance_1.id
  volume_id   = flexibleengine_blockstorage_volume_v2.test.id
}
`, testAccCompute_data, rName, rName)
}

func testAccComputeVolumeAttach_device(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_blockstorage_volume_v2" "test" {
  name              = "%s"
  availability_zone = data.flexibleengine_availability_zones.test.names[0]
  volume_type       = "SAS"
  size              = 10
}

resource "flexibleengine_compute_instance_v2" "instance_1" {
  name               = "%s"
  image_id           = data.flexibleengine_images_image.test.id
  flavor_id          = data.flexibleengine_compute_flavors_v2.test.flavors[0]
  security_groups    = [data.flexibleengine_networking_secgroup_v2.test.name]
  availability_zone  = data.flexibleengine_availability_zones.test.names[0]
  network {
    uuid = data.flexibleengine_vpc_subnet_v1.test.id
  }
}

resource "flexibleengine_compute_volume_attach_v2" "va_1" {
  instance_id = flexibleengine_compute_instance_v2.instance_1.id
  volume_id   = flexibleengine_blockstorage_volume_v2.test.id
  device      = "/dev/vdb"
}
`, testAccCompute_data, rName, rName)
}

func testAccComputeVolumeAttach_multiple(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_blockstorage_volume_v2" "test" {
  name              = "%[2]s"
  availability_zone = data.flexibleengine_availability_zones.test.names[0]
  volume_type       = "SAS"
  size              = 10
  multiattach       = true
}

resource "flexibleengine_compute_instance_v2" "test" {
  count              = 2
  name               = "%[2]s-${count.index}"
  image_id           = data.flexibleengine_images_image.test.id
  flavor_id          = data.flexibleengine_compute_flavors_v2.test.flavors[0]
  security_groups    = [data.flexibleengine_networking_secgroup_v2.test.name]
  availability_zone  = data.flexibleengine_availability_zones.test.names[0]

  network {
    uuid = data.flexibleengine_vpc_subnet_v1.test.id
  }
}

resource "flexibleengine_compute_volume_attach_v2" "test" {
  count       = 2
  instance_id = flexibleengine_compute_instance_v2.test[count.index].id
  volume_id   = flexibleengine_blockstorage_volume_v2.test.id
}
`, testAccCompute_data, rName, rName)
}
