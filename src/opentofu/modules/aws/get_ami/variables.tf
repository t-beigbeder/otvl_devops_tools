variable "ami_name_regex" {
  description = "The name regex for filtering AMI of EC2 Instances to run (e.g. amzn2-ami-amd)"
  type        = string
}

variable "ami_owner" {
  description = "The owner name or empty for filtering AMI of Instances to run (e.g. 099720109477)"
  type        = string
}
