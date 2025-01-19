output "ids" {
  value = [
    for i in openstack_compute_instance_v2.instances : i.id
  ]
}
output "ipv4s" {
  value = [
    for i in openstack_compute_instance_v2.instances : [
      for n in i.network : n.fixed_ip_v4
    ]
  ]
}
output "macs" {
  value = [
    for i in openstack_compute_instance_v2.instances :[
      for n in i.network : n.mac
    ]
  ]
}
