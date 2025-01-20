#!/bin/sh

## pre
rp=`realpath $0`
sd=`dirname $rp`
. $sd/env_install.sh
## endpre

if [ ! -d /root/localgit/otvl_rops ] ; then
  cd /root/localgit && \
    GIT_SSH_COMMAND="ssh -o StrictHostKeyChecking=accept-new -i /root/.ssh/id_rsa_gh_ssh" git clone git@github.com:t-beigbeder/otvl_rops.git && \
    true
fi