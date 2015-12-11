#!/bin/bash
set -e

echo "Installing dependencies..."
sudo apt-get update -y
sudo apt-get install -y unzip

echo "Fetching Consul..."
cd /tmp
wget https://releases.hashicorp.com/consul/0.6.0/consul_0.6.0_linux_amd64.zip -O consul.zip

echo "Installing Consul..."
unzip consul.zip >/dev/null
sudo chmod +x consul
sudo mv consul /usr/local/bin/consul
sudo mkdir -p /etc/consul.d
sudo mkdir -p /mnt/consul
sudo mkdir -p /etc/service

# Write the flags to a temporary file and move it into place
cat >/tmp/consul_flags << EOF
export CONSUL_FLAGS="-server -bootstrap-expect=3 -data-dir=/mnt/consul"
EOF
sudo mv /tmp/consul_flags /etc/service/consul
chmod 0644 /etc/service/consul

echo "Installing Upstart service..."
sudo mv /tmp/scripts/upstart.conf /etc/init/consul.conf
