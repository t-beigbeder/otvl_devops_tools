# Query OpenStack to generate an otvl_web inventory once
# usage: python otvl_os2ans.py <os-cloud> <yaml-inventory-host-file> [-lb]
# -lb indicates that bastion is local, meaning local subnet must be used
#
import json
import subprocess
import sys
import yaml


def run():
    res, lb = False, False
    if len(sys.argv) < 3:
        sys.stderr.write(f"Usage: {sys.argv[0]} <os-cloud> <yaml-inventory-host-file> [-lb]\n")
        return res
    if len(sys.argv) == 4:
        if sys.argv[3] != "-lb":
            sys.stderr.write(f"Usage: {sys.argv[0]} <os-cloud> <yaml-inventory-host-file> [-lb]\n")
            return res
        lb = True
    result = subprocess.run(f"openstack --os-cloud {sys.argv[1]} server list -f json --long".split(" "), stdout=subprocess.PIPE)
    if result.returncode != 0:
        return res
    osi = json.loads(result.stdout)
    servers = {}
    for server in osi:
        name = server["Name"]
        ip = None
        groups = []
        for nn, network in server["Networks"].items():
            for lip in network:
                ipfs = lip.split(".")
                if len(ipfs) == 4:
                    if lb and ipfs[0] == "172":
                        ip = lip
                    elif not lb and ipfs[0] != "172":
                        ip = lip
        for np, prop in server["Properties"].items():
            if np == "groups":
                groups = prop.split(",")
        servers[name] = {"ip": ip, "groups": groups}
    yo = {"all": {"hosts": {}, "children": {}}}
    for sn, sv in servers.items():
        yo["all"]["hosts"][sn] = {"ansible_host": sv["ip"]}
        for group in sv["groups"]:
            if group not in yo["all"]["children"]:
                yo["all"]["children"][group] = {"hosts": {}}
            yo["all"]["children"][group]["hosts"][sn] = None
    with open(sys.argv[2], "w") as osf:
        yaml.safe_dump(yo, osf)
    res = True
    return res


if __name__ == '__main__':
    gres = run()
    sys.exit(0 if gres else -1)
