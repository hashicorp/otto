{% extends "../app/dev/Vagrantfile-layer.tpl" %}

{% block vagrant_config %}
  dir = "/otto/foundation-layer-{{ foundation_id }}"
  config.vm.synced_folder '{{ foundation_dir }}', dir
  config.vm.provision "shell", inline: "cd #{dir} && bash #{dir}/layer.sh"
{% endblock %}
