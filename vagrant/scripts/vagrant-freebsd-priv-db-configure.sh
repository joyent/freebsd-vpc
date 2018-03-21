#!/bin/sh

set -o errexit

node_ip=$(ifconfig vmx0 | grep inet | awk '{ print $2; }')

certs_dir="/secrets/crdb"
ca_key_path="/secrets/crdb_ca"
ca_key="${ca_key_path}/ca.key"

mkdir -p "${certs_dir}"
mkdir -p "${ca_key_path}"

cockroach cert create-ca --certs-dir "${certs_dir}" --ca-key="${ca_key}"
cockroach cert create-node localhost "0.0.0.0" "${node_ip}" --certs-dir "${certs_dir}" --ca-key="${ca_key}"
cockroach cert create-client root --certs-dir "${certs_dir}" --ca-key="${ca_key}"

cockroach_flags="--certs-dir ${certs_dir}"
cockroach_flags="${cockroach_flags} --host=0.0.0.0 --advertise-host=${node_ip}"

echo 'cockroach_enable="YES"' >> /etc/rc.conf
echo "cockroach_flags=\"${cockroach_flags}\"" >> /etc/rc.conf

mkdir /home/vagrant/.cockroach-certs

chown -R cockroach:cockroach /secrets/crdb
chmod 0700 /secrets/crdb/ca.crt
chmod 0700 /secrets/crdb/node.crt
chmod 0700 /secrets/crdb/node.key

cp /secrets/crdb/ca.crt /home/vagrant/.cockroach-certs/ca.crt
cp /secrets/crdb/client.root.crt /home/vagrant/.cockroach-certs
cp /secrets/crdb/client.root.key /home/vagrant/.cockroach-certs
chown -R vagrant:vagrant /home/vagrant/.cockroach-certs
chmod 0700 /home/vagrant/.cockroach-certs/ca.crt
chmod 0700 /home/vagrant/.cockroach-certs/client.root.crt
chmod 0700 /home/vagrant/.cockroach-certs/client.root.key

service cockroach start

sleep 5
cockroach sql \
	--user root \
	--certs-dir /home/vagrant/.cockroach-certs \
	--execute "CREATE DATABASE IF NOT EXISTS triton"
