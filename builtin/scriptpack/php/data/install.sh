# php_install_prepare should be called before any other installation is
# done. This configures the system to be ready to install a PHP version.
php_install_prepare() {
  # Obvious: we're running non-interactive mode
  export DEBIAN_FRONTEND=noninteractive

  # Update apt once
  #apt_update_once

  # Our PPAs have unicode characters, so we need to set the proper lang.
  if [[ ! $(locale -a) =~ '^en_US\.utf8' ]]; then
      sudo locale-gen en_US.UTF-8
  fi
  export LANG=en_US.UTF-8

  sudo apt-get install -y python-software-properties software-properties-common apt-transport-https
}

# php_install_5_5 installs PHP 5.5.x
php_install_5_5() {
  add-apt-repository -y ppa:ondrej/php5
  apt-get update
  sudo apt-get install -y php5
}

# php_install_5_6 installs PHP 5.6.x
php_install_5_6() {
  add-apt-repository -y ppa:ondrej/php5-5.6
  apt-get update
  sudo apt-get install -y php5
}

# php_install_7_0 installs PHP 7.0.x
php_install_7_0() {
  add-apt-repository -y ppa:ondrej/php-7.0
  apt-get update
  sudo apt-get install -y php7.0
}
