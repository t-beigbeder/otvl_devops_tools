#!/bin/sh

## pre
rp=`realpath $0`
sd=`dirname $rp`
. $sd/env_install.sh
## endpre

install_rops() {
  cd /root/locgit && \
  log GIT_SSH_COMMAND="ssh -o StrictHostKeyChecking=accept-new -i /root/.ssh/id_rsa_gh_ssh" git clone $CI_ROPS_REPO && \
  GIT_SSH_COMMAND="ssh -o StrictHostKeyChecking=accept-new -i /root/.ssh/id_rsa_gh_ssh" git clone $CI_ROPS_REPO && \
  true
  return $?
}

if [ ! -d /root/locgit/otvl_rops ] ; then
  install_rops || exit 1
fi
if [ ! -d $CI_ENV_DIR ] ; then
  err "install environment $CI_INSTALL_ENV not found ($CI_ENV_DIR)"
  exit 1
fi

install_template /usr/local/bin/env_rops.sh 755 "${CI_LHN}" "${CI_LIP4}" "${CI_DOT_REPO}" "${CI_DOT_BRANCH}" "${CI_ROPS_REPO}" "${CI_INSTALL_ENV}"
