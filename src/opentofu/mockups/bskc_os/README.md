# Bootstrap for a k8s cluster on openstack

    ~/.tofurc
    plugin_cache_dir = "$HOME/.terraform.d/plugin-cache"
    provider_installation {
    dev_overrides {
    "hashicorp.com/edu/hashicups" = "/home/guest/go/bin"
    "tofu.otvl.org/otvl/bssms" = "/home/guest/go/bin"
    "registry.opentofu.org/carlpett/sops" = "/home/guest/.terraform.d/plugin-cache/registry.opentofu.org/carlpett/sops/1.1.1/linux_amd64"
    "registry.opentofu.org/terraform-provider-openstack/openstack" = "/home/guest/.terraform.d/plugin-cache/registry.opentofu.org/terraform-provider-openstack/openstack/1.42.0/linux_amd64"
    }
    direct {}
    }