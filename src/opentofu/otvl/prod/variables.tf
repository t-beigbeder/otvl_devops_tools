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
variable "ssh_key_name" {
  description = "The SSH key name to store the ssh_pub public key"
  type        = string
}
variable "ssh_pub" {
  description = "The SSH public key to authorize in created instances"
  type        = string
}
variable "instances_attrs" {
  description = "Attributes for instances to create"
  type        = list(object({
    name         = string
    logical_name = string
    groups       = list(string)
    ip_v4 = string
    image_name   = string
    flavor_name  = string
  }))
}

# ---------------------------------------------------------------------------------------------------------------------
# OPTIONAL PARAMETERS
# These parameters have reasonable defaults.
# ---------------------------------------------------------------------------------------------------------------------
