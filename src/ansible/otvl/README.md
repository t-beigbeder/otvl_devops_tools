# Ansible setup for otvl_web

    virtualenv venv
    venv/bin/pip install pip-tools
    venv/bin/pip-compile src/python/requirements.in
    venv/bin/pip-compile src/python/requirements-dev.in
    venv/bin/pip install -r src/python/requirements-dev.txt
    cp ansible_sample.cfg to ansible.cfg
    . venv/bin/activate
    # quick setup bastion host from dev host
    ansible-playbook -i /path/to/inventory src/ansible/otvl/otvl_bastion_v6.yml
    # from bastion with fast access to openstack API
    ansible-playbook -i /path/to/inventory src/ansible/otvl/otvl_web_v6.yml
