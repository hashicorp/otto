#!/bin/bash
set -e

oe() { $@ 2>&1 | logger -t otto > /dev/null; }
ol() { echo "[otto] $@"; }

dir=$(pwd)

# Source the variables
source ./vars

# Download and setup Consul directories
if ! command -v consul >/dev/null 2>&1; then
    ol "Installing Consul..."
    oe sudo apt-get update -y
    oe sudo apt-get install -y unzip
    cd /tmp
    oe wget https://dl.bintray.com/mitchellh/consul/0.5.2_linux_amd64.zip -O consul.zip
    oe unzip consul.zip
    oe sudo chmod +x consul
    oe sudo mv consul /usr/local/bin/consul
    oe sudo mkdir -p /etc/consul.d
    oe sudo mkdir -p /mnt/consul
    oe sudo mkdir -p /etc/service

    # Write the flags to a temporary file and move it into place
    cat >/tmp/consul_flags << EOF
export CONSUL_FLAGS="-data-dir=/mnt/consul -join=${consul_address}"
EOF
    oe chmod 0644 /tmp/consul_flags
    oe sudo mv /tmp/consul_flags /etc/service/consul

    # Setup Consul service and start it
    oe sudo cp ${dir}/upstart.conf /etc/init/consul.conf

    # Setup DNS
    ol "Installing dnsmasq for Consul..."
    oe sudo apt-get install -y dnsmasq
    echo "server=/consul/127.0.0.1#8600" > /tmp/dnsmasq
    oe sudo mv /tmp/dnsmasq /etc/dnsmasq.d/10-consul
    oe sudo /etc/init.d/dnsmasq restart
fi

# Write the service file
ol "Configuring consul service: {{ app_config.ServiceName }}"
cat <<DOC >/tmp/service.json
{
  "service": {
    "name": "{{ app_config.ServiceName }}",
    "tags": [],
    "port": {{ app_config.ServicePort }}
  }
}
DOC
oe chmod 0644 /tmp/service.json
oe sudo mv /tmp/service.json /etc/consul.d/service.{{ app_config.ServiceName }}.json
