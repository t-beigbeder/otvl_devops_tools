output "ext_net_id" {
  value       = data.openstack_networking_network_v2.ext_net.id
  description = "The external network id"
}

output "loc_net_id" {
  value       = openstack_networking_network_v2.loc_net.id
  description = "The local network id"
}

output "loc_subnet_id" {
  value       = openstack_networking_subnet_v2.loc_net_sn.id
  description = "The local subnet id"
}
