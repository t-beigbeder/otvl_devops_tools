terraform {
  required_providers {
    openstack = {
      source = "terraform-provider-openstack/openstack"
    }
  }
}

module "instances" {
  source             = "../../../../modules/instances"
  ext_net_id         = var.ext_net_id
  loc_net_id         = var.loc_net_id
  loc_subnet_id      = var.loc_subnet_id
  external_sg_id     = var.sproxy_sg_id
  ssh_key_name       = var.ssh_key_name
  ssh_pub            = var.ssh_pub
  dot_repo           = var.dot_repo
  dot_branch         = var.dot_branch
  rops_repo          = ""
  install_env        = ""
  sproxy_int_address = ""
  bssms_proxy_port   = var.bssms_proxy_port
  instances_attrs = [
    merge(
      var.instance_attr,
      {
        secrets_pri_key = ""
      }
    )
  ]
  go_version         = var.go_version
  user_data_template = "${path.module}/cloud-config.yaml"
}
