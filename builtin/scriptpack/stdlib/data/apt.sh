# apt_update_once can be called to run apt_update exactly once. This
# will create a temporary file in /tmp to prevent it from happening again.
apt_update_once() {
    if [ ! -f "/tmp/otto_apt_update_sentinel" ]; then
        apt_update
    fi
}

# apt_update updates the apt cache
apt_update() {
    sudo apt-get update
}
