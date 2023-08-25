package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccImsImageDataSource_basic(t *testing.T) {
	imageName := "OBS CentOS 7.4"
	osVersion := "CentOS 7.4 64bit"
	dataSourceName := "data.flexibleengine_images_image.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccImsImageDataSource_publicName(imageName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", imageName),
					resource.TestCheckResourceAttr(dataSourceName, "protected", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "visibility", "public"),
					resource.TestCheckResourceAttr(dataSourceName, "status", "active"),
				),
			},
			{
				Config: testAccImsImageDataSource_osVersion(osVersion),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "protected", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "visibility", "public"),
					resource.TestCheckResourceAttr(dataSourceName, "status", "active"),
				),
			},
			{
				Config: testAccImsImageDataSource_nameRegex("^OBS CentOS 7.4"),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "protected", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "visibility", "public"),
					resource.TestCheckResourceAttr(dataSourceName, "status", "active"),
				),
			},
		},
	})
}

func TestAccImsImageDataSource_testQueries(t *testing.T) {
	var rName = fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	dataSourceName := "data.flexibleengine_images_image.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccImsImageDataSource_base(rName),
			},
			{
				Config: testAccImsImageDataSource_queryName(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", rName),
					resource.TestCheckResourceAttr(dataSourceName, "protected", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "visibility", "private"),
					resource.TestCheckResourceAttr(dataSourceName, "status", "active"),
				),
			},
			{
				Config: testAccImsImageDataSource_queryTag(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
				),
			},
		},
	})
}

func testAccImsImageDataSource_publicName(imageName string) string {
	return fmt.Sprintf(`
data "flexibleengine_images_image" "test" {
  name        = "%s"
  visibility  = "public"
}
`, imageName)
}

func testAccImsImageDataSource_nameRegex(regexp string) string {
	return fmt.Sprintf(`
data "flexibleengine_images_image" "test" {
  architecture = "x86"
  name_regex   = "%s"
  visibility   = "public"
  most_recent  = true
}
`, regexp)
}

func testAccImsImageDataSource_osVersion(osVersion string) string {
	return fmt.Sprintf(`
data "flexibleengine_images_image" "test" {
  architecture = "x86"
  os_version   = "%s"
  visibility   = "public"
  most_recent  = "true"
}
`, osVersion)
}

func testAccImsImageDataSource_base(rName string) string {
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

func testAccImsImageDataSource_queryName(rName string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_images_image" "test" {
  most_recent = true
  name        = flexibleengine_images_image.test.name
}
`, testAccImsImageDataSource_base(rName))
}

func testAccImsImageDataSource_queryTag(rName string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_images_image" "test" {
  most_recent = true
  visibility  = "private"
  tag         = "foo=bar"
}
`, testAccImsImageDataSource_base(rName))
}
