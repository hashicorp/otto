#!/usr/bin/env bats

# Load the main library
. ${SCRIPTPACK_PHP_ROOT}/main.sh

@test "listing PHP versions should look like versions" {
  local list=($(php_version_list))
  local fail=0
  for i in "${list[@]}"; do
    if [[ ! $i =~ ^[0-9] ]]; then
        fail=1
        echo $i
    fi
  done

  [ $fail -eq 0 ]
}
