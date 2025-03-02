#!/bin/sh

## pre
rp=`realpath $0`
sd=`dirname $rp`
. $sd/env_install.sh
## endpre

mcurl="https://github.com/FiloSottile/mkcert/releases/download/${IV_MKCERT_VERSION}/mkcert-${IV_MKCERT_VERSION}-linux-amd64"

installmkcert() {
  is_root && \
  cmd apt-get install -y --no-install-recommends libnss3-tools && \
  cmd curl -L ${mcurl} -o /usr/local/bin/mkcert && \
  cmd chmod +755 /usr/local/bin/mkcert && \
  true
  return $?
}

if [ ! -f /usr/local/bin/mkcert ] ; then
  log "Installing mkcert from ${mcurl}"
  installmkcert
else
  mv=`mkcert -version`
  if [ "$mv" != "${IV_MKCERT_VERSION}" ] ; then
    log mkcert version mismatch "installed $mv" "wanted: ${IV_MKCERT_VERSION}"
    if [ "$IV_MKCERT_UPDATE" ] ; then
      log "Installing mkcert from ${mcurl}"
      cmd rm /usr/local/bin/mkcert && \
      installmkcert
    fi
  fi
fi