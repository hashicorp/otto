{% extends "compile:data/app/dev/Vagrantfile-layer.tpl" %}

{% block vagrant_config %}
  # Install Go build environment
  config.vm.provision "shell", inline: $script_app
{% endblock %}

{% block footer %}
$script_app = <<SCRIPT
set -e

oe() { $@ 2>&1 | logger -t otto > /dev/null; }
ol() { echo "[otto] $@"; }

ol "Updating Apt repo..."
export DEBIAN_FRONTEND=noninteractive
oe sudo apt-get update -y

ol "Downloading Node {{ node_version }}..."
oe wget -q -O /home/vagrant/node.tar.gz https://nodejs.org/dist/v{{ node_version }}/node-v{{ node_version }}-linux-x64.tar.gz

ol "Untarring Node..."
oe sudo tar -C /opt -xzf /home/vagrant/node.tar.gz

ol "Setting up PATH..."
oe sudo ln -s /opt/node-v{{ node_version }}-linux-x64/bin/node /usr/local/bin/node
oe sudo ln -s /opt/node-v{{ node_version }}-linux-x64/bin/npm /usr/local/bin/npm

ol "Installing build-essential for native packages..."
oe sudo apt-get install -y build-essential

ol "Installing GCC/G++ 4.8 (required for newer Node versions)..."
oe sudo apt-get install -y python-software-properties software-properties-common
oe sudo add-apt-repository -y ppa:ubuntu-toolchain-r/test
oe sudo apt-get update -y
oe sudo update-alternatives --remove-all gcc
oe sudo update-alternatives --remove-all g++
oe sudo apt-get install -y gcc-4.8
oe sudo apt-get install -y g++-4.8
oe sudo update-alternatives --install /usr/bin/gcc gcc /usr/bin/gcc-4.8 20
oe sudo update-alternatives --install /usr/bin/g++ g++ /usr/bin/g++-4.8 20
oe sudo update-alternatives --config gcc
oe sudo update-alternatives --config g++

ol "Installing Git..."
oe sudo apt-get install -y git

SCRIPT
{% endblock %}
