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
    cat $sd/install_secrets_spec.txt | while read line ; do
        echo line $line
    done
}

install_secrets() {
  is_root && \
  install_from_spec && \
  echo rm $json_secrets_file && \
  true
  return $?
}

if [ -f $json_secrets_file ] ; then
  install_secrets
fi
