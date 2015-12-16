{% extends "compile:data/app/dev/Vagrantfile.tpl" %}

{% block vagrant_box %}
  config.vm.box = "bento/ubuntu-14.04"
{% endblock %}

{% block vagrant_config %}
  # Install Docker
  config.vm.provision "docker"
{% endblock %}
