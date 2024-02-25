terraform {
  required_version = ">= 1.0.0, < 2.0.0"
  required_providers {
    openstack = {
      source  = "terraform-provider-openstack/openstack"
      version = "~> 1.42.0"
    }
  }
}

data "openstack_networking_network_v2" "ext_net" {
  name = var.ext_net_name
}

resource "openstack_networking_network_v2" "loc_net" {
  name           = var.loc_net_name
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "loc_net_sn" {
  network_id  = openstack_networking_network_v2.loc_net.id
  name        = var.loc_net_name
  enable_dhcp = "false"
  cidr        = var.loc_net_cidr
}
