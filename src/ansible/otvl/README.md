# Ansible setup for otvl_web

    virtualenv venv
    venv/bin/pip install pip-tools
    venv/bin/pip-compile src/python/requirements.in
    venv/bin/pip-compile src/python/requirements-dev.in
    venv/bin/pip install -r src/python/requirements-dev.txt
    cp ansible_sample.cfg to ansible.cfg
    . venv/bin/activate
    # quick setup bastion host from development host, copy dev env
    ansible-playbook -i /path/to/inventory src/ansible/otvl/otvl_bastion_v6.yml
    # from bastion host with fast access to openstack API
    ansible-playbook -i /path/to/inventory src/ansible/otvl/otvl_web_v6.yml

## Documentation

- [kubernetes.core.helm module](https://docs.ansible.com/ansible/latest/collections/kubernetes/core/helm_module.html)
