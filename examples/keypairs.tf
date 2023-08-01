resource "flexibleengine_compute_keypair_v2" "keypair" {
  name       = "${var.project}-terraform_key"
  public_key = file(var.ssh_pub_key)
}
