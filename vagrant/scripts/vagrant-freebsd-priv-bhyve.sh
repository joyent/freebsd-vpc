#!/bin/sh

set -o errexit

if ! zfs list guests/vm > /dev/null 2> /dev/null ; then
	zfs create -o compression=lz4 -o logbias=throughput guests/vm
fi

echo 'net.link.tap.up_on_open=1' >> /etc/sysctl.conf
sysrc -f /boot/loader.conf 'vmm_load=YES'
sysrc -f /boot/loader.conf 'vmmnet_load=YES'
sysrc -f /boot/loader.conf 'nmdm_load=YES'
sysrc -f /boot/loader.conf 'if_bridge_load=YES'
sysrc -f /boot/loader.conf 'if_tap_load=YES'
