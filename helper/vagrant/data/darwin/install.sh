#!/bin/bash
set -e

oe() { $@ 2>&1 | logger -t otto > /dev/null; }
ol() { echo "[otto] $@"; }

ol "Downloading Vagrant (this can take a few minutes)..."
curl -L -# $1 >vagrant.dmg

ol "Attaching Vagrant disk image..."
oe hdiutil attach vagrant.dmg

ol "Starting Vagrant installer..."
sudo installer -pkg /Volumes/Vagrant/Vagrant.pkg -target "/"

ol "Vagrant installed. Cleaning up..."
oe hdiutil detach /Volumes/Vagrant/
oe rm vagrant.dmg
