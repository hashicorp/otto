{% extends "compile:data/app/dev/Vagrantfile.tpl" %}

{% block default_shared_folder %}
  # Setup a synced folder from our working directory to /vagrant
  config.vm.synced_folder '{{ path.working }}', "{{ shared_folder_path }}",
    owner: "vagrant", group: "vagrant"
{% endblock %}

{% block vagrant_config %}
  {% if import_path != "" %}
  # Disable the default synced folder
  config.vm.synced_folder ".", "/vagrant", disabled: true
  {% endif %}

  # Install Go build environment
  config.vm.provision "shell", inline: $script_golang

  # Make it so that `vagrant ssh` goes directly to the correct dir
  config.vm.provision "shell", inline:
    %Q[echo "cd {{ shared_folder_path }}" >> /home/vagrant/.bashrc]
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
