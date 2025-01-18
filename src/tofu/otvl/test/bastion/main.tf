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

module "network" {
  source           = "../../modules/bastion/network"
  ext_net_name     = var.ext_net_name
  loc_net_name     = var.loc_net_name
  bastion_sg_name  = var.bastion_sg_name
  bssms_proxy_port = var.bssms_proxy_port
}

module "compute" {
  source        = "../../modules/bastion/compute"
  ext_net_id    = module.network.ext_net_id
  loc_net_id    = module.network.loc_net_id
  loc_subnet_id = module.network.loc_subnet_id
  bastion_sg_id = module.network.bastion_sg_id
  ssh_key_name  = var.ssh_key_name
  ssh_pub       = var.ssh_pub
  dot_branch    = var.dot_branch
  dot_repo      = var.dot_repo
  instance_attr = merge(
    var.instance_attr,
    {
      ip_v4 = var.bastion_loc_ip_v4
    }
  )

}
