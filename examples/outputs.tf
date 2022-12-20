output "webserver_url" {
  value = "http://${flexibleengine_vpc_eip.eip_1.address}/"
}
