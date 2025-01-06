#!/bin/sh

if [ ! -d /root/dot_init ] ; then
  apt-get update && apt-get upgrade && \
  apt-get install -y --no-install-recommends curl git jq && \
  git config --global credential.helper store && \
  mkdir /root/dot_init && \
  cd /root/dot_init && \
  git clone --single-branch --branch bdev9 https://github.com/t-beigbeder/otvl_devops_tools && \
  true || exit 1
fi
