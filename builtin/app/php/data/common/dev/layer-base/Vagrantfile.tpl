{% extends "compile:data/app/dev/Vagrantfile-layer.tpl" %}

{% block vagrant_config %}
  # Setup Ruby
  config.vm.provision "shell", inline: $script_app
{% endblock %}

{% block footer %}
$script_app = <<SCRIPT
#!/bin/bash
set -e

# Setup our scriptpacks
. /otto/scriptpacks/STDLIB/main.sh
. /otto/scriptpacks/PHP/main.sh

# Initialize
otto_init

# Make it so that `vagrant ssh` goes directly to the correct dir
vagrant_default_cd "vagrant" "/vagrant"

# Configuring SSH for faster login
vagrant_config_fast_ssh

# Install PHP
otto_output "Installing PHP Version {{php_version}} "
oe php_install "{{php_version}}"

otto_output "Installing supporting packages..."
oe sudo apt-get install -y \
  bzr git mercurial build-essential \
  curl

otto_output "Installing Composer..."
cd /tmp
curl -sS https://getcomposer.org/installer | php
oe sudo mv composer.phar /usr/local/bin/composer
SCRIPT
{% endblock %}
