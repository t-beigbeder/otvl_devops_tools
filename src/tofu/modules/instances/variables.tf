# ---------------------------------------------------------------------------------------------------------------------
# REQUIRED PARAMETERS
# You must provide a value for each of these parameters.
# ---------------------------------------------------------------------------------------------------------------------
variable "ext_net_id" {
  description = "The external network id"
  type        = string
}
variable "loc_net_id" {
  description = "The local network id"
  type        = string
}
variable "loc_subnet_id" {
  description = "The local subnet id"
  type        = string
}
variable "external_sg_id" {
  description = "external sg id"
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
variable "user_data_template" {
  description = "Template file for user data to be provided to the instance"
  type        = string
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
  description = "Golang version to download for building bssms"
  type        = string
}
variable "sproxy_int_address" {
  description = "The internal address of the bssms proxy, hostname or IP"
  type = string
}
variable "bssms_proxy_port" {
  description = "QUIC port of the bssms-proxy"
  type        = string
}
variable "instances_attrs" {
  description = "Attributes for instances to create"
  type = list(object({
    name            = string
    groups          = string
    otvl_meta       = string
    ip_v4           = string
    image_name      = string
    flavor_name     = string
    secrets_pri_key = string
  }))
}
