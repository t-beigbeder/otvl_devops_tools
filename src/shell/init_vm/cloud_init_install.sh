#!/bin/sh
gl=`curl http://169.254.169.254/openstack/latest/meta_data.json | jq -r .meta.groups | cut -d',' -f1- --output-delimiter=' '`
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
    echo "running $s"
    $s
    if [ $? -ne 0 ] ; then
      echo "$s failed, exiting"
      exit 1
    fi
  done
done