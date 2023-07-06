package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccImsImagesDataSource_basic(t *testing.T) {
	imageName := "OBS CentOS 7.4"
	osVersion := "CentOS 7.4 64bit"
	dataSourceName := "data.flexibleengine_images_images.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccImsImagesDataSource_publicName(imageName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.name", imageName),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.protected", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.visibility", "public"),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.status", "active"),
				),
			},
			{
				Config: testAccImsImagesDataSource_osVersion(osVersion),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.protected", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.visibility", "public"),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.status", "active"),
				),
			},
			{
				Config: testAccImsImagesDataSource_nameRegex("^OBS CentOS 7.4"),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.protected", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.visibility", "public"),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.status", "active"),
				),
			},
		},
	})
}

func TestAccImsImagesDataSource_testQueries(t *testing.T) {
	var rName = fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	dataSourceName := "data.flexibleengine_images_images.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccImsImagesDataSource_base(rName),
			},
			{
				Config: testAccImsImagesDataSource_queryName(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.name", rName),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.protected", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.visibility", "private"),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.status", "active"),
				),
			},
			{
				Config: testAccImsImagesDataSource_queryTag(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
				),
			},
		},
	})
}

func testAccImsImagesDataSource_publicName(imageName string) string {
	return fmt.Sprintf(`
data "flexibleengine_images_images" "test" {
  name       = "%s"
  visibility = "public"
}
`, imageName)
}

func testAccImsImagesDataSource_nameRegex(regexp string) string {
	return fmt.Sprintf(`
data "flexibleengine_images_images" "test" {
  architecture = "x86"
  name_regex   = "%s"
  visibility   = "public"
}
`, regexp)
}

func testAccImsImagesDataSource_osVersion(osVersion string) string {
	return fmt.Sprintf(`
data "flexibleengine_images_images" "test" {
  architecture = "x86"
  os_version   = "%s"
  visibility   = "public"
}
`, osVersion)
}

func testAccImsImagesDataSource_base(rName string) string {
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
  description = "created by Terraform AccTest"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testBaseNetwork(rName), rName)
}

func testAccImsImagesDataSource_queryName(rName string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_images_images" "test" {
  name = flexibleengine_images_image.test.name
}
`, testAccImsImagesDataSource_base(rName))
}

func testAccImsImagesDataSource_queryTag(rName string) string {
	return fmt.Sprintf(`
%s
data "flexibleengine_images_images" "test" {
  visibility = "private"
  tag        = "foo=bar"
}
`, testAccImsImagesDataSource_base(rName))
}
