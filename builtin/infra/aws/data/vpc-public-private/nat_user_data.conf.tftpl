#cloud-config
apt_upgrade: true
locale: en_US.UTF-8
runcmd:
 - [ sh, -c, "echo 1 > /proc/sys/net/ipv4/ip_forward;echo 655361 > /proc/sys/net/netfilter/nf_conntrack_max" ]
 - [ iptables, -N, LOGGINGF ]
 - [ iptables, -N, LOGGINGI ]
 - [ iptables, -A, LOGGINGF, -m, limit, --limit, 2/min, -j, LOG, --log-prefix, "IPTables-FORWARD-Dropped: ", --log-level, 4 ]
 - [ iptables, -A, LOGGINGI, -m, limit, --limit, 2/min, -j, LOG, --log-prefix, "IPTables-INPUT-Dropped: ", --log-level, 4 ]
 - [ iptables, -A, LOGGINGF, -j, DROP ]
 - [ iptables, -A, LOGGINGI, -j, DROP ]
 - [ iptables, -A, FORWARD, -s, ${vpc_cidr}, -j, ACCEPT ]
 - [ iptables, -A, FORWARD, -j, LOGGINGF ]
 - [ iptables, -P, FORWARD, DROP ]
 - [ iptables, -I, FORWARD, -m, state, --state, "ESTABLISHED,RELATED", -j, ACCEPT ]
 - [ iptables, -t, nat, -I, POSTROUTING, -s, ${vpc_cidr}, -d, 0.0.0.0/0, -j, MASQUERADE ]
 - [ iptables, -A, INPUT, -s, ${vpc_cidr}, -j, ACCEPT ]
 - [ iptables, -A, INPUT, -p, tcp, --dport, 22, -m, state, --state, NEW, -j, ACCEPT ]
 - [ iptables, -I, INPUT, -m, state, --state, "ESTABLISHED,RELATED", -j, ACCEPT ]
 - [ iptables, -I, INPUT, -i, lo, -j, ACCEPT ]
 - [ iptables, -A, INPUT, -j, LOGGINGI ]
 - [ iptables, -P, INPUT, DROP ]
