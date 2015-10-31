{% extends "compile:data/app/dev/Vagrantfile.tpl" %}

{% block vagrant_config %}
  config.vm.provision "shell", inline: $script_app
{% endblock %}

{% block footer %}
$script_app = <<SCRIPT
set -e

# Make it so that `vagrant ssh` goes directly to the correct dir
echo "cd /vagrant" >> /home/vagrant/.profile

# Configuring SSH for faster login
if ! grep "UseDNS no" /etc/ssh/sshd_config >/dev/null; then
  echo "UseDNS no" | sudo tee -a /etc/ssh/sshd_config >/dev/null
  oe sudo service ssh restart
fi

SCRIPT
{% endblock %}
