{% extends "compile:data/app/dev/Vagrantfile.tpl" %}

{% block vagrant_box %}
  config.vm.box = "bento/ubuntu-14.04"
{% endblock %}

{% block vagrant_config %}
  # Make `vagrant ssh` go directly to the right place
  config.vm.provision "shell", inline: %Q[echo "cd /vagrant" >> /home/vagrant/.bashrc]

  # Install Docker
  config.vm.provision "docker"
{% endblock %}
