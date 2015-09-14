apt-get -y install open-vm-tools
echo -n ".host:/ /mnt/hgfs vmhgfs rw,ttl=1,uid=my_uid,gid=my_gid,nobootwait 0 0" >> /etc/fstab
