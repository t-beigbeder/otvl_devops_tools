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
variable "hosting_sg_name" {
  description = "The security group name for external access"
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
variable "b64_id_rsa_rops" {
  description = "Encrypted private key for remote operations, base64 encoded"
  type = string
}
variable "sproxy_ext_address" {
  description = "The external address of the bssms proxy, hostname or IP"
  type = string
}
variable "sproxy_int_address" {
  description = "The internal address of the bssms proxy, hostname or IP"
  type = string
}
variable "bssms_proxy_port" {
  description = "The UDP port of bssms proxy"
  type = string
}
variable "dot_repo" {
  description = "Git repo devopstools"
  type        = string
}
variable "dot_branch" {
  description = "Git branch devopstools"
  type        = string
}
variable "rops_repo" {
  description = "Git repo remote operations"
  type        = string
}
variable "install_env" {
  description = "Environment to install"
  type        = string
}
variable "go_version" {
  description = "Version of the go runtime"
  type        = string
}
variable "instances_attrs" {
  description = "Attributes for instances to create"
  type = list(object({
    name        = string
    groups      = string
    otvl_meta   = string
    ip_v4       = string
    image_name  = string
    flavor_name = string
  }))
}

# ---------------------------------------------------------------------------------------------------------------------
# OPTIONAL PARAMETERS
# These parameters have reasonable defaults.
# ---------------------------------------------------------------------------------------------------------------------
variable "hosting_secrets_sops" {
  description = "sops enc.yaml containing dictionary of key/value per instance name"
  type = string
  default = "instances_secrets.enc.yaml"
}
