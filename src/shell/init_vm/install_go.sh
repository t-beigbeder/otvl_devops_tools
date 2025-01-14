#!/bin/sh

## pre
if [ `echo $0 | cut -c 1` = "/" ] ; then
  sd=`dirname $0`
else
  sd="${PWD}/`dirname $0`"
fi
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
  if [ "`go version | cut -f3 -d' '`" != "go${IV_GO_VERSION}" ] ; then
    log Go version mismatch "installed: `go version | cut -f3 -d' '`" "wanted: go${IV_GO_VERSION}"
    if [ "$IV_GO_UPDATE" ]; then
      log "Installing go from ${gourl}"
      cmd rm -r /usr/local/go && \
      installgo
    fi
  fi
fi
