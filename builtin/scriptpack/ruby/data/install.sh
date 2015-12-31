# ruby_install_prepare installs ruby-install and chruby
ruby_install_prepare() {
    rubyinstall_version="0.6.0"
    chruby_version="0.3.9"

    # If we don't have wget, install that
    if ! hash wget 2>/dev/null; then
        apt_update_once
        oe sudo apt-get install -y wget
    fi

    # Download/Install ruby-install
    pushd /tmp >/dev/null
    oe wget -O ruby-install-${rubyinstall_version}.tar.gz \
        https://github.com/postmodern/ruby-install/archive/v${rubyinstall_version}.tar.gz
    oe tar -xzvf ruby-install-${rubyinstall_version}.tar.gz
    cd ruby-install-${rubyinstall_version}/
    oe sudo make install
    popd >/dev/null

    # Download/Install chruby
    pushd /tmp >/dev/null
    oe wget -O chruby-${chruby_version}.tar.gz \
        https://github.com/postmodern/chruby/archive/v${chruby_version}.tar.gz
    oe tar -xzvf chruby-${chruby_version}.tar.gz
    cd chruby-${chruby_version}/
    oe sudo make install
    popd >/dev/null

    # Make /opt/rubies ahead of time so that chruby configures properly
    sudo mkdir -p /opt/rubies

    # Install chruby system-wide
    cat <<EOF >/tmp/chruby.sh
if [ -n "\$BASH_VERSION" ] || [ -n "\$ZSH_VERSION" ]; then
  source /usr/local/share/chruby/chruby.sh
fi
EOF
    sudo mv /tmp/chruby.sh /etc/profile.d/chruby.sh
    sudo chmod +x /etc/profile.d/chruby.sh
    source /etc/profile.d/chruby.sh
}

# ruby_install installs the given ruby version. The parameter should be
# the ruby version to install, such as "ruby-2.1.3"
ruby_install() {
    local version=$1

    # Install the ruby version
    oe sudo ruby-install $version -- --disable-install-rdoc

    # Resource chruby so it detects our new ruby
    . /etc/profile.d/chruby.sh

    # Configure the environment to use it
    echo "chruby ${version}" >> $HOME/.profile
    chruby ${version}
}
