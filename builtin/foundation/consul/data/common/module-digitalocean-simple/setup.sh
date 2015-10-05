#!/bin/bash

# Write the flags to a temporary file and move it into place
cat >/tmp/consul_flags << EOF
export CONSUL_FLAGS="-server -bootstrap -data-dir=/mnt/consul"
EOF
sudo mv /tmp/consul_flags /etc/service/consul
chmod 0644 /etc/service/consul

# Restart or start consul
sudo stop consul || true
sudo start consul
