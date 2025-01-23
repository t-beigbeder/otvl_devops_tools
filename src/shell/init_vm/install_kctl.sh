#!/bin/sh

## pre
rp=`realpath $0`
sd=`dirname $rp`
. $sd/env_install.sh
## endpre

kctl_url="https://dl.k8s.io/release/${IV_KCTL_VERSION}/bin/linux/amd64/kubectl"

install_kctl() {
  is_root && \
  cd /tmp && mkdir d$$ && cd d$$ && \
  cmd curl -LO "$kctl_url" && \
  cmd file kubectl && \
  file kubectl | grep "ELF 64-bit LSB executable" > /dev/null && \
  cmd chmod 755 kubectl && \
  mv kubectl /usr/local/bin && \
  cd .. && rm -r d$$ && \
  true
  return $?
}

if [ ! -f /usr/local/bin/kubectl ] ; then
  install_kctl
else
  cv=`kubectl version --client | head -1 | cut -d' ' -f3-`
  if [ "$cv" != "${IV_KCTL_VERSION}" ] ; then
    log kubectl version mismatch "installed: $cv" "wanted: ${IV_KCTL_VERSION}"
    if [ "$IV_KCTL_UPDATE" ]; then
      log "Installing kubectl from ${kctl_url}"
      cmd rm /usr/local/bin/kubectl && \
      install_kctl
    fi
  fi
fi
