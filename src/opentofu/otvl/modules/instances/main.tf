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

resource "openstack_compute_instance_v2" "this" {
  count       = length(var.instances_attrs)
  name        = var.instances_attrs[count.index].name
  image_name  = var.instances_attrs[count.index].image_name
  flavor_name = var.instances_attrs[count.index].flavor_name
  key_pair    = openstack_compute_keypair_v2.this.name
  user_data   = var.instance_user_data
  security_groups = []
  network {
    uuid = var.ext_net_id
  }
  network {
    uuid        = var.loc_net_id
    fixed_ip_v4 = var.instances_attrs[count.index].ip_v4
  }
  metadata = {
    "hostname"     = var.instances_attrs[count.index].name
    "logical_name" = var.instances_attrs[count.index].logical_name
  }
}

resource "openstack_networking_port_v2" "ext" {
  count       = length(var.instances_attrs)
  network_id  = var.ext_net_id
  security_group_ids = [var.ext_sg_id]
}

data "openstack_networking_port_v2" "ext" {
  count       = length(var.instances_attrs)
  device_id = openstack_compute_instance_v2.this[count.index].id
  network_id = var.ext_net_id
}
