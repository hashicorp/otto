#!/usr/bin/env bats

# Load the main library
. ${SCRIPTPACK_PHP_ROOT}/main.sh

@test "install PHP 5.6.x" {
  php_install_prepare
  php_install_5_6
  php --version
  [[ $(php --version) =~ 'PHP 5.6.' ]]
}
