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
variable "sproxy_sg_name" {
  description = "The security group name for external access"
  type        = string
}
variable "sproxy_port" {
  description = "The UDP port of bssms proxy"
  type = string
}
variable "sproxy_ssh_exposed" {
  description = "does the sproxy server expose SSH"
  type = bool
  default = false
}
