#!/bin/sh
. /usr/local/bin/env_rops.sh
cld=$HOME/locgit
grd="$cld/`basename $CI_DOT_REPO .git`"
if [ ! -d $grd ] ; then
  cmd mkdir -p $cld && \
  cd $cld && \
  cmd git clone --single-branch --branch $CI_DOT_BRANCH $CI_DOT_REPO
fi
