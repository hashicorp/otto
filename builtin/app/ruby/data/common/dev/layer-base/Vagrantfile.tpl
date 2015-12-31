{% extends "compile:data/app/dev/Vagrantfile-layer.tpl" %}

{% block vagrant_config %}
  # Setup locale
  config.vm.provision "shell", inline: $script_locale

  # Setup Ruby
  config.vm.provision "shell", inline: $script_app, privileged: false
{% endblock %}

{% block footer %}
$script_locale = <<SCRIPT
. /otto/scriptpacks/STDLIB/main.sh
. /otto/scriptpacks/RUBY/main.sh

otto_init
otto_init_locale
SCRIPT

$script_app = <<SCRIPT
# Load scriptpacks, init
. /otto/scriptpacks/STDLIB/main.sh
. /otto/scriptpacks/RUBY/main.sh
otto_init

# Configuring SSH for faster login
vagrant_config_fast_ssh

# Some params
export RUBY_VERSION="{{ ruby_version }}"

otto_output "Updating apt..."
apt_update_once

otto_output "Installing supporting packages..."
apt_install bzr git mercurial build-essential nodejs

otto_output "Installing Ruby ${RUBY_VERSION}. This can take a few minutes..."
ruby_install_prepare
ruby_install ruby-${RUBY_VERSION}

otto_output "Installing Bundler..."
oe gem install bundler --no-document

otto_output "Configuring Git to use SSH instead of HTTP so we can agent-forward private repo auth..."
oe git config --global url."git@github.com:".insteadOf "https://github.com/"
SCRIPT
{% endblock %}
