#!/bin/sh
. /usr/local/bin/env_rops.sh
cld=$HOME/locgit
grd="$cld/`basename $CI_ROPS_REPO .git`"
if [ ! -d $grd ] ; then
  cmd mkdir -p $cld && \
  cd $cld && \
  log GIT_SSH_COMMAND="ssh -o StrictHostKeyChecking=accept-new -i /root/.ssh/id_rsa_gh_ssh" git clone $CI_ROPS_REPO && \
  GIT_SSH_COMMAND="ssh -o StrictHostKeyChecking=accept-new -i /root/.ssh/id_rsa_gh_ssh" git clone $CI_ROPS_REPO && \
  true
fi
