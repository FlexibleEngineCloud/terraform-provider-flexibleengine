package acceptance

import "fmt"

// testSecGroup can be referred as `flexibleengine_networking_secgroup_v2.test`
func testSecGroup(name string) string {
	return fmt.Sprintf(`
resource "flexibleengine_networking_secgroup_v2" "test" {
  name                 = "%s"
  delete_default_rules = true
}
`, name)
}

// testVpc can be referred as `flexibleengine_vpc_v1.test` and `flexibleengine_vpc_subnet_v1.test`
func testVpc(name string) string {
	return fmt.Sprintf(`
resource "flexibleengine_vpc_v1" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "test" {
  name       = "%[1]s"
  vpc_id     = flexibleengine_vpc_v1.test.id
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
}
`, name)
}

// testBaseNetwork vpc, subnet, security group
func testBaseNetwork(name string) string {
	return fmt.Sprintf(`
# base security group without default rules
%s

# base vpc and subnet
%s
`, testSecGroup(name), testVpc(name))
}

// testBaseComputeResources vpc, subnet, security group, availability zone, keypair, image, flavor
func testBaseComputeResources(name string) string {
	return fmt.Sprintf(`
# base test resources
%s

data "flexibleengine_availability_zones" "test" {}

data "flexibleengine_compute_flavors_v2" "test" {
  availability_zone = data.flexibleengine_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core          = 2
  memory_size       = 4
}

data "flexibleengine_images_image" "test" {
  name = "OBS Ubuntu 18.04"
}
`, testBaseNetwork(name))
}
