#!/bin/sh

set -o errexit

echo 'dbus_enable="YES"' >> /etc/rc.conf
echo 'avahi_daemon_enable="YES"' >> /etc/rc.conf

sed -i '' 's/^hosts: files dns$/hosts: files dns mdns/' /etc/nsswitch.conf

service dbus start
service avahi-daemon start
