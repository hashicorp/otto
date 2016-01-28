# java_install_prepare installs the things necessary for any Java installs
java_install_prepare() {
  # Obvious: we're running non-interactive mode
  export DEBIAN_FRONTEND=noninteractive

  # Update apt once
  apt_update_once

  # Our PPAs have unicode characters, so we need to set the proper lang.
  otto_init_locale

  # Install stuff we need
  oe sudo apt-get install -y python-software-properties software-properties-common apt-transport-https
}

# java_install_8 installs Java 8
java_install_8() {
  oe sudo add-apt-repository ppa:webupd8team/java -y
  apt_update

  # Accept the Oracle license agreement
  echo debconf shared/accepted-oracle-license-v1-1 select true | oe sudo debconf-set-selections

  # Install Java
  apt_install oracle-java8-installer oracle-java8-set-default
}

# java_gradle_install installs the specified Gradle version.
java_gradle_install() {
  local version="$1"

  oe sudo add-apt-repository ppa:cwchien/gradle -y
  apt_update
  apt_install "gradle-${version}"
}

# java_maven_install installs the specified Maven version.
java_maven_install() {
  local version="$1"

  if hash mvn 2>/dev/null;
      then
      echo "Maven is already Installed"
      mvn --version | grep "Apache Maven"
  else
      maven_home="/usr/lib/apache/maven/${version}"
      #   provision the classpath
      oe sudo sh -c 'echo "export MAVEN_HOME=' + ${maven_home} + '" >> /etc/environment'
      oe sudo sh -c 'echo "PATH=\"$PATH:$MAVEN_HOME/bin\"" >> /etc/environment'
      oe sudo sh -c 'echo "export PATH" >> /etc/environment'
      #   create the directories
      oe sudo mkdir -p ${maven_home}
      #   download and extract the files
      cd ${maven_home}
      oe sudo wget http://mirrors.koehn.com/apache/maven/maven-3/${version}/binaries/apache-maven-${version}-bin.tar.gz
      oe sudo tar -zxvf apache-maven-${version}-bin.tar.gz
      #   massage the directory structure and cleanup
      oe sudo rm apache-maven-${version}-bin.tar.gz
      oe sudo mv apache-maven-${version}/* .
      oe sudo rm -rf apache-maven-${version}/
      oe source /etc/environment
  fi
}

# java_lein_install installs the specified Leiningen version.
java_lein_install() {
  local version="$1"

  if hash ./lein 2>/dev/null;
      then
      echo "Leiningen is already Installed"
      ./lein -v | grep "Leiningen"
  else
      oe sudo curl https://raw.githubusercontent.com/technomancy/leiningen/${version}/bin/lein --create-dirs -o ~/bin/lein
      oe sudo chmod a+x ~/bin/lein
  fi

}
