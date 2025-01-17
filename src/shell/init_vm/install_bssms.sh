#!/bin/sh

## pre
rp=`realpath $0`
sd=`dirname $rp`
. $sd/env_install.sh
## endpre

configure_service() {
  cat > /etc/systemd/system/bssms_proxy.service <<EOF
[Unit]
Description=Runs bssms proxy

[Service]
ExecStart=/usr/local/bin/bssms px -la 0.0.0.0:9443 --host `hostname` -ut
WorkingDirectory=/root

[Install]
WantedBy=multi-user.target
EOF

}

install_bssms() {
  is_root && \
  cd ../../go/bssms && \
  /usr/local/go/bin/go build -o /usr/local/bin ./... && \
  configure_service && \
  systemctl daemon-reload && \
  systemctl enable bssms_proxy && \
  systemctl start bssms_proxy && \
  true
  return $?
}

if [ ! -f /usr/local/bin/bssms ] ; then
  install_bssms
fi