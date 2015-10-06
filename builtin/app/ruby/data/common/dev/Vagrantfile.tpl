# Generated by Otto, do not edit!
#
# This is the Vagrantfile generated by Otto for the development of
# this application/service. It should not be hand-edited. To modify the
# Vagrantfile, use the Appfile.

Vagrant.configure("2") do |config|
  config.vm.box = "hashicorp/precise64"

  # Host only network
  config.vm.network "private_network", ip: "{{ dev_ip_address }}"

  # Setup a synced folder from our working directory to /vagrant
  config.vm.synced_folder '{{ path.working }}', "/vagrant",
    owner: "vagrant", group: "vagrant"

  # Enable SSH agent forwarding so getting private dependencies works
  config.ssh.forward_agent = true

  # Foundation configuration (if any)
  {% for dir in foundation_dirs.dev %}
  dir = "/otto/foundation-{{ forloop.Counter }}"
  config.vm.synced_folder '{{ dir }}', dir
  config.vm.provision "shell", inline: "cd #{dir} && bash #{dir}/main.sh"
  {% endfor %}

  # Load all our fragments here for any dependencies.
  {% for fragment in dev_fragments %}
  {{ fragment|read }}
  {% endfor %}

  # Set locale to en_US.UTF-8
  config.vm.provision "shell", inline: $script_locale

  # Install Ruby build environment
  config.vm.provision "shell", inline: $script_ruby, privileged: false

  config.vm.provider :parallels do |p, o|
    o.vm.box = "parallels/ubuntu-12.04"
  end
end

$script_locale = <<SCRIPT
  oe() { eval "$@" 2>&1 | logger -t otto > /dev/null; }
  ol() { echo "[otto] $@"; }

  ol "Setting locale to en_US.UTF-8..."
  oe locale-gen en_US.UTF-8
  oe update-locale LANG=en_US.UTF-8 LC_ALL=en_US.UTF-8
SCRIPT

$script_ruby = <<SCRIPT
set -o nounset -o errexit -o pipefail -o errtrace

error() {
   local sourcefile=$1
   local lineno=$2
   echo "ERROR at ${sourcefile}:${lineno}; Last logs:"
   grep otto /var/log/syslog | tail -n 20
}
trap 'error "${BASH_SOURCE}" "${LINENO}"' ERR

# otto-exec: execute command with output logged but not displayed
oe() { eval "$@" 2>&1 | logger -t otto > /dev/null; }

# otto-log: output a prefixed message
ol() { echo "[otto] $@"; }

# Make it so that `vagrant ssh` goes directly to the correct dir
echo "cd /vagrant" >> /home/vagrant/.bashrc

# Configuring SSH for faster login
if ! grep "UseDNS no" /etc/ssh/sshd_config >/dev/null; then
  echo "UseDNS no" | sudo tee -a /etc/ssh/sshd_config >/dev/null
  oe sudo service ssh restart
fi

export DEBIAN_FRONTEND=noninteractive

ol "Adding apt repositories and updating..."
oe sudo apt-get update
oe sudo apt-get install -y python-software-properties software-properties-common apt-transport-https
oe sudo add-apt-repository -y ppa:chris-lea/node.js
oe sudo apt-add-repository -y ppa:brightbox/ruby-ng
oe sudo apt-get update

# TODO: parameterize ruby version as input
export RUBY_VERSION="{{ ruby_version }}"

ol "Installing Ruby ${RUBY_VERSION} and supporting packages..."
export DEBIAN_FRONTEND=noninteractive
oe sudo apt-get install -y bzr git mercurial build-essential \
  software-properties-common \
  nodejs \
  ruby$RUBY_VERSION ruby$RUBY_VERSION-dev

ol "Configuring Ruby environment..."
echo 'export GEM_HOME=$HOME/.gem\nexport PATH=$HOME/.gem/bin:$PATH' >> $HOME/.ruby_env
echo 'source $HOME/.ruby_env' >> $HOME/.bashrc
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
detect_gem_deps mysql2 "libmysqld-dev"
detect_gem_deps nokogiri "zlib1g-dev"
detect_gem_deps pg "libpq-dev"
detect_gem_deps rmagick "libmagickwand-dev"
detect_gem_deps sqlite3 "libsqlite3-dev"
detect_gem_deps libxml-ruby "libxml-dev"

if [ -n "${gem_deps_queue-}" ]; then
  ol "Installing native gem system dependencies..."
  oe sudo apt-get install -y "${gem_deps_queue[@]}"
fi

ol "Installing Bundler..."
oe gem install bundler --no-document

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

ol "Configuring Git to use SSH instead of HTTP so we can agent-forward private repo auth..."
oe git config --global url."git@github.com:".insteadOf "https://github.com/"
SCRIPT
