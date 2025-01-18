#!/bin/sh

## pre
rp=`realpath $0`
sd=`dirname $rp`
. $sd/env_install.sh
## endpre

install_fail2ban() {
  is_root && \
  apt-get install -y --no-install-recommends fail2ban && \
  systemctl status fail2ban && \
  true
  return $?
}

if [ ! -f /usr/bin/fail2ban-server ] ; then
  install_fail2ban
fi