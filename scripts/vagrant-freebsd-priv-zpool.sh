#!/bin/sh

if zpool import | grep -q guests ; then
	zpool import guests
else
	zpool create guests /dev/da1
fi
