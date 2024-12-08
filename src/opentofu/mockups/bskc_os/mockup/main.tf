provider "openstack" {
  alias = "osdp"
}

terraform {
  required_version = ">= 1.6.2, < 2.0.0"

  required_providers {
    openstack = {
      source  = "terraform-provider-openstack/openstack"
      version = "~> 1.42.0"
    }
  }

  /*
    backend "s3" {
      bucket         = "default-tf-bucket"
      key            = "otvl/prod/terraform.tfstate"
      region         = "here"
    }
  */

}
