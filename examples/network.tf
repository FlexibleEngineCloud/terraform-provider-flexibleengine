resource "flexibleengine_vpc_v1" "vpc" {
  name = "${var.project}-vpc"
  cidr = var.vpc_cidr
}

resource "flexibleengine_vpc_subnet_v1" "subnet" {
  name       = "${var.project}-subnet"
  cidr       = var.subnet_cidr
  gateway_ip = var.gateway_ip
  vpc_id     = flexibleengine_vpc_v1.vpc.id
  dns_list   = ["1.1.1.1", "9.9.9.9"]

}
