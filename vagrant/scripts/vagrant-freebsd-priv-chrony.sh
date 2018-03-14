#!/bin/sh

set -o errexit

export ASSUME_ALWAYS_YES=yes

pkg update
pkg install -y net/chrony

cat <<EOT >> /usr/local/etc/chrony.conf
pool 0.freebsd.pool.ntp.org iburst
 
minsources 2
driftfile /var/db/chrony/drift
makestep 1.0 3
EOT

echo 'chronyd_enable="YES"' >> /etc/rc.conf

service chronyd start
