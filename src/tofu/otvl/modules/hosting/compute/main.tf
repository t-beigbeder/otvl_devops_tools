terraform {
  required_providers {
    openstack = {
      source = "terraform-provider-openstack/openstack"
    }
  }
}

module "instances" {
  source               = "../../../../modules/instances"
  ext_net_id           = var.ext_net_id
  loc_net_id           = var.loc_net_id
  loc_subnet_id        = var.loc_subnet_id
  external_sg_id       = var.ext_net_id
  ssh_key_name         = var.ssh_key_name
  ssh_pub              = var.ssh_pub
  dot_repo             = var.dot_repo
  dot_branch           = var.dot_branch
  bssms_proxy_hostname = var.bastion_loc_ip_v4
  bssms_proxy_port     = var.bssms_proxy_port
  instances_attrs = [
    merge(
      var.instances_attrs,
      {
        secrets_pri_key = ""
      }
    )
  ]
  go_version         = ""
  user_data_template = "${path.module}/cloud-config.yaml"
  yaml_secrets       = ""
}