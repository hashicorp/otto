#!/usr/bin/env bats

# Load the main library
. ${SCRIPTPACK_JAVA_ROOT}/main.sh

@test "install Java 8" {
  java_install_prepare
  java_install_8
  [[ $(java -version) =~ "java version \"1.8." ]]

  # Install Gradle
  java_gradle_install "2.10"
  [[ $(gradle -version) =~ "Gradle 2.10" ]]

  # Install Maven
  java_maven_install
  [[ $(mvn --version) =~ "Maven 3.3.9" ]]

  # Install Leiningen
  java_lein_install
  [[ $(./lein -v) =~ "Leiningen 2.5.3" ]]
}
