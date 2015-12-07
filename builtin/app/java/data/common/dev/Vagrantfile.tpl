{% extends "compile:data/app/dev/Vagrantfile.tpl" %}

{% block vagrant_config %}
  config.vm.provision "shell", inline: $script_app
{% endblock %}

{% block footer %}
$script_app = <<SCRIPT
. /otto/scriptpacks/STDLIB/main.sh
. /otto/scriptpacks/JAVA/main.sh
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

ol "Downloading Java 8..."
oe sudo aptitude install software-properties-common python-software-properties -y
oe sudo aptitude update -y
oe sudo add-apt-repository ppa:webupd8team/java -y
oe sudo aptitude update -y
oe sudo aptitude install -y --force-yes oracle-java8-installer oracle-java8-set-default

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
sudo wget www.scala-lang.org/files/archive/scala-2.10.4.deb
oe sudo dpkg -i scala-2.10.4.deb
oe sudo aptitude update
oe sudo aptitude install scala

ol "Downloading SBT..."
wget http://scalasbt.artifactoryonline.com/scalasbt/sbt-native-packages/org/scala-sbt/sbt/0.12.4/sbt.deb
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
