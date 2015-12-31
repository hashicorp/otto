#!/usr/bin/env bats

# Load the main library
. ${SCRIPTPACK_RUBY_ROOT}/main.sh

@test "install ruby-install" {
  ruby_install_rubyinstall
  ruby-install
}
