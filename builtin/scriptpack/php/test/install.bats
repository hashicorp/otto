#!/usr/bin/env bats

# Load the main library
. ${SCRIPTPACK_PHP_ROOT}/main.sh

@test "install from string '5.5'" {
  php_install "5.5"
  php --version
  [[ $(php --version) =~ 'PHP 5.5.' ]]
}

@test "install from string '5.6'" {
  php_install "5.6"
  php --version
  [[ $(php --version) =~ 'PHP 5.6.' ]]
}

@test "install from string '7.0'" {
  php_install "7.0"
  php --version
  [[ $(php --version) =~ 'PHP 7.0.' ]]
}
