output "ids" {
  value       = local.ids
  description = "The list of filtered subnets ids"
}

output "cidrs" {
  value       = local.cidrs
  description = "The list of filtered subnets CIDRs"
}

output "list" {
  value       = local.list
  description = "The list of filtered subnets"
}

output "byid" {
  value       = local.byid
  description = "The map of filtered subnets by id"
}

output "byname" {
  value       = local.byname
  description = "The map of filtered subnets by tag Name"
}

output "default_vpc" {
  value = data.aws_vpc.default
  description = "The default VPC in this account"
}