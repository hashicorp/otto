#!/usr/bin/env bats

# Load the main library
. ${SCRIPTPACK_PHP_ROOT}/main.sh

@test "install PHP composer" {
  php_install "5.5"
  php_install_composer
  [[ $(composer -V) =~ 'Composer version' ]]
}
