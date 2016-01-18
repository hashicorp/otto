$preshell = <<SCRIPT
if command -v apt-get >/dev/null 2>&1; then
    export DEBIAN_FRONTEND=noninteractive
    sudo apt-get install -y apt-transport-https >/dev/null 2>&1
fi
SCRIPT

# http://foo-o-rama.com/vagrant--stdin-is-not-a-tty--fix.html
config.vm.provision "fix-no-tty", type: "shell" do |s|
    s.privileged = false
    s.inline = "sudo sed -i '/tty/!s/mesg n/tty -s \\&\\& mesg n/' /root/.profile"
end

config.vm.provision "shell", inline: $preshell
config.vm.provision "docker" do |d|
  d.run "{{ name }}", args: "{{ run_args }}", image: "{{ docker_image }}"
end

# Sync our own dep folder in there
config.vm.synced_folder '{{ path.working }}', "{{ path.guest_working }}"

# Foundation configuration for dev dep
{% for dir in foundation_dirs.dev_dep %}
dir = "/otto/foundation-{{ name }}-{{ forloop.Counter }}"
config.vm.synced_folder '{{ dir }}', dir
config.vm.provision "shell", inline: "cd #{dir} && bash #{dir}/main.sh"
{% endfor %}
