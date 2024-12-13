output "ext_ports" {
  value = data.openstack_networking_port_v2.ext
  description = "The ext ports"
}

output "instances" {
    value = resource.openstack_compute_instance_v2.this
    description = "The instances"
}