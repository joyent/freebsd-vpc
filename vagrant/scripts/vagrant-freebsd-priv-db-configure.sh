#!/bin/sh

set -o errexit

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
chmod 0700 /secrets/crdb/ca.crt
chmod 0700 /secrets/crdb/node.crt
chmod 0700 /secrets/crdb/node.key
chmod 0700 /home/vagrant/.cockroach-certs/ca.crt
chmod 0700 /home/vagrant/.cockroach-certs/client.root.crt
chmod 0700 /home/vagrant/.cockroach-certs/client.root.key

service cockroach start
