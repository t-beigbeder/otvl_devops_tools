# OpenTofu setup for otvl_web

## Tool setup

### Openstack

See ansible
[setup](../../ansible/otvl/README.md)
concerning openstack client installation.

Set OS_* env variables according to
[documentation](https://docs.openstack.org/python-openstackclient/3.4.1/man/openstack.html#environment-variables),
then check:

    openstack flavor list

### OpenTofu

Avoid storing 300 MB for aws plugin (if applicable) per tf launch

    $ vi ~/.tofurc
    plugin_cache_dir = "$HOME/.terraform.d/plugin-cache"

## Deploy using IaC

    cd src/opentofu/otvl/prod
    cp values.auto.tfvars.sample values.auto.tfvars
    vi values.auto.tfvars.sample
    tofu init
    tofu apply

## Help

### General

- [GitHub terraform up and running](https://github.com/brikis98/terraform-up-and-running-code)
- [h2 use tf variables](https://spacelift.io/blog/how-to-use-terraform-variables)
- [Terraform Variables - A Standard](https://lachlanwhite.com/posts/terraform/10-11-2021-terraform-variables-a-standard/)

### Specific cases

- [Destroy all but protected resource](https://stackoverflow.com/questions/55265203/terraform-delete-all-resources-except-one)

Eg

    tofu state list
    tofu state rm module.prereqs.aws_kms_key.kms_key_for_infra
    tofu destroy

## Documentation

- [Open Tofu documentation](https://opentofu.org/docs/cli/)
- [Terraform](https://developer.hashicorp.com/terraform?product_intent=terraform)
- [AWS Provider](https://registry.terraform.io/providers/hashicorp/aws/latest/docs)
- [OS Provider](https://registry.terraform.io/providers/terraform-provider-openstack/openstack/latest/docs)
