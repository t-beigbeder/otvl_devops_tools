#!/bin/sh

enableswap() {
  if [ -f /swapfile ] ; then return 0 ; fi
  fallocate -l 1G /swapfile && \
  chmod 600 /swapfile && \
  mkswap /swapfile && \
  swapon /swapfile && \
  echo "/swapfile swap swap defaults 0 0" >> /etc/fstab && \
  true
}

echo `date`: command $0 is starting
. /etc/k3s_ha_aws/ec2_exports
enableswap || exit 1

if [ "${ec2_profile}" = "k3s-ha-bastion" -a "${ec2_bastion_instance_has_fail2ban}" != "false" ] ; then
  cat > /etc/fail2ban/jail.d/defaults-debian.conf <<EOF
[DEFAULT]
# Debian 12 has no log files, just journalctl
backend = systemd

[sshd]
enabled = true
EOF
  systemctl restart fail2ban.service || exit 1
fi

echo `date`: command $0 is exiting
exit 0
