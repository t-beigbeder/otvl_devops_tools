
disp() {
    echo >&2 "$@"
}

log() {
    disp "`date -Iseconds`" "$@"
}

err() {
    disp "`date -Iseconds`" "ERROR:" "$@"
}

cmd() {
    log running "$@"
    "$@"
    st=$?
    if [ $st -ne 0 ] ; then
      err "running" "$@" "failed"
    fi
    return $st
}

cmd_if() {
    echo -n "Run" "$@" "? "
    read ans
    if [ "$ans" = "y" ] ; then
        cmd "$@"
    fi
}

is_root() {
  if [ `id -u` != 0 ] ; then
    disp "Must be root"
    return 1
  fi
  return 0
}

export CI_LHN=@1@
export CI_LIP4=@2@
export CI_DOT_REPO=@3@
export CI_DOT_BRANCH=@4@
export CI_ROPS_REPO=@5@
export CI_INSTALL_ENV=@6@

export CI_ENV_DIR=/root/locgit/`basename $CI_ROPS_REPO .git`/$CI_INSTALL_ENV
