#!/usr/bin/env bats

# Load the main library
. ${SCRIPTPACK_PHP_ROOT}/main.sh

@test "install PHP 7.0.x" {
  php_install_prepare
  php_install_7_0
  php --version
  [[ $(php --version) =~ 'PHP 7.0.' ]]
}
