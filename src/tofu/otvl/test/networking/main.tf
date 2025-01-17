provider "openstack" {
}

terraform {
  required_version = ">= 1.9.0, < 2.0.0"
  required_providers {
    openstack = {
      source  = "terraform-provider-openstack/openstack"
      version = "~> 3.0.0"
    }
  }
}

module "networking" {
  source = "../../modules/networking"
  ext_net_name = var.ext_net_name
  loc_net_name = var.loc_net_name
  loc_net_cidr = var.loc_net_cidr
}
