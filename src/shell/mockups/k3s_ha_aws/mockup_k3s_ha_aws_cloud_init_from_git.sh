#!/bin/sh

enableswap() {
  if [ -f /swapfile ] ; then return 0 ; fi
  fallocate -l 1G /swapfile && \
  chmod 600 /swapfile && \
  mkswap /swapfile && \
  echo "/swapfile swap swap defaults 0 0" >> /etc/fstab && \
  true
}

echo `date`: command $0 is starting

cat > /etc/fail2ban/jail.d/defaults-debian.conf <<EOF
[DEFAULT]
# Debian 12 has no log files, just journalctl
backend = systemd

[sshd]
enabled = true
EOF

enableswap && \
systemctl restart fail2ban.service && \
true || exit 1

exit 0
