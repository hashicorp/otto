# oe is shorthand for "otto_execute" and executes a command with no
# output. The output is logged to the syslog. This output can be read
# later but is also expected to be picked up from otto_error.
oe() {
    "$@" 2>&1 | logger -t otto >/dev/null
}
