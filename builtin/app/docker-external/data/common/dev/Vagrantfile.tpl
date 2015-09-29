Vagrant.configure("2") do |config|
  config.vm.box = "hashicorp/precise64"

  # Host only network
  config.vm.network "private_network", ip: "{{ dev_ip_address }}"

  # Setup some stuff
  config.vm.provision "shell", inline: $script

  # Foundation configuration (if any)
  {% for dir in foundation_dirs.dev %}
  dir = "/otto/foundation-{{ forloop.Counter }}"
  config.vm.synced_folder "{{ dir }}", dir
  config.vm.provision "shell", inline: "cd #{dir} && bash #{dir}/main.sh"
  {% endfor %}

  # Read in the fragment that we use as a dep
  eval(File.read("{{ fragment_path }}"), binding)

  ["vmware_fusion", "vmware_workstation"].each do |name|
    config.vm.provider(name) do |p|
      p.enable_vmrun_ip_lookup = false
    end
  end

  config.vm.provider :parallels do |p, o|
    o.vm.box = "parallels/ubuntu-12.04"
  end
end

$script = <<SCRIPT
set -e

# otto-exec: execute command with output logged but not displayed
oe() { $@ 2>&1 | logger -t otto > /dev/null; }

# otto-log: output a prefixed message
ol() { echo "[otto] $@"; }

# Configuring SSH for faster login
if ! grep "UseDNS no" /etc/ssh/sshd_config >/dev/null; then
  echo "UseDNS no" | sudo tee -a /etc/ssh/sshd_config >/dev/null
  oe sudo service ssh restart
fi

export DEBIAN_FRONTEND=noninteractive

ol "Installing HTTPS driver for Apt..."
oe sudo apt-get update
oe sudo apt-get install -y apt-transport-https
SCRIPT
