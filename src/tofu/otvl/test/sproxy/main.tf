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
  source             = "../../modules/sproxy/network"
  ext_net_name       = var.ext_net_name
  loc_net_name       = var.loc_net_name
  sproxy_sg_name     = var.sproxy_sg_name
  sproxy_port        = "9443"
  sproxy_ssh_exposed = var.sproxy_ssh_exposed
}

module "compute" {
  source           = "../../modules/sproxy/compute"
  ext_net_id       = module.network.ext_net_id
  loc_net_id       = module.network.loc_net_id
  loc_subnet_id    = module.network.loc_subnet_id
  sproxy_sg_id     = module.network.sproxy_sg_id
  ssh_key_name     = var.ssh_key_name
  ssh_pub          = var.ssh_pub
  dot_repo         = var.dot_repo
  dot_branch       = var.dot_branch
  bssms_proxy_port = var.bssms_proxy_port
  go_version       = var.go_version
  instance_attr    = var.instance_attr
}
