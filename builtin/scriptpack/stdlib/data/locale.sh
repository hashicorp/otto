# otto_init_locale sets the locale up for UTF-8
otto_init_locale() {
  otto_output "Setting locale to en_US.UTF-8..."
  if [[ ! $(locale -a) =~ '^en_US\.utf8' ]]; then
      oe sudo locale-gen en_US.UTF-8
  fi
  oe sudo update-locale LANG=en_US.UTF-8 LC_ALL=en_US.UTF-8
  export LANG=en_US.UTF-8
}
