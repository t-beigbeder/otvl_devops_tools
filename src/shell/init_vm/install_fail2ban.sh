#!/bin/sh

## pre
rp=`realpath $0`
sd=`dirname $rp`
. $sd/env_install.sh
## endpre

patch_fail2ban_install() {
  cat > /etc/fail2ban/jail.d/defaults-debian.conf <<EOF
[DEFAULT]
# Debian 12 has no log files, just journalctl
backend = systemd

[sshd]
enabled = true
EOF

}

install_fail2ban() {
  is_root && \
  apt-get install -y --no-install-recommends fail2ban && \
  patch_fail2ban_install && \
  systemctl restart fail2ban && \
  systemctl status fail2ban && \
  true
  return $?
}

if [ ! -f /usr/bin/fail2ban-server ] ; then
  install_fail2ban
fi