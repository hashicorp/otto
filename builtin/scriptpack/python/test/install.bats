#!/usr/bin/env bats

# Load the main library
. ${SCRIPTPACK_PYTHON_ROOT}/main.sh

@test "install from string '2.7'" {
  python_install "2.7"

  # Verify Python installed
  python --version
  [[ $(python --version 2>&1) =~ 'Python 2.7.' ]]

  # Verify pip installed
  pip --version
  [[ $(pip --version 2>&1) =~ 'pip' ]]

  # Verify virtualenv
  virtualenv --version
  [[ $(virtualenv --version 2>&1) =~ '13.' ]]
}

@test "install from string '3.5'" {
  python_install "3.5"

  # Verify Python installed
  python --version
  [[ $(python --version 2>&1) =~ 'Python 3.5.' ]]

  # Verify pip installed
  pip --version
  [[ $(pip --version 2>&1) =~ 'pip' ]]

  # Verify virtualenv
  virtualenv --version
  [[ $(virtualenv --version 2>&1) =~ '13.' ]]
}
