output "bastion_ip" {
  value = aws_instance.k3s_ha_bastion_instance.public_ip
}

output "k3s_server_ips" {
  value = local.aws_instance_ips
}