# php_install_composer installs Composer. This requires PHP is already installed.
php_install_composer() {
  curl -sS https://getcomposer.org/installer | php
  sudo mv composer.phar /usr/local/bin/composer
}
