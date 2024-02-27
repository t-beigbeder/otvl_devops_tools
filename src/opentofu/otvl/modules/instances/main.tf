terraform {
  required_version = ">= 1.0.0, < 2.0.0"
  required_providers {
    openstack = {
      source  = "terraform-provider-openstack/openstack"
      version = "~> 1.42.0"
    }
  }
}

resource "openstack_compute_keypair_v2" "this" {
  name       = var.ssh_key_name
  public_key = var.ssh_pub
}

resource "openstack_networking_port_v2" "ext" {
  count              = length(var.instances_attrs)
  network_id         = var.ext_net_id
  security_group_ids = [count.index == 0 ? var.bastion_sg_id : var.ext_sg_id]
}

resource "openstack_networking_port_v2" "loc" {
  count      = length(var.instances_attrs)
  network_id = var.loc_net_id
  fixed_ip {
    subnet_id  = var.loc_subnet_id
    ip_address = var.instances_attrs[count.index].ip_v4
  }
  admin_state_up = "true"
}

resource "openstack_compute_instance_v2" "this" {
  count           = length(var.instances_attrs)
  name            = var.instances_attrs[count.index].name
  image_name      = var.instances_attrs[count.index].image_name
  flavor_name     = var.instances_attrs[count.index].flavor_name
  key_pair        = openstack_compute_keypair_v2.this.name
  user_data       = var.instance_user_data
  security_groups = []
  network {
    port = openstack_networking_port_v2.ext[count.index].id
  }
  network {
    port = openstack_networking_port_v2.loc[count.index].id
  }
  metadata = {
    "groups"    = var.instances_attrs[count.index].groups
    "otvl_meta" = var.instances_attrs[count.index].otvl_meta
  }
}

data "openstack_networking_port_v2" "ext" {
  count      = length(var.instances_attrs)
  network_id = var.ext_net_id
  device_id  = openstack_compute_instance_v2.this[count.index].id
}
