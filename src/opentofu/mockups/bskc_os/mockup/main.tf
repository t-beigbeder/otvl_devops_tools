provider "openstack" {
  alias = "osdp"
}

terraform {
  required_version = ">= 1.6.2, < 2.0.0"

  required_providers {
    openstack = {
      source  = "terraform-provider-openstack/openstack"
      version = "~> 1.42.0"
    }
    bssms = {
      source = "tofu.otvl.org/otvl/bssms"
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

provider "bssms" {
  proxy_address = "localhost:9443"
}

module "networking" {
  source          = "../modules/networking"
  ext_net_name    = var.ext_net_name
  loc_net_name    = var.loc_net_name
  loc_net_cidr    = var.loc_net_cidr
  bastion_sg_name = var.bastion_sg_name
  ext_sg_name     = var.ext_sg_name
}

module "instances" {
  source               = "../modules/instances"
  ext_net_id           = module.networking.ext_net_id
  loc_net_id           = module.networking.loc_net_id
  loc_subnet_id        = module.networking.loc_subnet_id
  ssh_key_name         = var.ssh_key_name
  ssh_pub              = var.ssh_pub
  instances_attrs      = var.instances_attrs
  bastion_sg_id        = module.networking.bastion_sg_id
  ext_sg_id            = module.networking.ext_sg_id
  dot_repo             = var.dot_repo
  dot_branch           = var.dot_branch
  go_version           = var.go_version
  bssms_proxy_hostname = var.bssms_proxy_hostname
  bssms_proxy_port     = var.bssms_proxy_port
}
