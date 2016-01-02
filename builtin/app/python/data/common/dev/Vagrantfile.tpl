{% extends "compile:data/app/dev/Vagrantfile.tpl" %}

{% block vagrant_config %}
  config.vm.provision "shell", inline: $script_app
{% endblock %}

{% block footer %}
$script_app = <<SCRIPT
. /otto/scriptpacks/STDLIB/main.sh
. /otto/scriptpacks/PYTHON/main.sh
otto_init

# Make it so that the python venv is automatically sourced
echo ". /home/vagrant/virtualenv/bin/activate" >> "/home/vagrant/.bashrc"

otto_output "Setting up virtualenv in /home/vagrant/virtualenv..."
oe virtualenv --python=/usr/bin/python{{python_version}} "/home/vagrant/virtualenv"
oe chown -R vagrant:vagrant "/home/vagrant/virtualenv"

otto_output "Configuring Git to use SSH instead of HTTP so we can agent-forward private repo auth..."
oe git config --global url."git@github.com:".insteadOf "https://github.com/"
SCRIPT
{% endblock %}
