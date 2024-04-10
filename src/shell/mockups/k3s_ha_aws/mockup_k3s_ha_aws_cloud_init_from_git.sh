#!/bin/sh

enableswap() {
  if [ -f /swapfile ] ; then return 0 ; fi
  fallocate -l 1G /swapfile && \
  chmod 600 /swapfile && \
  mkswap /swapfile && \
  echo "/swapfile swap swap defaults 0 0" >> /etc/fstab && \
  true
}

git_repo="https://github.com/t-beigbeder/otvl_devops_tools"
git_branch="bdev2"
git_local="/root/clinit/otvl_devops_tools"
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
echo `date`: command $0 is exiting, will reboot in 10s
sleep 10
reboot
exit 0
