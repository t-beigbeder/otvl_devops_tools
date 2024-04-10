# ---------------------------------------------------------------------------------------------------------------------
# REQUIRED PARAMETERS
# You must provide a value for each of these parameters.
# ---------------------------------------------------------------------------------------------------------------------

variable "ec2_bastion_instance_type" {
  description = "The type of EC2 Bastion Instance to run (e.g. t3.micro)"
  type        = string
}

variable "ec2_instance_ami_name_regex" {
  description = "The name regex for filtering AMI of EC2 Instances to run (e.g. amzn2-ami-amd or debian-12)"
  type        = string
}

variable "ec2_instance_ami_owner" {
  description = "The owner name or empty for filtering AMI of EC2 Instances to run (e.g. 099720109477 or 136693071363)"
  type        = string
}

variable "ec2_k3s_server_nb_per_subnet" {
  description = "Number of EC2 instances per subnet for K3s server nodes"
}

variable "ec2_k3s_server_instance_type" {
  description = "The type of EC2 Instance for K3s server nodes to run (e.g. t3.micro)"
  type        = string
}

# ---------------------------------------------------------------------------------------------------------------------
# OPTIONAL PARAMETERS
# These parameters have reasonable defaults.
# ---------------------------------------------------------------------------------------------------------------------

variable "ec2_bastion_instance_key_name" {
  description = "The key name for ssh to EC2 Bastion Instance or empty if no ssh"
  type        = string
  default     = ""
}

variable "ec2_instance_user_data" {
  description = "User data to be passed to the instance"
  type        = string
  default     = null
}
