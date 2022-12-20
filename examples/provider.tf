terraform {
  required_providers {
    flexibleengine = {
      source  = "FlexibleEngineCloud/flexibleengine"
      version = "1.35.1"
    }
    random = {
      source = "hashicorp/random"
    }
  }
}


provider "random" {}

provider "flexibleengine" {
  access_key  = var.access_key
  secret_key  = var.secret_key
  tenant_name = var.tenant_name
  domain_name = var.domain_name
  region      = var.region
}
