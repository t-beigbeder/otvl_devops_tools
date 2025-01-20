#!/bin/sh

## pre
if [ `echo $0 | cut -c 1` = "/" ] ; then
  sd=`dirname $0`
else
  sd="${PWD}/`dirname $0`"
fi
. $sd/env_install.sh
## endpre

json_secrets_file=/root/clinit/bssms.json

install_from_spec() {
  st=0
  cat $sd/install_secrets_spec.txt | while read line ; do
      set `echo $line cut -d' ' -f1-3`
      key=$1
      file=$2
      mod=$3
      jq -r .$key < $json_secrets_file > $file 2> /dev/null
      if [ -s $file ] ; then
        rm $file
        st=1
      fi
      chmod $3 $file
  done
  return $st
}

install_secrets() {
  is_root && \
  mkdir -p /root/.ssh /root/.config/sops/age && \
  chmod go-rwx /root/.ssh /root/.config/sops/age && \
  install_from_spec && \
  echo rm $json_secrets_file && \
  true
  return $?
}

if [ -f $json_secrets_file ] ; then
  install_secrets
fi
