#!/usr/bin/env bats

# Load the main library
. ${SCRIPTPACK_PYTHON_ROOT}/main.sh

@test "pip packages with a native extension" {
  # Write a test Gemfile
  local dir=$(mktemp -d)
  cd $dir
  cat <<EOF >requirements.txt
psycopg2==2.6.1
EOF

  # Install deps
  python_pip_apt

  # Check for it
  dpkg-query -l libpq-dev
}
