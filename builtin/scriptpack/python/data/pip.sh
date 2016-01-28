# python_requirements_contains checks if the requriements.txt file
# has a certain pip package in it
python_requirements_contains() {
    local name=$1

    if [ -f requirements.txt ]; then
        grep -e "$name" requirements.txt > /dev/null
        return $?
    fi

    return 1
}

# python_pip_apt installs packages for pip packages that are detected.
python_pip_apt() {
    _python_pip_packages=()
    _python_requirements_check psycopg2 "libpq-dev"

    if [ -n "${_python_pip_packages-}" ]; then
        otto_output "Installing native pip package system dependencies..."
        apt_update_once
        apt_install "${_python_pip_packages[@]}"
    fi
}

# Internal functions for accumulating the queue of things to install
# for a pip package.
_python_pip_packages=()
_python_requirements_check() {
    local package=$1
    local deps=$2

    if python_requirements_contains $package; then
        otto_output "Detected the pip package: ${package}"
        _python_pip_packages+=($deps)
    fi
}
