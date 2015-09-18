#!/bin/bash
set -e

# Install cURL if we have to
apt-get update -y
apt-get install -y curl

# Install Docker
curl -sSL https://get.docker.com/ | sh

# Create the container
docker create {{ run_args }} --name="{{ name }}" {{ docker_image }}

# Write the service
cat >/etc/init/{{ name }}.conf <<EOF
description "Docker container: {{ name }}"

start on filesystem and started docker
stop on runlevel [!2345]

respawn

post-stop exec sleep 5

script
  /usr/bin/docker start {{ name }}
end script
EOF

# Start the service
start {{ name }}
