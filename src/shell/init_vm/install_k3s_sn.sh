#!/bin/sh

## pre
rp=`realpath $0`
sd=`dirname $rp`
. $sd/env_install.sh
## endpre

if [ -f /etc/systemd/system/k3s.service ] ; then
  exit 0
fi
mkdir -p /etc/rancher/k3s && \
  chmod 755 /etc/rancher/k3s && \
  install_template /etc/rancher/k3s/rancher-config.yaml 644 "${CI_LHN}-loc"
  true
# curl -sfL https://get.k3s.io | sh -s - --docker
