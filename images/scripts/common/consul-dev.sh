#!/bin/sh -eux

# Make the flags a single-server setup
cat >/etc/service/consul << EOF
export CONSUL_FLAGS="-server -bootstrap -data-dir=/mnt/consul"
EOF
