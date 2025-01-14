#!/bin/sh

## pre
if [ `echo $0 | cut -c 1` = "/" ] ; then
  sd=`dirname $0`
else
  sd="${PWD}/`dirname $0`"
fi
. $sd/env_install.sh
## endpre

ageurl="https://github.com/FiloSottile/age/releases/download/${IV_AGE_VERSION}/age-${IV_AGE_VERSION}-linux-amd64.tar.gz"

installage() {
  is_root && \
  cmd mkdir /tmp/td.$$ && cd /tmp/td.$$ && \
  cmd curl -L ${ageurl} -o age.tgz && \
  cmd tar -xzf age.tgz && \
  cmd mv age/age age/age-keygen /usr/local/bin && \
  cd .. && rm -r /tmp/td.$$ && \
  true
  return $?
}

if [ ! -f /usr/local/bin/age ] ; then
  installage
else
  if [ "`age -version`" != "${IV_AGE_VERSION}" ] ; then
    log Age version mismatch "installed: `age -version`" "wanted: ${IV_AGE_VERSION}"
    if [ "$IV_AGE_UPDATE" ]; then
      log "Installing age from ${ageurl}"
      cmd rm /usr/local/bin/age && \
      installage
    fi
  fi
fi
