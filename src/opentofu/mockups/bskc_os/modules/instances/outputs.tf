output "ports" {
  value = [
    data.openstack_networking_port_v2.ext,
    data.openstack_networking_port_v2.loc
  ]
  description = "The instances ports"
}
output "ip_v4" {
  value = local.ip_v4_addresses
  description = "The instances IP v4 addresses"
}
output "mac" {
  value = local.mac_addresses
  description = "The instances MAC addresses"
}