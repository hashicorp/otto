# Setup all the error cases for scripts. We do this early hoping this is
# the first thing loaded.

# Load all our functions
. ${SCRIPTPACK_STDLIB_ROOT}/execute.sh
. ${SCRIPTPACK_STDLIB_ROOT}/error.sh
. ${SCRIPTPACK_STDLIB_ROOT}/output.sh
. ${SCRIPTPACK_STDLIB_ROOT}/vagrant.sh

# Ubuntu
. ${SCRIPTPACK_STDLIB_ROOT}/apt.sh

# Function for initializing any scripts. We put this in a function because
# our BATS tests don't call this (it messes them up).
otto_init() {
  set -o nounset -o errexit -o pipefail -o errtrace
}

# Setup a trap for all errors to be logged
trap 'otto_error "${BASH_SOURCE}" "${LINENO}"' ERR
