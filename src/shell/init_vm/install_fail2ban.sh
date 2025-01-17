#!/bin/sh

## pre
rp=`realpath $0`
sd=`dirname $rp`
. $sd/env_install.sh
## endpre

install_fail2ban() {
  is_root && \
  apt-get install -y --no-install-recommends fail2ban && \
  true
  return $?
}