terraform {
  required_providers {
    openstack = {
      source  = "terraform-provider-openstack/openstack"
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
  network_id = openstack_networking_network_v2.loc_net.id
  name       = var.loc_net_name
  cidr       = var.loc_net_cidr
}
