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

func testAccDcsV1Instance_network(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_vpc_v1" "vpc_1" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "subnet_1" {
  name       = "%[1]s"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = flexibleengine_vpc_v1.vpc_1.id
}

resource "flexibleengine_networking_secgroup_v2" "secgroup_1" {
  name        = "%[1]s"
  description = "secgroup_1"
}
`, rName)
}

func testAccDcsV1Instance_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_dcs_flavors" "test" {
  cache_mode     = "ha"
  capacity       = 0.125
  engine_version = "5.0"
}

data "flexibleengine_dcs_product_v1" "product1" {
  engine         = "Redis"
  engine_version = "4.0;5.0"
  cache_mode     = "cluster"
  capacity       = 8
  replica_count  = 2
}

resource "flexibleengine_dcs_instance_v1" "instance_1" {
  name              = "%s"
  engine            = "Redis"
  engine_version    = "5.0"
  password          = "FlexibleEngine_test"
  product_id        = data.flexibleengine_dcs_product_v1.product1.id
  capacity          = 8
  vpc_id            = flexibleengine_vpc_v1.vpc_1.id
  network_id        = flexibleengine_vpc_subnet_v1.subnet_1.id
  available_zones   = ["eu-west-0a", "eu-west-0b"]

  save_days       = 1
  backup_type     = "manual"
  begin_at        = "00:00-01:00"
  period_type     = "weekly"
  backup_at       = [1]
}
`, testAccDcsV1Instance_network(rName), rName)
}
