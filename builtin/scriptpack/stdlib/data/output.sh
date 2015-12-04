# otto_output can be used to write output to the screen that is clearly
# marked as coming from Otto. It is recommended for all messages to the user.
otto_output() {
    echo "[otto] $@"
}
