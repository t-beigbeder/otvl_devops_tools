output "ext_net_id" {
  value = module.networking.ext_net_id
  description = "The external network id"
}
output "loc_net_id" {
  value = module.networking.loc_net_id
  description = "The local network id"
}
output "ports" {
  value = module.instances.ports
  description = "The ports"
}
output "ip_v4" {
  value = module.instances.ip_v4
  description = "The instances IP v4 addresses"
}
output "mac" {
  value = module.instances.mac
  description = "The instances MAC addresses"
}
output "secrets" {
  value = nonsensitive(yamldecode(data.sops_file.bskc_os_secrets.raw))
}