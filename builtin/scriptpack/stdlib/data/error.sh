# otto_error is called when an error occurs. The trap is configured
# globally from main.sh.
otto_error() {
   local sourcefile=$1
   local lineno=$2
   echo "ERROR at ${sourcefile}:${lineno}; Last logs:"
   grep otto /var/log/syslog | tail -n 20
}
