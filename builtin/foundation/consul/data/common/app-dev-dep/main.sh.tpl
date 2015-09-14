#!/bin/bash
set -e

oe() { $@ 2>&1 | logger -t otto > /dev/null; }
ol() { echo "[otto] $@"; }

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

# Reload consul. It is okay if this fails.
oe consul reload
