{% extends "compile:data/app/dev/Vagrantfile.tpl" %}

{% block vagrant_config %}
  # Disable the default synced folder
  config.vm.synced_folder ".", "/vagrant", disabled: true

  # Read in the fragment that we use as a dep
  eval(File.read("{{ fragment_path }}"), binding)

  # Setup some stuff
  config.vm.provision "shell", inline: $script
{% endblock %}

{% block footer %}
$script = <<SCRIPT
set -e

# otto-exec: execute command with output logged but not displayed
oe() { $@ 2>&1 | logger -t otto > /dev/null; }

# otto-log: output a prefixed message
ol() { echo "[otto] $@"; }

# Configuring SSH for faster login
if ! grep "UseDNS no" /etc/ssh/sshd_config >/dev/null; then
  echo "UseDNS no" | sudo tee -a /etc/ssh/sshd_config >/dev/null
  oe sudo service ssh restart
fi

export DEBIAN_FRONTEND=noninteractive

ol "Installing HTTPS driver for Apt..."
oe sudo apt-get update
oe sudo apt-get install -y apt-transport-https
SCRIPT
{% endblock %}
