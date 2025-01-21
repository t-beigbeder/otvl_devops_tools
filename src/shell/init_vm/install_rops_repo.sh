#!/bin/sh

## pre
rp=`realpath $0`
sd=`dirname $rp`
. $sd/env_install.sh
## endpre

install_rops() {
  cd /root/locgit && \
  GIT_SSH_COMMAND="ssh -o StrictHostKeyChecking=accept-new -i /root/.ssh/id_rsa_gh_ssh" git clone git@github.com:t-beigbeder/otvl_rops.git && \
  true
  return $?

}

if [ ! -d /root/localgit/otvl_rops ] ; then
  install_rops
fi