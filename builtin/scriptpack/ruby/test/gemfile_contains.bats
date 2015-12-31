#!/usr/bin/env bats

# Load the main library
. ${SCRIPTPACK_RUBY_ROOT}/main.sh

@test "gemfile contains a gem" {
  # Write a test Gemfile
  local dir=$(mktemp -d)
  cd $dir
  cat <<EOF >Gemfile
gem "nokogiri"
EOF

  if ! ruby_gemfile_contains "nokogiri"; then
      echo "nokogiri not found"
      exit 1
  fi

  if ruby_gemfile_contains "foo"; then
      echo "foo should not be found"
      exit 1
  fi
}
