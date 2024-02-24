provider "openstack" {
  alias = "osdp"
}

terraform {
  required_version = ">= 1.0.0, < 2.0.0"

  required_providers {
    openstack = {
      source  = "terraform-provider-openstack/openstack"
      version = "~> 1.42.0"
    }
  }

  /*
    backend "s3" {
      bucket         = "default-tf-bucket"
      key            = "otvl/prod/terraform.tfstate"
      region         = "here"
    }
  */

}

module "networking" {
  source       = "../modules/networking"
  ext_net_name = var.ext_net_name
  loc_net_name = var.loc_net_name
  loc_net_cidr = var.loc_net_cidr
}

module "instances" {
  source             = "../modules/instances"
  ext_net_id         = module.networking.ext_net_id
  loc_net_id         = module.networking.loc_net_id
  ssh_key_name       = var.ssh_key_name
  ssh_pub            = var.ssh_pub
  instances_attrs    = var.instances_attrs
  instance_user_data = var.instance_user_data
}
