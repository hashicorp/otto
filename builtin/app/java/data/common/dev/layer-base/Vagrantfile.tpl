{% extends "compile:data/app/dev/Vagrantfile-layer.tpl" %}

{% block vagrant_config %}
  config.vm.provision "shell", inline: $script_app
{% endblock %}

{% block footer %}
$script_app = <<SCRIPT
#!/bin/bash
set -e

# Setup our scriptpacks
. /otto/scriptpacks/STDLIB/main.sh
. /otto/scriptpacks/JAVA/main.sh

# Initialize
otto_init

# Make it so that `vagrant ssh` goes directly to the correct dir
vagrant_default_cd "vagrant" "/vagrant"

# Install Java
otto_output "Preparing to install Java..."
oe java_install_prepare
otto_output "Installing Java 8..."
oe java_install_8
otto_output "Installing Gradle {{ gradle_version }}..."
oe java_gradle_install "{{ gradle_version }}"
otto_output "Installing Maven {{ maven_version }}..."
oe java_maven_install "{{ maven_version }}"
otto_output "Installing Scala {{ scala_version }}..."
oe java_scala_install "{{ scala_version }}"
otto_output "Installing sbt {{ sbt_version }}..."
oe java_sbt_install "{{ sbt_version }}"
otto_output "Installing Leiningen {{ lein_version }}..."
oe java_lein_install "{{ lein_version }}"

otto_output "Installing Git..."
oe sudo apt-get install -y git

SCRIPT
{% endblock %}
