#!/bin/sh
if [ "$IV_OS_VM" ] ; then
  gl=`curl http://169.254.169.254/openstack/latest/meta_data.json | jq -r .meta.groups | cut -d',' -f1- --output-delimiter=' '`
else
  gl=$IV_META_GROUPS
fi
if [ -z "$gl" ] ; then
  echo "no group in .meta.groups, nothing to install"
  exit 0
fi
for g in $gl ; do
  if [ ! -d $g ] ; then
    echo "$g: nothing to install"
    continue
  fi
  echo "installing tools for group $g"
  for s in $g/* ; do
    bs=`basename $s`
    echo "$g: running $bs"
    $g/$bs
    if [ $? -ne 0 ] ; then
      echo "$g: $bs failed, exiting"
      exit 1
    fi
  done
done