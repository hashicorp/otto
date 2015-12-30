{% extends "compile:data/app/dev/Vagrantfile-layer.tpl" %}

{% block vagrant_config %}
  # Setup Ruby
  config.vm.provision "shell", inline: $script_app
{% endblock %}

{% block footer %}
$script_app = <<SCRIPT
#!/bin/bash

# Setup our scriptpacks
. /otto/scriptpacks/STDLIB/main.sh
. /otto/scriptpacks/PHP/main.sh

# Make it so that `vagrant ssh` goes directly to the correct dir
echo "cd /vagrant" >> /home/vagrant/.bashrc

# Configuring SSH for faster login
if ! grep "UseDNS no" /etc/ssh/sshd_config >/dev/null; then
  echo "UseDNS no" | sudo tee -a /etc/ssh/sshd_config >/dev/null
  oe sudo service ssh restart
fi

# Install PHP
otto_output "Installing PHP Version {{php_version}} "
php_install "{{php_version}}"

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
