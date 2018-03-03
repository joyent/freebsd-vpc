#!/bin/sh

if ! zfs list guests/vm > /dev/null 2> /dev/null ; then
	zfs create -o compression=lz4 -o logbias=throughput guests/vm
fi

echo 'net.link.tap.up_on_open=1' >> /etc/sysctl.conf
echo 'vmm_load="YES"' >> /boot/loader.conf
echo 'nmdm_load="YES"' >> /boot/loader.conf
echo 'if_bridge_load="YES"' >> /boot/loader.conf
echo 'if_tap_load="YES"' >> /boot/loader.conf
