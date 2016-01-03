#!/usr/bin/env bats

# Load the main library
. ${SCRIPTPACK_SKELETON_ROOT}/main.sh

@test "foo function" {
  output=$(skeleton_foo)
  [[ $output =~ 'hello' ]]
}
