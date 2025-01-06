#!/bin/sh

if [ ! -d /usr/local/go ] ; then
  curl -L https://go.dev/dl/go1.23.4.linux-amd64.tar.gz -o go.$$.tgz && \
  tar -C /usr/local -xzf go.$$.tgz && rm go.$$.tgz && \
  true
fi
