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
  count = length(var.instances_attrs)
  name = var.instances_attrs[count.index].name
}

locals {
  instances_attrs = [
    for index, ia in var.instances_attrs :
    merge(
      ia,
      {
        secrets_pri_key = bssms_installable.this[index].pri_key
      }
    )
  ]
}

module "instances" {
  source             = "../../../../modules/instances"
  ext_net_id         = var.ext_net_id
  loc_net_id         = var.loc_net_id
  loc_subnet_id      = var.loc_subnet_id
  external_sg_id     = var.hosting_sg_id
  ssh_key_name       = var.ssh_key_name
  ssh_pub            = var.ssh_pub
  dot_repo           = var.dot_repo
  dot_branch         = var.dot_branch
  rops_repo          = var.rops_repo
  install_env        = var.install_env
  b64_id_rsa_rops    = var.b64_id_rsa_rops
  sproxy_int_address = var.bastion_loc_ip_v4
  bssms_proxy_port   = var.bssms_proxy_port
  instances_attrs    = local.instances_attrs
  go_version         = var.go_version
  user_data_template = "${path.module}/cloud-config.yaml"
}

resource "bssms_secrets" "this" {
  count = length(var.instances_attrs)
  name            = var.instances_attrs[count.index].name
  pri_key         = resource.bssms_installable.this[count.index].pri_key
  pub_key         = resource.bssms_installable.this[count.index].pub_key
  server_uuid     = module.instances.ids[count.index]
  ip_v4_addresses = module.instances.ipv4s[count.index]
  mac_addresses   = module.instances.macs[count.index]
  yaml_secrets    = var.yaml_secrets
}
