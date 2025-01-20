#!/bin/sh

## pre
rp=`realpath $0`
sd=`dirname $rp`
. $sd/env_install.sh
## endpre

tofuurl="https://github.com/opentofu/opentofu/releases/download/v${IV_TOFU_VERSION}/tofu_${IV_TOFU_VERSION}_linux_amd64.tar.gz"

installtofu() {
  is_root && \
  cmd mkdir /tmp/td.$$ && cd /tmp/td.$$ && \
  cmd curl -L ${tofuurl} -o tofu.tgz && \
  cmd tar -xzf tofu.tgz && \
  cmd mv tofu /usr/local/bin && \
  cd .. && rm -r /tmp/td.$$ && \
  true
  return $?
}

if [ ! -f /usr/local/bin/tofu ] ; then
  installtofu
else
  if [ "`tofu version|head -1|cut -d' ' -f2`" != "v${IV_TOFU_VERSION}" ] ; then
    log tofu version mismatch "installed: `tofu version|head -1|cut -d' ' -f2`" "wanted: ${IV_TOFU_VERSION}"
    if [ "$IV_TOFU_UPDATE" ]; then
      log "Installing tofu from ${tofuurl}"
      cmd rm /usr/local/bin/tofu
      installtofu
    fi
  fi
fi
