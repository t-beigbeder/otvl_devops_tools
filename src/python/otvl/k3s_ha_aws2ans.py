# Query AWS to generate an ansible inventory once
# usage: python aws2ans.py <yaml-inventory-host-file> [-lb]
# -lb indicates that bastion is local, meaning private ip adresses subnet must be used
#
import json
import subprocess
import sys
import yaml


def run():
    debug = True

    res, lb = False, False
    if len(sys.argv) < 2:
        sys.stderr.write(f"Usage: {sys.argv[0]} <yaml-inventory-host-file> [-lb]\n")
        return res
    if len(sys.argv) == 3:
        if sys.argv[2] != "-lb":
            sys.stderr.write(f"Usage: {sys.argv[0]} <yaml-inventory-host-file> [-lb]\n")
            return res
        lb = True
    if debug:
        with open("/home/guest/locgit/otvl/otvl_devops_tools/tmp/ec2-dci.json") as fi:
            osi = json.load(fi)
    else:
        result = subprocess.run(f"openstack --os-cloud {sys.argv[1]} server list -f json --long".split(" "), stdout=subprocess.PIPE)
        if result.returncode != 0:
            return res
        osi = json.loads(result.stdout)
    servers = {}
    for rsv in osi["Reservations"]:
        for inst in rsv["Instances"]:
            for tag in inst["Tags"]:
                if tag["Key"] == "Name":
                    name = tag["Value"]
                    if "k3s-ha-bastion" in name:
                        groups = ["bastion_group"]
                    elif "k3s-ha-server" in name:
                        groups = ["bastion_controlled_group", "k3s_ha_server_group"]
                    elif "k3s-ha-node" in name:
                        groups = ["bastion_controlled_group", "k3s_ha_node_group"]
                    else:
                        groups = []
                    servers[name] = servers[name] = {"ip": inst["PrivateIpAddress"] if lb else inst["PublicIpAddress"], "groups": groups}
    yo = {"all": {"hosts": {}, "children": {}}}
    for sn, sv in servers.items():
        yo["all"]["hosts"][sn] = {"ansible_host": sv["ip"]}
        for group in sv["groups"]:
            if group not in yo["all"]["children"]:
                yo["all"]["children"][group] = {"hosts": {}}
            yo["all"]["children"][group]["hosts"][sn] = None
    with open(sys.argv[1], "w") as osf:
        yaml.safe_dump(yo, osf)
    res = True
    return res


if __name__ == '__main__':
    gres = run()
    sys.exit(0 if gres else -1)
