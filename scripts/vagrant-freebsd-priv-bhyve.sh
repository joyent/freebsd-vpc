#!/bin/sh

if ! zfs list guests/vm > /dev/null 2> /dev/null ; then
	zfs create -o compression=lz4 -o logbias=throughput guests/vm
fi
