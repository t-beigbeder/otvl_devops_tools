# Query AWS to generate an ansible inventory once
# usage: python aws2ans.py <yaml-inventory-host-file> [-lb]
# -lb indicates that bastion is local, meaning private ip adresses subnet must be used
#
import json
import subprocess
import sys
import yaml


def run():
    res, lb = False, False
    if len(sys.argv) < 2:
        sys.stderr.write(f"Usage: {sys.argv[0]} <yaml-inventory-host-file> [-lb]\n")
        return res
    if len(sys.argv) == 3:
        if sys.argv[2] != "-lb":
            sys.stderr.write(f"Usage: {sys.argv[0]} <yaml-inventory-host-file> [-lb]\n")
            return res
        lb = True
    result = subprocess.run(f"aws ec2 describe-instances".split(" "), stdout=subprocess.PIPE)
    if result.returncode != 0:
        return res
    osi = json.loads(result.stdout)
    servers = {}
    server_count = 0
    for rsv in osi["Reservations"]:
        for inst in rsv["Instances"]:
            if inst["State"]["Name"] != "running":
                continue
            for tag in inst["Tags"]:
                if tag["Key"] == "Name":
                    name = tag["Value"]
                    if "k3s-ha-bastion" in name:
                        groups = ["bastion_group"]
                    elif "k3s-ha-server" in name:
                        groups = ["bastion_controlled_group", "k3s_ha_server_group"]
                        server_count += 1
                        if server_count == 2:
                            groups.append("k3s_build_group")
                    elif "k3s-ha-node" in name:
                        groups = ["bastion_controlled_group", "k3s_ha_node_group"]
                    else:
                        groups = []
                    servers[name] = servers[name] = {"ip": inst["PrivateIpAddress"] if lb or "PublicIpAddress" not in inst else inst["PublicIpAddress"], "groups": groups}
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
