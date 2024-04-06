variable "subnets_name_filter" {
  description = "Some information to filter the subnets by their name"
  type        = string
}

variable "vpc_is_default" {
  description = "Use default VPC, else only one VPC should be available"
  type        = bool
}
