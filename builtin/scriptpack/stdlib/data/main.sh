# Setup all the error cases for scripts. We do this early hoping this is
# the first thing loaded.
set -o nounset -o errexit -o pipefail -o errtrace

# Load all our functions
. ${SCRIPTPACK_STDLIB_ROOT}/execute.sh
. ${SCRIPTPACK_STDLIB_ROOT}/error.sh
. ${SCRIPTPACK_STDLIB_ROOT}/output.sh

# Setup a trap for all errors to be logged
trap 'otto_error "${BASH_SOURCE}" "${LINENO}"' ERR
