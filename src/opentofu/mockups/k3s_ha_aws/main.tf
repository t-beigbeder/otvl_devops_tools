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

module "gsn" {
  source = "../../modules/aws/get_subnets"
  subnets_name_filter = ""
  vpc_is_default = true
}
