#!/usr/bin/env bats

# Load the main library
. ${SCRIPTPACK_RUBY_ROOT}/main.sh

@test "gemfile with a native extension" {
  # Write a test Gemfile
  local dir=$(mktemp -d)
  cd $dir
  cat <<EOF >Gemfile
gem "curb"
EOF

  # Install deps
  ruby_gemfile_apt

  # Check for it
  dpkg-query -l libcurl3
}
