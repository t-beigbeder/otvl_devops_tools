#!/bin/sh

enableswap() {
  if [ -f /swapfile ] ; then return 0 ; fi
  fallocate -l 1G /swapfile && \
  chmod 600 /swapfile && \
  mkswap /swapfile && \
  echo "/swapfile swap swap defaults 0 0" >> /etc/fstab && \
  true
}

enablevenv() {
  if [ -d /srv/venv/otvl_cloud_init ] ; then return 0 ; fi
  virtualenv -p python3 /srv/venv/otvl_cloud_init && \
  /srv/venv/otvl_cloud_init/bin/pip install pyudev PyYAML cryptography && \
  true
}

getosmeta() {
  if [ -f /root/clinit/osmeta.json ] ; then return 0 ; fi
  curl http://169.254.169.254/openstack/2018-08-27/meta_data.json | jq .meta > /root/clinit/osmeta.json && \
  true
}

git_repo="https://github.com/t-beigbeder/otvl_devops_tools"
git_branch="bdev1"
git_local="/root/clinit/otvl_devops_tools"
echo `date`: command $0 is starting

mkdir -p /root/bin /srv/venv /srv/otvl/iaas/python
cat > /root/bin/otvl_display_net_conf.sh <<EOF
#!/bin/sh
logger -t  otvl_display_net_conf -s "command $0 is starting"
echo "otvl_display_net_conf `date`": command $0 is starting > /dev/console
cat /etc/network/interfaces > /tmp/tf.$$
ls -l /etc/network/interfaces.d/50-cloud-init.cfg >> /tmp/tf.$$
rm -f /etc/network/interfaces.d/50-cloud-init.cfg
ip ad show >> /tmp/tf.$$
ip route show >> /tmp/tf.$$
logger -t  otvl_display_net_conf -s -f /tmp/tf.$$
cat /tmp/tf.$$ | sed -e "s/^/otvl_display_net_conf `date` /" > /dev/console
rm /tmp/tf.$$
logger -t  otvl_display_net_conf -s "command $0 is exiting"
echo "otvl_display_net_conf `date`": command $0 is exiting > /dev/console
exit 0
EOF
chmod 700 /root/bin/otvl_display_net_conf.sh

cat > /etc/systemd/system/otvl_display_net_conf.service <<EOF
[Unit]
Description=Display network configuration at startup

[Service]
Type=oneshot
ExecStart=/root/bin/otvl_display_net_conf.sh

[Install]
WantedBy=multi-user.target
EOF

cat > /root/otvl_cloud_init_py_check.py <<EOF
import sys
import pyudev
context = pyudev.Context()
for device in context.list_devices():
    sys.stderr.write("otvl_cloud_init_py_check: listing {0}\n".format(device))
EOF

cat > /etc/systemd/system/otvl_network_configurator.service <<EOF
[Unit]
Description=Runs otvl_network_configurator
Wants=otvl_network_configurator.timer

[Service]
ExecStart=/srv/venv/otvl_cloud_init/bin/python /srv/otvl/iaas/python/otvl_network_configurator.py
WorkingDirectory=/srv/otvl/iaas/python

[Install]
WantedBy=multi-user.target
EOF

cat > /etc/systemd/system/otvl_network_configurator.timer <<EOF
[Unit]
Description=Run otvl_network_configurator every 2 minutes
Requires=otvl_network_configurator.service

[Timer]
Unit=otvl_network_configurator.service
OnCalendar=*:0/2

[Install]
WantedBy=timers.target
EOF

tmp=`ip -4 -o address show | grep dynamic`
external_ip=`echo $tmp | cut -d' ' -f4 | cut -d/ -f1`
nic_dev=`echo $tmp | cut -d' ' -f2`
cat > /etc/network/interfaces <<EOF
# This file describes the network interfaces available on your system
# and how to activate them. For more information, see interfaces(5).

source /etc/network/interfaces.d/*

# The loopback network interface
auto lo
iface lo inet loopback

# The primary network interface
allow-hotplug $nic_dev
iface $nic_dev inet dhcp
EOF

sed -i -e 's=#precedence ::ffff:0:0/96  100=precedence ::ffff:0:0/96  100=' /etc/gai.conf

systemctl enable /etc/systemd/system/otvl_display_net_conf.service && \
virtualenv -p python3 /srv/venv/otvl_cloud_init && \
/srv/venv/otvl_cloud_init/bin/pip install pyudev && \
/srv/venv/otvl_cloud_init/bin/python /root/otvl_cloud_init_py_check.py && \
enableswap && \
enablevenv && \
cp src/python/otvl/otvl_network_configurator.py /srv/otvl/iaas/python/ && \
getosmeta && \
true || exit 1

echo `date`: command $0 is exiting || exit 1
exit 0
