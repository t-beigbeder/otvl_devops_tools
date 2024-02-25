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
  network_id = openstack_networking_network_v2.loc_net.id
  name       = var.loc_net_name
  cidr       = var.loc_net_cidr
}

resource "openstack_networking_secgroup_v2" "bastion" {
  name = var.bastion_sg_name
}

resource "openstack_networking_secgroup_rule_v2" "bastion_http" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 80
  port_range_max    = 80
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.bastion.id
}

resource "openstack_networking_secgroup_rule_v2" "bastion_https" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 443
  port_range_max    = 443
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.bastion.id
}

resource "openstack_networking_secgroup_rule_v2" "bastion_ssh" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 22
  port_range_max    = 22
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.bastion.id
}

resource "openstack_networking_secgroup_rule_v2" "bastion_ssh6" {
  direction         = "ingress"
  ethertype         = "IPv6"
  protocol          = "tcp"
  port_range_min    = 22
  port_range_max    = 22
  remote_ip_prefix  = "::/0"
  security_group_id = openstack_networking_secgroup_v2.bastion.id
}

resource "openstack_networking_secgroup_rule_v2" "bastion_icmp" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "icmp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.bastion.id
}
#
#resource "openstack_networking_secgroup_rule_v2" "bastion_out_ipv4" {
#  direction         = "egress"
#  ethertype         = "IPv4"
#  remote_ip_prefix  = "0.0.0.0/0"
#  security_group_id = openstack_networking_secgroup_v2.bastion.id
#}
#
#resource "openstack_networking_secgroup_rule_v2" "bastion_out_ipv6" {
#  direction         = "egress"
#  ethertype         = "IPv6"
#  remote_ip_prefix  = "::/0"
#  security_group_id = openstack_networking_secgroup_v2.bastion.id
#}

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

resource "openstack_networking_secgroup_rule_v2" "ext_icmp" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "icmp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.ext.id
}
#
#resource "openstack_networking_secgroup_rule_v2" "ext_out_ipv4" {
#  direction         = "egress"
#  ethertype         = "IPv4"
#  remote_ip_prefix  = "0.0.0.0/0"
#  security_group_id = openstack_networking_secgroup_v2.ext.id
#}
#
#resource "openstack_networking_secgroup_rule_v2" "ext_out_ipv6" {
#  direction         = "egress"
#  ethertype         = "IPv6"
#  remote_ip_prefix  = "::/0"
#  security_group_id = openstack_networking_secgroup_v2.ext.id
#}
