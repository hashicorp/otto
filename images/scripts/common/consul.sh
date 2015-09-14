#!/bin/sh -eux

# Install unzip since Consul is distributed as a zip
apt-get update -y
apt-get install -y unzip

# Download Consul
cd /tmp
wget https://dl.bintray.com/mitchellh/consul/0.5.2_linux_amd64.zip -O consul.zip

# Install Consul
unzip consul.zip >/dev/null
rm consul.zip
chmod +x consul
mv consul /usr/local/bin/consul
mkdir -p /etc/consul.d
mkdir -p /mnt/consul
mkdir -p /etc/service

# Install the upstart service
cat >/etc/init/consul.conf <<EOF
description "Consul agent"

start on runlevel [2345]
stop on runlevel [!2345]

respawn

# This is to avoid Upstart re-spawning the process upon \`consul leave\`
normal exit 0 INT

# stop consul will not mark node as failed but left
kill signal INT

script
  if [ -f "/etc/service/consul" ]; then
    . /etc/service/consul
  fi

  # Make sure to use all our CPUs, because Consul can block a scheduler thread
  export GOMAXPROCS=\`nproc\`

  # Get the public IP
  BIND=\`ifconfig eth0 | grep "inet addr" | awk '{ print substr(\$2,6) }'\`

  exec /usr/local/bin/consul agent \\
    -config-dir="/etc/consul.d" \\
    -bind=\$BIND \\
    \${CONSUL_FLAGS} \\
    >>/var/log/consul.log 2>&1
end script
EOF
