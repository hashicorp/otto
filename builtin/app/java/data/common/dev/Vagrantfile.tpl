{% extends "compile:data/app/dev/Vagrantfile.tpl" %}

{% block vagrant_config %}
  config.vm.provision "shell", inline: $script_app
{% endblock %}

end

$script_app = <<SCRIPT
#!/bin/bash
set -e

# Setup our scriptpacks
. /otto/scriptpacks/STDLIB/main.sh
. /otto/scriptpacks/JAVA/main.sh

# Initialize
otto_init

oe() { $@ 2>&1 | logger -t otto > /dev/null; }
ol() { echo "[otto] $@"; }

# Make it so that `vagrant ssh` goes directly to the correct dir
vagrant_default_cd "vagrant" "/vagrant"

# Configuring SSH for faster login
if ! grep "UseDNS no" /etc/ssh/sshd_config >/dev/null; then
  echo "UseDNS no" | sudo tee -a /etc/ssh/sshd_config >/dev/null
  oe sudo service ssh restart
fi

export DEBIAN_FRONTEND=noninteractive
ol "Upgrading Outdated Apt Packages..."
oe sudo aptitude update -y
oe sudo aptitude upgrade -y
oe sudo do-release-upgrade -f DistUpgradeViewNonInteractive

ol "Installing requirements to add ppa repositories for Java and Gradle installs."
oe sudo aptitude install software-properties-common python-software-properties -y

ol "Downloading Java 8..."
echo oracle-java8-installer shared/accepted-oracle-license-v1-1 select true | sudo /usr/bin/debconf-set-selections
oe sudo add-apt-repository ppa:webupd8team/java -y
oe sudo aptitude update -y
oe sudo aptitude install oracle-java8-installer -y
ol "Setting environment variables for Java 8.."
oe sudo aptitude install oracle-java8-set-default -y

ol "Downloading Gradle {{ gradle_version }}..."
oe sudo add-apt-repository ppa:cwchien/gradle -y
oe sudo aptitude update -y
oe sudo apt-cache search gradle
oe sudo aptitude install gradle-{{ gradle_version }} -y

ol "Downloading Maven..."
oe sudo aptitude update -y
oe sudo aptitude install maven -y

ol "Downloading Scala..."
oe sudo aptitude remove scala-library scala
oe sudo wget www.scala-lang.org/files/archive/scala-{{ scala_version }}.deb
oe sudo dpkg -i scala-{{ scala_version }}.deb
oe sudo aptitude update
oe sudo aptitude install scala

ol "Downloading SBT..."
oe wget https://bintray.com/artifact/download/sbt/debian/sbt-{{ sbt_version }}.deb
oe sudo dpkg -i sbt.deb
oe sudo aptitude update
oe sudo aptitude install sbt

ol "Installing Git..."
oe sudo add-apt-repository ppa:git-core/ppa -y
oe sudo aptitude update -y
oe sudo aptitude install git -y
git config --global url."git@github.com:".insteadOf "https://github.com/"

SCRIPT
{% endblock %}
