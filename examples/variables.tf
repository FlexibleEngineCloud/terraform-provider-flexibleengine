### FlexibleEngine Credentials
variable "access_key" {
  # If you don't fill this in, you will be prompted for it
  #default = "your_access_key"
}

variable "secret_key" {
  # If you don't fill this in, you will be prompted for it
  #default = "your_secret_key'
}

variable "domain_name" {
  # If you don't fill this in, you will be prompted for it
  #default = "your_domainname"
}

variable "tenant_name" {
  default = "eu-west-0"
}

variable "region" {
  default = "eu-west-0"
}

### Project Settings
variable "project" {
  default = "terraform"
}

variable "vpc_cidr" {
  default = "192.168.10.0/24"
}

variable "subnet_cidr" {
  default = "192.168.10.0/24"
}

variable "gateway_ip" {
  default = "192.168.10.1"
}

variable "ssh_pub_key" {
  default = "~/.ssh/id_rsa.pub"
}

### VM (Instance) Settings
variable "instance_count" {
  default = "1"
}

variable "disk_size_gb" {
  default = "0"
}

variable "flavor_name" {
  default = "s6.medium.2"
}

variable "image_name" {
  default = "OBS Ubuntu 22.04"
}
