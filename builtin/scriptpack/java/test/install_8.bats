#!/usr/bin/env bats

# Load the main library
. ${SCRIPTPACK_JAVA_ROOT}/main.sh

@test "install Java 8" {
  java_install_prepare
  java_install_8
  [[ $(java -version) =~ "java version \"1.8." ]]

  # Install Gradle
  java_gradle_install "2.7"
  [[ $(gradle -version) =~ "Gradle 2.7" ]]

  # Install Maven
  java_maven_install
  [[ $(mvn --version) =~ "Maven 3.0.5" ]]
}
