
terraform {
  required_version = ">= 1.6.2, < 2.0.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.44"
    }
  }
}

data "aws_ami" "this" {
  most_recent = true
  name_regex  = var.ami_name_regex
  filter {
    name   = "architecture"
    values = ["x86_64"]
  }
  owners = var.ami_owner == "" ? null : [var.ami_owner]
}
