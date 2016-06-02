# php_install_prepare should be called before any other installation is
# done. This configures the system to be ready to install a PHP version.
php_install_prepare() {
  # Obvious: we're running non-interactive mode
  export DEBIAN_FRONTEND=noninteractive

  # Update apt once
  apt_update_once

  # Our PPAs have unicode characters, so we need to set the proper lang.
  otto_init_locale

  sudo apt-get install -y python-software-properties software-properties-common apt-transport-https

  add-apt-repository -y ppa:ondrej/php
}

# php_install installs an arbitrary PHP version given in the argument
php_install() {
    local version="$1"
    case $version in
    7.0)
        php_install_prepare
        php_install_7_0
        ;;
    5.6)
        php_install_prepare
        php_install_5_6
        ;;
    5.5)
        php_install_prepare
        php_install_5_5
        ;;
    *)
        echo "Unknown PHP version: ${version}"
        exit 1
        ;;
    esac
}

# php_install_5_5 installs PHP 5.5.x
php_install_5_5() {
  apt-get update
  sudo apt-get install -y php5.5
}

# php_install_5_6 installs PHP 5.6.x
php_install_5_6() {
  apt-get update
  sudo apt-get install -y php5.6
}

# php_install_7_0 installs PHP 7.0.x
php_install_7_0() {
  apt-get update
  sudo apt-get install -y php7.0
}
