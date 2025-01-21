#!/bin/sh

## pre
rp=`realpath $0`
sd=`dirname $rp`
. $sd/env_install.sh
## endpre

install_k3s() {
  mkdir -p /etc/rancher/k3s /var/lib/rancher/k3s/server/manifests && \
  chmod 755 /etc/rancher/k3s /var/lib/rancher/k3s/server/manifests && \
  install_template /etc/rancher/k3s/rancher-config.yaml 644 "${CI_LHN}-loc" && \
  install_template /var/lib/rancher/k3s/server/manifests/traefik-config.yaml 644 && \
  install_template /etc/rancher/k3s/registries.yaml 644 && \
  log "will run curl -sfL https://get.k3s.io | sh -s - --docker" && \
  curl -sfL https://get.k3s.io | sh -s - --docker && \
  true
  return $?
}

if [ ! -f /etc/systemd/system/k3s.service ] ; then
  install_k3s
  exit 0
else
  exit 0
fi
