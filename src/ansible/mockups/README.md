# Ansible setup for mockups

## Dev or bastion host

    virtualenv venv
    venv/bin/pip install pip-tools
    venv/bin/pip-compile src/python/requirements.in
    venv/bin/pip-compile src/python/requirements-dev.in
    venv/bin/pip install -r src/python/requirements-dev.txt
    cp ansible_sample.cfg to ansible.cfg
    . venv/bin/activate

## How to run the K3s HA cluster sample on AWS

    # create the ansible inventory
    export AWS_ACCESS_KEY_ID=this
    export AWS_SECRET_ACCESS_KEY=that
    export AWS_DEFAULT_REGION=eu-west-3
    venv/bin/python src/python/otvl/k3s_ha_aws2ans.py /path/to/inventory/hosts.yml
    # quick setup bastion host from development host, copy dev env
    ansible-playbook -i /path/to/inventory src/ansible/mockups/mockups_bastion.yml
    # from bastion host with access to private networks
    venv/bin/python src/python/otvl/k3s_ha_aws2ans.py /path/to/inventory/hosts.yml -lb
    ansible-playbook -i /path/to/inventory src/ansible/mockups/mockup_k3s_ha.yml

## How to create a K8s development environment on a single VM

Ansible inventory for a host `k3allhost`

    all:
      hosts:
        dzmk3a:
      children:
        k3all_group:
          hosts:
            dzmk3a:

Run

    ansible-playbook -i /path/to/inventory src/ansible/mockups/mockups_k3all.yml

This will install a VM with:

- this git repository with its venv
- kubectl, Helm, opentofu, nerdctl, buildkit
- K3s single node
- Golang, kubebuider
