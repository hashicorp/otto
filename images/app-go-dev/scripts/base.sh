# Basic packages
apt-get update
apt-get -y install linux-headers-$(uname -r)

# Passwordless sudo
sed -i -e '/Defaults\s\+env_reset/a Defaults\texempt_group=sudo' /etc/sudoers
sed -i -e 's/%sudo  ALL=(ALL:ALL) ALL/%sudo  ALL=NOPASSWD:ALL/g' /etc/sudoers

# Faster SSH sign-in
echo "UseDNS no" >> /etc/ssh/sshd_config
