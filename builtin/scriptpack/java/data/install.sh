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

# java_maven_install installs Maven. Eventually this will take an argument.
java_maven_install() {
  apt_install maven
}
