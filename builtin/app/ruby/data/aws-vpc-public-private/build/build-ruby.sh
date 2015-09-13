#!/bin/bash

# otto-exec: execute command with output logged but not displayed
oe() { $@ 2>&1 | logger -t otto > /dev/null; }

# otto-log: output a prefixed message
ol() { echo "[otto] $@"; }

# Wait for cloud-config to complete, which can interfere with apt commands if
# it's still running
until [[ -f /var/lib/cloud/instance/boot-finished ]]; do
  sleep 1
done

# Prevent any apt operations from ever asking for input
export DEBIAN_FRONTEND=noninteractive

oe sudo apt-get update

# TODO: parameterize ruby version somehow
export RUBY_VERSION="2.2"
ol "Installing Ruby ${RUBY_VERSION} & Passenger..."
oe sudo apt-get install -y python-software-properties
oe sudo apt-add-repository -y ppa:brightbox/ruby-ng
oe sudo apt-get update
oe sudo apt-get install -y ruby$RUBY_VERSION ruby$RUBY_VERSION-dev apache2

ol "Installing Bundler..."
oe gem install bundler --no-ri --no-rdoc

ol "Installing VCSs for bundle install..."
oe sudo apt-get install -y git bzr mercurial

ol "Installing build-essential for native gem builds..."
oe sudo apt-get install -y build-essential
