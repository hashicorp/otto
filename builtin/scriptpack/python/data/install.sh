# python_install_prepare should be called before any installation is done.
python_install_prepare() {
  # Obvious: we're running non-interactive mode
  export DEBIAN_FRONTEND=noninteractive

  # Update apt once
  apt_update_once

  # Our PPAs have unicode characters, so we need to set the proper lang.
  otto_init_locale

  oe sudo apt-get install -y python-software-properties software-properties-common \
      apt-transport-https wget
}

# python_install installs an arbitrary Python version given in the argument.
#
# The python version must be two sections, ex. "2.4", "2.7"
python_install() {
    local version="$1"
    case $version in
    2.3)
        _python_install_raw "2.3"
        ;;
    2.4)
        _python_install_raw "2.4"
        ;;
    2.5)
        _python_install_raw "2.5"
        ;;
    2.6)
        _python_install_raw "2.6"
        ;;
    2.7)
        _python_install_raw "2.7"
        ;;
    3.1)
        _python_install_raw "3.1"
        ;;
    3.2)
        _python_install_raw "3.2"
        ;;
    3.3)
        _python_install_raw "3.3"
        ;;
    3.4)
        _python_install_raw "3.4"
        ;;
    3.5)
        _python_install_raw "3.5"
        ;;
    *)
        echo "Unknown Python version: ${version}" >&2
        exit 1
        ;;
    esac
}

_python_install_raw() {
    local version="$1"

    # Install Python proper
    python_install_prepare
    apt_install python${version} python${version}-dev

    # Install pip
    wget -q -O - https://bootstrap.pypa.io/get-pip.py | oe sudo python${version}

    # Install virtualenv
    oe sudo pip install virtualenv
}
