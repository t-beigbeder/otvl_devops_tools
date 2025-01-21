#!/bin/sh

## pre
rp=`realpath $0`
sd=`dirname $rp`
export IV_SD=$sd
. $sd/env_install.sh
## endpre

if [ -z "$IV_OS_VM" ] ; then
  . $sd/testdata/env_test_local.sh || exit 1
fi

echo "$CI_LIP4 ${CI_LHN}-loc" >> /etc/hosts
cmd cat /etc/hosts

export CI_ENV_DIR=/root/locgit/`basename $CI_ROPS_DIR .git`/$CI_INSTALL_ENV
if [ ! -d $CI_ENV_DIR ] ; then
  err "install environment $CI_INSTALL_ENV not found ($CI_ENV_DIR)"
  exit 1
fi

# fetching features groups from meta
if [ "$IV_OS_VM" ] ; then
  c="curl http://169.254.169.254/openstack/latest/meta_data.json"
  log running $c
  gl=`curl $c | jq -r .meta.groups | cut -d',' -f1- --output-delimiter=' '`
else
  gl=$IV_META_GROUPS
fi

# running install scripts for each feature group
if [ -z "$gl" ] ; then
  log "no group in .meta.groups, nothing to install"
  exit 0
fi
for g in $gl ; do
  if [ ! -d $g ] ; then
    log "$g: nothing to install"
    continue
  fi
  log "installing tools for group $g"
  for s in $g/* ; do
    bs=`basename $s`
    cmd $g/$bs || exit 1
  done
done