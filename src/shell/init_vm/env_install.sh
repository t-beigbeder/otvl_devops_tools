
setvarif() {
  if [ -z "`env | grep $1`" ]; then
    eval $1=$2
    eval export $1
  fi
}


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

install_template() {
  target=$1
  relpath=`echo $target | cut -c 2-`
  template=$IV_SD/etc/templates/$relpath
  if [ ! -f $template ] ; then
    err template $template not found
    return 1
  fi
  mod=$2
  shift 2
  cmd cp $template $target || return 1
  log instantiate $target
  c=0
  for v in "$@" ; do
    c=`expr $c + 1`
    cmd sed -i -e "s=@${c}@=${v}=" $target || return 1
  done
  cmd chmod $mod $target
}

setvarif IV_OS_VM 1
setvarif IV_GO_VERSION 1.23.4
setvarif IV_GO_UPDATE 1
setvarif IV_AGE_VERSION v1.2.1
setvarif IV_AGE_UPDATE 1
setvarif IV_SOPS_VERSION 3.9.3
setvarif IV_SOPS_UPDATE 1
setvarif IV_TOFU_VERSION 1.9.0
setvarif IV_TOFU_UPDATE 1
