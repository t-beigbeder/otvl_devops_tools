output "ext_ports" {
  value = data.openstack_networking_port_v2.ext
  description = "The ext ports"
}
