#!/bin/sh

## pre
rp=`realpath $0`
sd=`dirname $rp`
. $sd/env_install.sh
## endpre

sopsurl="https://github.com/getsops/sops/releases/download/v${IV_SOPS_VERSION}/sops-v${IV_SOPS_VERSION}.linux.amd64"

installsops() {
  is_root && \
  cmd curl -L ${sopsurl} -o /usr/local/bin/sops && \
  cmd chmod 755 /usr/local/bin/sops && \
  true
  return $?
}

if [ ! -f /usr/local/bin/sops ] ; then
  installsops
else
  if [ "`sops -version|cut -d' ' -f2`" != "${IV_SOPS_VERSION}" ] ; then
    log sops version mismatch "installed: `sops  -version|cut -d' ' -f2`" "wanted: ${IV_SOPS_VERSION}"
    if [ "$IV_SOPS_UPDATE" ]; then
      log "Installing sops from ${sopsurl}"
      cmd rm /usr/local/bin/sops && \
      installsops
    fi
  fi
fi
