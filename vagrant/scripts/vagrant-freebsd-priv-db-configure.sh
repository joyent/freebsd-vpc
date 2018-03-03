#!/bin/sh

node_internal=$(ifconfig vmx1 | grep inet | awk '{ print $2; }')
join_addrs="172.27.10.11:26257,172.27.10.12:26257,172.27.10.13:26257"

cockroach_flags="--certs-dir /secrets/crdb"
cockroach_flags="${cockroach_flags} --advertise-host=${node_internal}"
cockroach_flags="${cockroach_flags} --join=${join_addrs}"

echo 'cockroach_enable="YES"' >> /etc/rc.conf
echo "cockroach_flags=\"${cockroach_flags}\"" >> /etc/rc.conf

cat <<EOT >> /home/vagrant/.profile
export COCKROACH_HOST=${node_internal}
EOT

chown -R cockroach:cockroach /secrets/crdb

service cockroach start
