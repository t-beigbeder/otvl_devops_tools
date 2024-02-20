output "ext_net_id" {
  value = module.networking.ext_net_id
  description = "The external network id"
}
output "loc_net_id" {
  value = module.networking.loc_net_id
  description = "The local network id"
}
output "ext_ports" {
  value = module.instances.ext_ports
  description = "The ext ports"
}
