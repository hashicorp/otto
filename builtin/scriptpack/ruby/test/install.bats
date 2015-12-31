#!/usr/bin/env bats

# Load the main library
. ${SCRIPTPACK_RUBY_ROOT}/main.sh

@test "install ruby version" {
  ruby_install_prepare
  ruby_install ruby-2.2
  ruby --version
  [[ $(ruby --version) =~ 'ruby 2.2.' ]]

  # Test gem installs which should also work
  gem install bundler --no-document
  bundle --version
  [[ $(bundle --version) =~ 'Bundler version' ]]
}
