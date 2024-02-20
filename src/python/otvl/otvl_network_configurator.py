import sys
import logging
import os
import json
import subprocess
import re
import copy
import traceback

import pyudev
import yaml


class NetConf:
    logger = logging.getLogger(__module__ + '.' + __qualname__)  # noqa

    def __init__(self, name):
        self.logger.debug("__init__")
        self.name = name
        self.previous_state = None
        self.current_state = None
        self.stf = "{0}/otvl_network_configurator.json".format(os.getenv("OTVL_NC_STD", "/srv/otvl/iaas/data"))
        self.nif = os.getenv("OTVL_NC_NIF", "/etc/network/interfaces")
        self.previous_nis = None
        self.current_nis = None
        self.fake_ifup = os.getenv("OTVL_NC_FAKEIFUP")
        self.rcr = os.getenv("OTVL_NC_RCR", "/srv/otvl/iaas/data/resolv.conf.reference")
        self.rcf = os.getenv("OTVL_NC_RCF", "/etc/resolv.conf")
        self.if_changed = False
        self.force_if_changed = os.getenv("OTVL_NC_FIC")
        self.ncf = os.getenv("OTVL_NC_NCF", "/srv/otvl/iaas/config/network_config.yml")
        with open(self.ncf) as fncf:
            self.network_config = yaml.load(fncf, yaml.FullLoader)

    def _load_state(self):
        try:
            with open(self.stf, encoding="utf-8") as fp:
                self.previous_state = json.load(fp)
        except FileNotFoundError:
            self.previous_state = {"devices": {}}
        self.current_state = copy.deepcopy(self.previous_state)

    def _save_state(self):
        with open(self.stf, "w", encoding="utf-8") as fp:
            json.dump(self.current_state, fp, indent=2)

    def _load_nif(self):
        self.logger.debug(f"_load_nif nif {self.nif}")
        self.previous_nis = {}
        try:
            with open(self.nif, encoding="utf-8") as fp:
                current_iface = ""
                for line in fp:
                    if line[-1] == "\n":
                        line = line[0:-1]
                    words = re.split(r"\s+", line)
                    self.logger.debug(f"_load_nif words {words}")
                    if len(words) >= 2 and words[0] == "allow-hotplug":
                        # allow-hotplug enp7s0
                        self.previous_nis[words[1]] = {}
                    elif len(words) >= 4 and words[0] == "iface":
                        # iface enp7s0 inet dhcp
                        # iface eth0 inet static
                        try:
                            self.previous_nis[words[1]]["iface"] = " ".join(words)
                            self.previous_nis[words[1]]["ip"] = ""
                            current_iface = words[1] if words[3] == "static" else ""
                        except KeyError:
                            pass
                    elif len(words) >= 2 and words[0] == "address" and current_iface:
                        # address 192.0.2.7/24
                        self.previous_nis[current_iface]["ip"] = words[2]
        except FileNotFoundError:
            pass
        self.current_nis = copy.deepcopy(self.previous_nis)

    def _update_nif(self):
        self.logger.debug(f"_update_nif nif {self.nif}")
        for itf in self.current_state["devices"].keys():
            if itf not in self.current_nis:
                self.logger.debug(f"_update_nif misses {itf}")
                break
            self.logger.debug(f"_update_nif already {itf}")
        else:
            return
        with open(self.nif, "w", encoding="utf-8") as fp:
            index = 0
            for itf in self.current_state["devices"].keys():
                network = self.network_config["networks"][index]
                index += 1
                self.logger.debug(f"_update_nif append {itf} for {network}")
                if "cidr" in network and network["cidr"]:
                    net_mask = network["cidr"].split("/")[1]
                else:
                    net_mask = ""
                host_ip = "{0}/{1}".format(network["host_ip"], net_mask) if network["host_ip"] else "dhcp"
                self.logger.info(f"Declared interface {itf} in file {self.nif} for IP {host_ip}")
                fp.write(f"\n# added by otvl_network_configurator\nallow-hotplug {itf}\n")
                if host_ip == "dhcp":
                    fp.write(f"iface {itf} inet dhcp\n")
                else:
                    fp.write(f"iface {itf} inet static\n address {host_ip}\n")

    def _ip_for(self, dev_name):
        cmd = f"ip -4 -o address show dev {dev_name}"
        line = subprocess.run(cmd.split(" "), stdout=subprocess.PIPE).stdout.decode("utf-8")
        try:
            words = re.split(r"\s+", line)
            ip = words[3].split("/")[0]
        except IndexError:
            ip = ""
        return ip

    def _run_cmd(self, args, warn_if_error=True):
        if type(args) is str:
            args = args.split(" ")
        self.logger.info(f"Running {args}")
        result_cmd = subprocess.run(args, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        self.logger.debug(f"_run_cmd {args} {result_cmd}")
        if result_cmd.returncode != 0:
            diag = f"command {result_cmd.args} returned {result_cmd.returncode} error {result_cmd.stderr.decode('utf-8')}"  # noqa
            if warn_if_error:
                self.logger.warning(diag)
                return True, result_cmd.stdout.decode('utf-8')
            else:
                self.logger.error(diag)
                return False, result_cmd.stdout.decode('utf-8')
        return True, result_cmd.stdout.decode('utf-8')

    def _ifup_new(self):
        for itf in self.current_state["devices"].keys():
            if itf in self.current_nis:
                self.logger.debug(f"_ifup_new already in {self.nif} {itf}")
                continue
            if self.current_state["devices"][itf]:
                self.logger.debug(f"_ifup_new already has ip {itf}")
                continue
            if self.fake_ifup:
                cmd = ["sh", "-c", f'ifup {itf}']
            else:
                cmd = ["ifup", itf]
            self._run_cmd(cmd)
            self.if_changed = True

    def _change_rcf(self):
        with open(self.rcf, "r") as fpf:
            with open(self.rcr, "r") as fpr:
                ref = fpr.read()
                cur = fpf.read()
        if ref == cur:
            return
        with open(self.rcf, "w") as fpf:
            self.logger.info(f"Updating {self.rcf} from {self.rcr}")
            fpf.write(ref)

    def _run_fakeip(self, args):
        if self.fake_ifup:
            args = ["sh", "-c", "fakeip " + args]
        return self._run_cmd(args, warn_if_error=False)

    def _display_routes(self, json_lines):
        routes = json.loads(json_lines)
        self.logger.info(f"Current routes are:")
        for route in routes:
            self.logger.info(f"\t{route}")

    def _ip_route_change(self):
        self.logger.debug("_ip_route_change: start")

        # ip route show
        lines = subprocess.run(["ip", "-j", "route", "show"], stdout=subprocess.PIPE).stdout.decode("utf-8")
        # default via 51.83.104.1 dev eth0
        default_ip = None
        for route in json.loads(lines):
            if route["dst"] != "default":
                continue
            self.logger.debug(f"_ip_route_change: default {route}")
            # get IP for eth0
            default_ip = self.current_state["devices"][route['dev']]
            break

        # identify network for this IP
        for nv in self.network_config["networks"]:
            if nv["prefix"] is None:
                continue
            if default_ip and default_ip.startswith(nv["prefix"]):
                route_network = nv
                break
        else:
            route_network = None

        # if this network is not external route del
        if route_network is None or route_network["key"] == "external":
            if not self.force_if_changed:
                return
        self._display_routes(lines)
        args = "ip route del default"
        if self.fake_ifup:
            args = ["sh", "-c", "fakeip " + args]
        self._run_cmd(args, warn_if_error=False)

        # ifdown/ifup
        for dev, ip in self.current_state["devices"].items():
            for nv in self.network_config["networks"]:
                if nv["prefix"] is None:
                    continue
                if ip and ip.startswith(nv["prefix"]):
                    dev_network = nv
                    break
            else:
                dev_network = None
            if dev_network and dev_network["key"] != "external":
                self.logger.debug(f"_ip_route_change: internal {dev_network} for {dev} {ip}")
                continue
            self.logger.debug(f"_ip_route_change: external or undefined {dev_network} for {dev} {ip}")
            self._run_fakeip(f"ifdown {dev}")
            self._run_fakeip(f"ifup {dev}")

        # ip route show
        res, lines = self._run_cmd("ip -j route show")
        self._display_routes(lines)

    def _ip_route_check_default(self):
        self.logger.debug("_ip_route_check_default: start")

        # ip route show
        lines = subprocess.run(["ip", "-j", "route", "show"], stdout=subprocess.PIPE).stdout.decode("utf-8")
        # default via 51.83.104.1 dev eth0
        default_ip = None
        default_route = None
        for route in json.loads(lines):
            if route["dst"] != "default":
                continue
            self.logger.debug(f"_ip_route_check_default: default {route}")
            default_route = route
            default_ip = self.current_state["devices"][route['dev']]
            break

        # identify network for this IP
        for nv in self.network_config["networks"]:
            if nv["prefix"] is None:
                continue
            if default_ip and default_ip.startswith(nv["prefix"]):
                route_network = nv
                break
        else:
            route_network = None

        if default_route and (route_network is None or route_network["key"] == "external"):
            return
        self.logger.info(f"_ip_route_check_default default_route {default_route} route_network {route_network}")
        self._display_routes(lines)

        # ifdown/ifup
        for dev, ip in self.current_state["devices"].items():
            for nv in self.network_config["networks"]:
                if nv["prefix"] is None:
                    continue
                if ip and ip.startswith(nv["prefix"]):
                    dev_network = nv
                    break
            else:
                dev_network = None
            if dev_network and dev_network["key"] != "external":
                self.logger.debug(f"_ip_route_change: internal {dev_network} for {dev} {ip}")
                continue
            self.logger.debug(f"_ip_route_change: external or undefined {dev_network} for {dev} {ip}")
            self._run_fakeip(f"ifdown {dev}")
            self._run_fakeip(f"ifup {dev}")

        # ip route show
        res, lines = self._run_cmd("ip -j route show")
        self._display_routes(lines)

    def _do_it(self):
        self._load_state()
        context = pyudev.Context()
        for device in context.list_devices(subsystem='net'):
            if device.properties["DEVPATH"].startswith("/devices/virtual/net/"):
                continue
            self.logger.debug(f"listing {device} {device.properties['DEVPATH']}")
            itf = device.properties["INTERFACE"]
            self.current_state["devices"][itf] = self._ip_for(itf)

        self._load_nif()
        self._update_nif()
        self._ifup_new()
        self._change_rcf()

        if self.if_changed or self.force_if_changed:
            self._ip_route_change()

        self._ip_route_check_default()

        self._save_state()
        return True

    def run(self):
        try:
            self.logger.debug("starting")
            result = self._do_it()
            self.logger.debug("exiting")
            return result
        except Exception as e:
            self.logger.error(traceback.format_exc())
            self.logger.error(
                'An unkonwn error occured, please contact the support - {0} {1}'.format(
                    type(e), e))
        return False


def _logging_setup():
    logging.basicConfig(
        level=os.getenv('OTVL_NC_LOGGING', 'INFO'),
        format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')
    bl = logging.getLogger()
    lfh = logging.FileHandler(os.getenv('OTVL_NC_LOGFILE', '/var/log/otvl_nc.log'), encoding="utf-8")
    lfh.setLevel(bl.level)
    lfh.setFormatter(logging.Formatter('%(asctime)s - %(name)s - %(levelname)s - %(message)s'))
    bl.addHandler(lfh)


if __name__ == '__main__':
    _logging_setup()
    cmd_name = os.path.basename(sys.argv[0])
    res = NetConf(cmd_name).run()
    sys.exit(0 if res else -1)
