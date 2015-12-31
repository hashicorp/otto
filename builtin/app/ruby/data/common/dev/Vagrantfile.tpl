{% extends "compile:data/app/dev/Vagrantfile.tpl" %}

{% block vagrant_config %}
  config.vm.provision "shell", inline: $script_app, privileged: false
{% endblock %}

{% block footer %}
$script_app = <<SCRIPT
. /otto/scriptpacks/STDLIB/main.sh
. /otto/scriptpacks/RUBY/main.sh
otto_init

# Make it so that `vagrant ssh` goes directly to the correct dir
vagrant_default_cd "vagrant" "/vagrant"

# Go to our working directory and install gems
cd /vagrant
ruby_gemfile_apt
otto_output "Bundling gem dependencies..."
bundle

{% if app_type == "rails" %}
  otto_output "Detected Rails application"

  if has_gem pg; then
    otto_output "Detected the pg gem, installing PostgreSQL..."
    . /etc/default/locale
    oe sudo apt-get install -y postgresql-9.1
    oe sudo -u postgres createuser --superuser vagrant
  fi

  otto_output "Preparing the database..."
  oe "bundle exec rake db:setup || bundle exec rake db:migrate"
{% endif %}

SCRIPT
{% endblock %}
