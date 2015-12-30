# vagrant_config_fast_ssh configures SSH to be faster.
vagrant_config_fast_ssh() {
  if ! grep "UseDNS no" /etc/ssh/sshd_config >/dev/null; then
    echo "UseDNS no" | sudo tee -a /etc/ssh/sshd_config >/dev/null
    oe sudo service ssh restart
  fi
}
