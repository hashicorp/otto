{% extends "compile:data/app/dev/Vagrantfile-layer.tpl" %}

{% block vagrant_config %}
  # Setup Python
  config.vm.provision "shell", inline: $script_app
{% endblock %}

{% block footer %}
$script_app = <<SCRIPT
#!/bin/bash
set -e

# Setup our scriptpacks
. /otto/scriptpacks/STDLIB/main.sh
. /otto/scriptpacks/PYTHON/main.sh

# Initialize
otto_init

# Make it so that `vagrant ssh` goes directly to the correct dir
vagrant_default_cd "vagrant" "/vagrant"

# Configuring SSH for faster login
vagrant_config_fast_ssh

# Install Python
otto_output "Installing Python Version {{python_version}} "
oe python_install "{{python_version}}"

otto_output "Installing supporting packages..."
oe sudo apt-get install -y \
  bzr git mercurial build-essential \
  libpq-dev zlib1g-dev software-properties-common \
  libsqlite3-dev \
  curl
SCRIPT
{% endblock %}
