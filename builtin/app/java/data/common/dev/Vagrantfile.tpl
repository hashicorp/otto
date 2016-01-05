{% extends "compile:data/app/dev/Vagrantfile.tpl" %}

{% block vagrant_config %}
  config.vm.provision "shell", inline: $script_app
{% endblock %}

{% block footer %}
$script_app = <<SCRIPT
. /otto/scriptpacks/STDLIB/main.sh
. /otto/scriptpacks/JAVA/main.sh
otto_init

# Make it so that `vagrant ssh` goes directly to the correct dir
vagrant_default_cd "vagrant" "/vagrant"

# Install Java
otto_output "Preparing to install Java..."
java_install_prepare
otto_output "Installing Java 8..."
java_install_8
otto_output "Installing Gradle {{ gradle_version }}..."
java_gradle_install "{{ gradle_version }}"
otto_output "Installing Maven..."
java_maven_install
SCRIPT
{% endblock %}
