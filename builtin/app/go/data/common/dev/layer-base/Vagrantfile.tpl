{% extends "compile:data/app/dev/Vagrantfile-layer.tpl" %}

{% block vagrant_config %}
  # Install Go build environment
  config.vm.provision "shell", inline: $script_golang
{% endblock %}

{% block footer %}
$script_golang = <<SCRIPT
set -e

oe() { $@ 2>&1 | logger -t otto > /dev/null; }
ol() { echo "[otto] $@"; }

# If we have Go, then do nothing
if command -v go >/dev/null 2>&1; then
    ol "Go already installed! Otto won't install Go."
    exit 0
fi

ol "Downloading Go {{ dev_go_version }}..."
oe wget -q -O /home/vagrant/go.tar.gz https://storage.googleapis.com/golang/go{{ dev_go_version }}.linux-amd64.tar.gz

ol "Untarring Go..."
oe sudo tar -C /usr/local -xzf /home/vagrant/go.tar.gz

ol "Making GOPATH..."
oe sudo mkdir -p /opt/gopath
fstype=$(find /opt/gopath -mindepth 0 -maxdepth 0 -type d -printf "%F")
find /opt/gopath -fstype ${fstype} -print0 | xargs -0 -n 100 chown vagrant:vagrant

ol "Setting up PATH..."
echo 'export PATH=/opt/gopath/bin:/usr/local/go/bin:$PATH' >> /home/vagrant/.bashrc
echo 'export GOPATH=/opt/gopath' >> /home/vagrant/.bashrc

ol "Installing VCSs for go get..."
oe sudo apt-get update -y
oe sudo apt-get install -y git bzr mercurial

ol "Configuring Go to use SSH instead of HTTP..."
git config --global url."git@github.com:".insteadOf "https://github.com/"
SCRIPT
{% endblock %}
