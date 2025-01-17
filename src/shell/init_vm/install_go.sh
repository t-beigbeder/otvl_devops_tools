#!/bin/sh

## pre
rp=`realpath $0`
sd=`dirname $rp`
. $sd/env_install.sh
## endpre

gourl="https://go.dev/dl/go${IV_GO_VERSION}.linux-amd64.tar.gz"

installgo() {
  is_root && \
  cmd curl -L $gourl -o go.$$.tgz && \
  cmd tar -C /usr/local -xzf go.$$.tgz && rm go.$$.tgz && \
  true
  return $?
}

if [ ! -d /usr/local/go ] ; then
  installgo
else
  PATH=/usr/local/go/bin:$PATH
  if [ "`go version | cut -f3 -d' '`" != "go${IV_GO_VERSION}" ] ; then
    log Go version mismatch "installed: `go version | cut -f3 -d' '`" "wanted: go${IV_GO_VERSION}"
    if [ "$IV_GO_UPDATE" ]; then
      log "Installing go from ${gourl}"
      is_root && \
      cmd rm -r /usr/local/go && \
      installgo
    fi
  fi
fi
