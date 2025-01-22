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
  ingress_host=`yq -r .ctr.ingress_host < $env_conf` && \
  login=`yq -r .ctr.login < $env_conf` && \
  skip_verify=`yq -r .ctr.skip_verify < $env_conf` && \
  password=`cat /root/.config/.otvl/.secrets/ctr_password.txt` && \
  install_template /etc/rancher/k3s/registries.yaml 644 $ingress_host $login $skip_verify $password && \
  log "will run curl -sfL https://get.k3s.io | sh -s - --docker" && \
  curl -sfL https://get.k3s.io | sh -s - --docker && \
  true
  return $?
}

env_conf=$CI_ENV_DIR/config.yaml
if [ ! -f "$env_conf" ] ; then
  err "$env_conf not found"
  exit 1
fi
if [ ! -f /etc/systemd/system/k3s.service ] ; then
  install_k3s
  exit 0
else
  exit 0
fi
