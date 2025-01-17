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

data "openstack_networking_network_v2" "loc_net" {
  name = var.loc_net_name
}

data "openstack_networking_subnet_v2" "loc_net_sn" {
  network_id = data.openstack_networking_network_v2.loc_net.id
  name       = var.loc_net_name
}

resource "openstack_networking_secgroup_v2" "ext" {
  name = var.ext_sg_name
}

resource "openstack_networking_secgroup_rule_v2" "ext_http" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 80
  port_range_max    = 80
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.ext.id
}

resource "openstack_networking_secgroup_rule_v2" "ext_https" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 443
  port_range_max    = 443
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.ext.id
}

resource "openstack_networking_secgroup_rule_v2" "ext_http3" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 9443
  port_range_max    = 9443
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.ext.id
}

resource "openstack_networking_secgroup_rule_v2" "ext_icmp" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "icmp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.ext.id
}
