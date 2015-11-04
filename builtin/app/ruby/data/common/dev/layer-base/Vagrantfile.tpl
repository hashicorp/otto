{% extends "compile:data/app/dev/Vagrantfile-layer.tpl" %}

{% block vagrant_config %}
  # Setup locale
  config.vm.provision "shell", inline: $script_locale

  # Setup Ruby
  config.vm.provision "shell", inline: $script_app, privileged: false
{% endblock %}

{% block footer %}
$script_locale = <<SCRIPT
  oe() { eval "$@" 2>&1 | logger -t otto > /dev/null; }
  ol() { echo "[otto] $@"; }

  ol "Setting locale to en_US.UTF-8..."
  oe locale-gen en_US.UTF-8
  oe update-locale LANG=en_US.UTF-8 LC_ALL=en_US.UTF-8
SCRIPT

$script_app = <<SCRIPT
set -o nounset -o errexit -o pipefail -o errtrace

error() {
   local sourcefile=$1
   local lineno=$2
   echo "ERROR at ${sourcefile}:${lineno}; Last logs:"
   grep otto /var/log/syslog | tail -n 20
}
trap 'error "${BASH_SOURCE}" "${LINENO}"' ERR

# otto-exec: execute command with output logged but not displayed
oe() { eval "$@" 2>&1 | logger -t otto > /dev/null; }

# otto-log: output a prefixed message
ol() { echo "[otto] $@"; }

# Configuring SSH for faster login
if ! grep "UseDNS no" /etc/ssh/sshd_config >/dev/null; then
  echo "UseDNS no" | sudo tee -a /etc/ssh/sshd_config >/dev/null
  oe sudo service ssh restart
fi

export RUBY_INSTALL_VERSION="0.5.0"
export RUBY_VERSION="{{ ruby_version }}"

ol "Adding apt repositories and updating..."
oe sudo apt-get update

ol "Installing supporting packages..."
export DEBIAN_FRONTEND=noninteractive
oe sudo apt-get install -y bzr git mercurial build-essential nodejs

ol "Installing ruby-install v${RUBY_INSTALL_VERSION}..."
wget -O ruby-install-${RUBY_INSTALL_VERSION}.tar.gz \
  https://github.com/postmodern/ruby-install/archive/v${RUBY_INSTALL_VERSION}.tar.gz
tar -xzvf ruby-install-${RUBY_INSTALL_VERSION}.tar.gz
cd ruby-install-${RUBY_INSTALL_VERSION}/
sudo make install

ol "Installing Ruby ${RUBY_VERSION}. This can take a few minutes..."
sudo ruby-install ruby ${RUBY_VERSION} -- --disable-install-rdoc

ol "Configuring Ruby environment..."
echo 'export GEM_HOME=$HOME/.gem' >> $HOME/.ruby_env
echo "export PATH=\$HOME/.gem/bin:/opt/rubies/ruby-${RUBY_VERSION}/bin:\$PATH" >> $HOME/.ruby_env
echo 'source $HOME/.ruby_env' >> $HOME/.profile
source $HOME/.ruby_env

ol "Installing Bundler..."
oe gem install bundler --no-document

ol "Configuring Git to use SSH instead of HTTP so we can agent-forward private repo auth..."
oe git config --global url."git@github.com:".insteadOf "https://github.com/"
SCRIPT
{% endblock %}
