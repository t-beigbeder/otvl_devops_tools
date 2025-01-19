#!/bin/sh

## pre
rp=`realpath $0`
sd=`dirname $rp`
. $sd/env_install.sh
## endpre

patch_fail2ban_install() {
  cat > /etc/fail2ban/jail.local <<EOF
[DEFAULT]
# Debian 12 has no log files, just journalctl
backend = systemd

# "bantime" is the number of seconds that a host is banned.
bantime  = 1d
# "maxretry" is the number of failures before a host get banned.
maxretry = 5
# A host is banned if it has generated "maxretry" during the last "findtime"
findtime  = 1h

[sshd]
enabled = true

EOF

}

install_fail2ban() {
  is_root && \
  apt-get install -y --no-install-recommends fail2ban python3-systemd && \
  patch_fail2ban_install && \
  systemctl restart fail2ban && \
  systemctl status fail2ban && \
  true
  return $?
}

if [ ! -f /usr/bin/fail2ban-server ] ; then
  install_fail2ban
fi
