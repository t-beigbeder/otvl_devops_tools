output "ports" {
  value = [
    data.openstack_networking_port_v2.ext,
    data.openstack_networking_port_v2.loc
  ]
  description = "The instance ports"
}
