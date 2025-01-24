terraform {
  required_providers {
    openstack = {
      source = "terraform-provider-openstack/openstack"
    }
    bssms = {
      source = "tofu.otvl.org/otvl/bssms"
    }
  }
}

resource "bssms_installable" "this" {
  name = var.instance_attr.name
}

module "instances" {
  source               = "../../../../modules/instances"
  ext_net_id           = var.ext_net_id
  loc_net_id           = var.loc_net_id
  loc_subnet_id        = var.loc_subnet_id
  external_sg_id       = var.bastion_sg_id
  ssh_key_name         = var.ssh_key_name
  ssh_pub              = var.ssh_pub
  dot_repo             = var.dot_repo
  dot_branch           = var.dot_branch
  rops_repo            = var.rops_repo
  install_env          = var.install_env
  b64_id_rsa_rops      = var.b64_id_rsa_rops
  bssms_proxy_hostname = var.bssms_proxy_address
  bssms_proxy_port     = var.bssms_proxy_port
  instances_attrs = [
    merge(
      var.instance_attr,
      {
        secrets_pri_key = ""
      }
    )
  ]
  go_version           = var.go_version
  user_data_template   = "${path.module}/cloud-config.yaml"
}

resource "bssms_secrets" "this" {
  name            = var.instance_attr.name
  pri_key         = resource.bssms_installable.this.pri_key
  pub_key         = resource.bssms_installable.this.pub_key
  server_uuid     = module.instances[0].ids
  ip_v4_addresses = module.instances.ipv4s[0]
  mac_addresses   = module.instances.macs[0]
  yaml_secrets    = var.yaml_secrets
}
