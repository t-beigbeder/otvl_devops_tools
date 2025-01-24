provider "openstack" {
}

terraform {
  required_version = ">= 1.9.0, < 2.0.0"
  required_providers {
    openstack = {
      source  = "terraform-provider-openstack/openstack"
      version = "~> 3.0.0"
    }
    sops = {
      source  = "carlpett/sops"
      version = "~> 1.1.1"
    }
    bssms = {
      source = "tofu.otvl.org/otvl/bssms"
    }
  }
}

provider "bssms" {
  proxy_address = "${var.sproxy_ext_address}:${var.bssms_proxy_port}"
}

module "network" {
  source          = "../../modules/bastion/network"
  ext_net_name    = var.ext_net_name
  loc_net_name    = var.loc_net_name
  bastion_sg_name = var.bastion_sg_name
}

data "sops_file" "hosting_secret" {
  source_file = var.hosting_secrets_sops
}

module "compute" {
  source             = "../../modules/bastion/compute"
  ext_net_id         = module.network.ext_net_id
  loc_net_id         = module.network.loc_net_id
  loc_subnet_id      = module.network.loc_subnet_id
  bastion_sg_id      = module.network.bastion_sg_id
  ssh_key_name       = var.ssh_key_name
  ssh_pub            = var.ssh_pub
  b64_id_rsa_rops    = var.b64_id_rsa_rops
  dot_repo           = var.dot_repo
  dot_branch         = var.dot_branch
  rops_repo          = var.rops_repo
  install_env        = var.install_env
  sproxy_int_address = var.sproxy_int_address
  bssms_proxy_port   = var.bssms_proxy_port
  go_version         = var.go_version
  yaml_secrets       = data.sops_file.hosting_secret.raw
  instance_attr      = var.instance_attr
}
