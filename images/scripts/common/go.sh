#!/bin/sh -eux

GO_VERSION="1.5.1"

# Download and install Go
wget -q -O /tmp/go.tar.gz https://storage.googleapis.com/golang/go${GO_VERSION}.linux-amd64.tar.gz
tar -C /usr/local -xzf /tmp/go.tar.gz
rm /tmp/go.tar.gz

# Setup GOPATH
mkdir -p /opt/gopath

# Setup the PATH data in the bashrc for vagrant
echo 'export PATH=/opt/gopath/bin:/usr/local/go/bin:$PATH' >> /home/vagrant/.bashrc
echo 'export GOPATH=/opt/gopath' >> /home/vagrant/.bashrc

# Install various VCS for specific platforms
if command -v apt-get >/dev/null 2>&1; then
    apt-get install -y git bzr mercurial
fi
