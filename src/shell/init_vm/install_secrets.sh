#!/bin/sh

## pre
rp=`realpath $0`
sd=`dirname $rp`
. $sd/env_install.sh
## endpre

json_secrets_file=/root/clinit/bssms.json

install_from_spec() {
  st=0
  cat $IV_SD/etc/install_secrets_spec.txt | while read line ; do
      if [ -z "$line" ] ; then
        continue
      fi
      set $line
      if [ $# -ne 3 ] ; then
        echo "bad line $line ($#)"
        continue
      fi
      key=$1
      file=$2
      mod=$3
      jq -r .$key < $json_secrets_file > $file 2> /dev/null
      ct=`cat $file`
      if [ "$ct" = "null" ] ; then
        echo key $key file $file is empty
        rm $file
      else
        chmod $3 $file
      fi
  done
  return $st
}

install_secrets() {
  ld="/root/.ssh /root/.config/sops/age /root/.config/.otvl/.secrets"
  is_root && \
  mkdir -p $ld && \
  chmod go-rwx $ld && \
  install_from_spec && \
  echo rm $json_secrets_file && \
  true
  return $?
}

if [ -f $json_secrets_file ] ; then
  install_secrets
fi
