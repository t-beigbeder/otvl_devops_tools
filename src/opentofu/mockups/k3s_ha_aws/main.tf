provider "aws" {
  region = "eu-west-3"
}

terraform {
  required_version = ">= 1.0.0, < 2.0.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.44.0"
    }
  }

  #   backend "s3" {
  #     bucket         = "default-tf-bucket"
  #     key            = "mockups/k3s_ha_aws/terraform.tfstate"
  #     region         = "eu-west-3"
  #     encrypt        = true
  #   }
}

module "get_default_subnets" {
  source              = "../../modules/aws/get_subnets"
  subnets_name_filter = "-default"
  vpc_is_default      = true
}

module "sg_bastion" {
  source         = "../../modules/aws/mk_sg"
  name           = "bastion"
  default_vpc_id = module.get_default_subnets.default_vpc.id
  ingress_rules = [{
    from_port          = 22
    to_port            = 22
    protocol           = "tcp"
    cidr_blocks        = ["0.0.0.0/0"]
    ipv6_cidr_blocks   = ["::/0"]
    security_group_ids = []
  }]
  egress_allow_all = true
  tags             = {}
}

module "get_ami" {
  source         = "../../modules/aws/get_ami"
  ami_name_regex = var.ec2_instance_ami_name_regex
  ami_owner      = var.ec2_instance_ami_owner
}

resource "aws_instance" "bastion_instance" {
  ami           = module.get_ami.ami.id
  instance_type = var.ec2_bastion_instance_type
  key_name      = var.ec2_bastion_instance_key_name
  user_data = base64encode(templatefile("${path.module}/cloud-config.yaml", {
    ec2_git_repo         = var.ec2_git_repo
    ec2_git_branch       = var.ec2_git_branch
    ec2_profile          = "k3s-ha-bastion"
    ec2_hostname         = "k3s-ha-bastion"
    ec2_k3s_server_count = var.ec2_k3s_server_nb_per_subnet * length(module.get_default_subnets.ids)
  }))
  vpc_security_group_ids = [module.sg_bastion.security_group.id]
  tags = {
    Name = "k3s-ha-bastion"
  }
}

module "sg_k3s_server" {
  source         = "../../modules/aws/mk_sg"
  name           = "k3s_server"
  default_vpc_id = module.get_default_subnets.default_vpc.id
  ingress_rules = [{
    from_port          = 22
    to_port            = 22
    protocol           = "tcp"
    cidr_blocks        = []
    ipv6_cidr_blocks   = []
    security_group_ids = [module.sg_bastion.security_group.id]
  }]
  egress_allow_all = true
  tags             = {}
}

resource "aws_instance" "k3s_server_instance" {
  count         = var.ec2_k3s_server_nb_per_subnet * length(module.get_default_subnets.ids)
  ami           = module.get_ami.ami.id
  instance_type = var.ec2_bastion_instance_type
  key_name      = var.ec2_bastion_instance_key_name
  user_data = base64encode(templatefile("${path.module}/cloud-config.yaml", {
    ec2_git_repo         = var.ec2_git_repo
    ec2_git_branch       = var.ec2_git_branch
    ec2_profile          = "k3s-ha-server"
    ec2_hostname         = format("k3s-ha-server-%d", count.index)
    ec2_k3s_server_count = var.ec2_k3s_server_nb_per_subnet * length(module.get_default_subnets.ids)
  }))
  subnet_id              = module.get_default_subnets.ids[count.index % length(module.get_default_subnets.ids)]
  vpc_security_group_ids = [module.sg_k3s_server.security_group.id]
  tags = {
    Name = format("k3s-ha-server-%d", count.index)
  }
}

locals {
  aws_instance_ips = [
    for v in aws_instance.k3s_server_instance : v.public_ip
  ]
}