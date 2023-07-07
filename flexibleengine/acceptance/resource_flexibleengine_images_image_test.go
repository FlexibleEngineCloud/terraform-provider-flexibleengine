package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/ims/v2/cloudimages"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/ims"
)

func TestAccImsImage_basic(t *testing.T) {
	var image cloudimages.Image

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	rNameUpdate := rName + "-update"
	resourceName := "flexibleengine_images_image.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckImsImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImsImage_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImsImageExists(resourceName, &image),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttr(resourceName, "min_ram", "0"),
					resource.TestCheckResourceAttr(resourceName, "max_ram", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				Config: testAccImsImage_update(rName, rNameUpdate, 1024, 4096),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImsImageExists(resourceName, &image),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttr(resourceName, "description",
						"created by Terraform AccTest for update"),
					resource.TestCheckResourceAttr(resourceName, "min_ram", "1024"),
					resource.TestCheckResourceAttr(resourceName, "max_ram", "4096"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccImsImage_update(rName, rNameUpdate, 4096, 8192),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImsImageExists(resourceName, &image),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttr(resourceName, "description",
						"created by Terraform AccTest for update"),
					resource.TestCheckResourceAttr(resourceName, "min_ram", "4096"),
					resource.TestCheckResourceAttr(resourceName, "max_ram", "8192"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccImsImage_update(rName, rNameUpdate, 0, 0),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImsImageExists(resourceName, &image),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttr(resourceName, "description",
						"created by Terraform AccTest for update"),
					resource.TestCheckResourceAttr(resourceName, "min_ram", "0"),
					resource.TestCheckResourceAttr(resourceName, "max_ram", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
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

func TestAccImsImage_wholeImage_withServer(t *testing.T) {
	var image cloudimages.Image

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	rNameUpdate := rName + "-update"
	resourceName := "flexibleengine_images_image.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckImsImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImsImage_wholeImage_withServer(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImsImageExists(resourceName, &image),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				Config: testAccImsImage_wholeImage_withServer_update(rName, rNameUpdate, 1024, 4096),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImsImageExists(resourceName, &image),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttr(resourceName, "min_ram", "1024"),
					resource.TestCheckResourceAttr(resourceName, "max_ram", "4096"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccImsImage_wholeImage_withServer_update(rName, rNameUpdate, 4096, 8192),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImsImageExists(resourceName, &image),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttr(resourceName, "min_ram", "4096"),
					resource.TestCheckResourceAttr(resourceName, "max_ram", "8192"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccImsImage_wholeImage_withServer_update(rName, rNameUpdate, 0, 0),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImsImageExists(resourceName, &image),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttr(resourceName, "min_ram", "0"),
					resource.TestCheckResourceAttr(resourceName, "max_ram", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"vault_id"},
			},
		},
	})
}

func TestAccImsImage_wholeImage_withBackup(t *testing.T) {
	var image cloudimages.Image

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	rNameUpdate := rName + "-update"
	resourceName := "flexibleengine_images_image.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckImsBackupId(t)
		},
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckImsImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImsImage_wholeImage_withBackup(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImsImageExists(resourceName, &image),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				Config: testAccImsImage_wholeImage_withBackup_update(rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImsImageExists(resourceName, &image),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

func testAccCheckImsImageDestroy(s *terraform.State) error {
	cfg := testAccProvider.Meta().(*config.Config)
	imageClient, err := cfg.ImageV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating Image: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_images_image" {
			continue
		}

		_, err := ims.GetCloudImage(imageClient, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("image still exists")
		}
	}

	return nil
}

func testAccCheckImsImageExists(n string, image *cloudimages.Image) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("IMS Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		cfg := testAccProvider.Meta().(*config.Config)
		imageClient, err := cfg.ImageV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating Image: %s", err)
		}

		found, err := ims.GetCloudImage(imageClient, rs.Primary.ID)
		if err != nil {
			return err
		}

		*image = *found
		return nil
	}
}

func testAccImsImage_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "flexibleengine_availability_zones" "test" {}

data "flexibleengine_compute_flavors_v2" "test" {
  availability_zone = data.flexibleengine_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core          = 2
  memory_size       = 4
}

resource "flexibleengine_compute_instance_v2" "test" {
  name               = "%[2]s"
  image_name         = "OBS Ubuntu 18.04"
  flavor_id          = data.flexibleengine_compute_flavors_v2.test.flavors[0]
  security_groups    = [flexibleengine_networking_secgroup_v2.test.name]
  availability_zone  = data.flexibleengine_availability_zones.test.names[0]

  network {
    uuid = flexibleengine_vpc_subnet_v1.test.id
  }
}

resource "flexibleengine_images_image" "test" {
  name        = "%[2]s"
  instance_id = flexibleengine_compute_instance_v2.test.id

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testBaseNetwork(rName), rName)
}

func testAccImsImage_update(rName, rNameUpdate string, minRAM, maxRAM int) string {
	return fmt.Sprintf(`
%[1]s

data "flexibleengine_availability_zones" "test" {}

data "flexibleengine_compute_flavors_v2" "test" {
  availability_zone = data.flexibleengine_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core          = 2
  memory_size       = 4
}

resource "flexibleengine_compute_instance_v2" "test" {
  name               = "%[2]s"
  image_name         = "OBS Ubuntu 18.04"
  flavor_id          = data.flexibleengine_compute_flavors_v2.test.flavors[0]
  security_groups    = [flexibleengine_networking_secgroup_v2.test.name]
  availability_zone  = data.flexibleengine_availability_zones.test.names[0]

  network {
    uuid = flexibleengine_vpc_subnet_v1.test.id
  }
}

resource "flexibleengine_images_image" "test" {
  name        = "%[3]s"
  instance_id = flexibleengine_compute_instance_v2.test.id
  description = "created by Terraform AccTest for update"
  min_ram     = %[4]d
  max_ram     = %[5]d

  tags = {
    foo  = "bar"
    key  = "value1"
    key2 = "value2"
  }
}
`, testBaseNetwork(rName), rName, rNameUpdate, minRAM, maxRAM)
}

func testAccImsImage_wholeImage_base(rName string) string {
	return fmt.Sprintf(`
data "flexibleengine_availability_zones" "test" {}

data "flexibleengine_compute_flavors_v2" "test" {
  availability_zone = data.flexibleengine_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core          = 2
  memory_size       = 4
}

resource "flexibleengine_compute_instance_v2" "test" {
  name               = "%[1]s"
  image_name         = "OBS Ubuntu 18.04"
  flavor_id          = data.flexibleengine_compute_flavors_v2.test.flavors[0]
  security_groups    = [flexibleengine_networking_secgroup_v2.test.name]
  availability_zone  = data.flexibleengine_availability_zones.test.names[0]

  network {
    uuid = flexibleengine_vpc_subnet_v1.test.id
  }
}

resource "flexibleengine_cbr_vault" "test" {
  name             = "%[1]s"
  type             = "server"
  consistent_level = "app_consistent"
  protection_type  = "backup"
  size             = 200

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, rName)
}

func testAccImsImage_wholeImage_withServer(rName string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "flexibleengine_images_image" "test" {
  name        = "%[3]s"
  instance_id = flexibleengine_compute_instance_v2.test.id
  description = "created by Terraform AccTest"
  vault_id    = flexibleengine_cbr_vault.test.id

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testBaseNetwork(rName), testAccImsImage_wholeImage_base(rName), rName)
}

func testAccImsImage_wholeImage_withServer_update(rName, updateName string, minRAM, maxRAM int) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "flexibleengine_images_image" "test" {
  name        = "%[3]s"
  instance_id = flexibleengine_compute_instance_v2.test.id
  description = "created by Terraform AccTest"
  vault_id    = flexibleengine_cbr_vault.test.id
  min_ram     = %[4]d
  max_ram     = %[5]d

  tags = {
    foo  = "bar"
    key  = "value1"
    key2 = "value2"
  }
}
`, testBaseNetwork(rName), testAccImsImage_wholeImage_base(rName), updateName, minRAM, maxRAM)
}

func testAccImsImage_wholeImage_withBackup(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_images_image" "test" {
  name        = "%[1]s"
  backup_id   = "%[2]s"
  description = "created by Terraform AccTest"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, rName, OS_IMS_BACKUP_ID)
}

func testAccImsImage_wholeImage_withBackup_update(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_images_image" "test" {
  name        = "%[1]s"
  backup_id   = "%[2]s"
  description = "created by Terraform AccTest"

  tags = {
    foo  = "bar"
    key  = "value1"
    key2 = "value2"
  }
}
`, rName, OS_IMS_BACKUP_ID)
}
