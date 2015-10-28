{% extends "compile:data/app/dev/Vagrantfile.tpl" %}

{% block vagrant_config %}
  config.vm.provision "shell", inline: $script_app, privileged: false
{% endblock %}

{% block footer %}
$script_app = <<SCRIPT
set -e

# otto-exec: execute command with output logged but not displayed
oe() { eval "$@" 2>&1 | logger -t otto > /dev/null; }

# otto-log: output a prefixed message
ol() { echo "[otto] $@"; }

# Make it so that `vagrant ssh` goes directly to the correct dir
echo "cd /vagrant" >> $HOME/.profile

source $HOME/.ruby_env

has_gem() {
  gem_name=$1

  if [ -f Gemfile.lock ]; then
    grep -e " $gem_name \(" Gemfile.lock > /dev/null
    return $?
  fi

  if [ -f Gemfile ]; then
    grep -e "gem .$gem_name." Gemfile > /dev/null
    return $?
  fi

  return 1
}

gem_deps_queue=()

detect_gem_deps() {
  gem_name=$1; apt_deps=$2

  if has_gem $gem_name; then
    ol "Detected the $gem_name gem"
    gem_deps_queue+=($apt_deps)
  fi
}

cd /vagrant
detect_gem_deps curb "libcurl3 libcurl3-gnutls libcurl4-openssl-dev"
detect_gem_deps capybara-webkit "libqt4-dev"
detect_gem_deps mysql2 "libmysqlclient-dev"
detect_gem_deps nokogiri "zlib1g-dev"
detect_gem_deps pg "libpq-dev"
detect_gem_deps rmagick "libmagickwand-dev"
detect_gem_deps sqlite3 "libsqlite3-dev"
detect_gem_deps libxml-ruby "libxml-dev"
detect_gem_deps paperclip "imagemagick"

if [ -n "${gem_deps_queue-}" ]; then
  ol "Installing native gem system dependencies..."
  oe sudo apt-get install -y "${gem_deps_queue[@]}"
fi

ol "Bundling gem dependencies..."
oe bundle

{% if app_type == "rails" %}
  ol "Detected Rails application"

  if has_gem pg; then
    ol "Detected the pg gem, installing PostgreSQL..."
    . /etc/default/locale
    oe sudo apt-get install -y postgresql-9.1
    oe sudo -u postgres createuser --superuser vagrant
  fi

  ol "Preparing the database..."
  oe "bundle exec rake db:setup || bundle exec rake db:migrate"
{% endif %}

SCRIPT
{% endblock %}
