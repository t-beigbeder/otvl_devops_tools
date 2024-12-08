# ---------------------------------------------------------------------------------------------------------------------
# REQUIRED PARAMETERS
# You must provide a value for each of these parameters.
# ---------------------------------------------------------------------------------------------------------------------
variable "ext_net_name" {
  description = "The name of the external network"
  type        = string
}
variable "loc_net_name" {
  description = "The name of the local network"
  type        = string
}
variable "loc_net_cidr" {
  description = "The CIDR of the local network"
  type        = string
}
variable "bastion_sg_name" {
  description = "The security group name for bastion access"
  type        = string
}
variable "ext_sg_name" {
  description = "The security group name for external access"
  type        = string
}
