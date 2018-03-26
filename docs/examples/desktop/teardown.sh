#!/bin/sh --

set -x
set -e
set -u

. input.sh

vpc ethlink destroy --ethlink-id=${ETHLINK0_ID}
vpc switch port disconnect --port-id=${VPCP0_ID} --interface-id=${HOSTLINK_ID}
vpc vmnic destroy --vmnic-id=${HOSTLINK_ID}
vpc switch port remove --switch-id=${VPCSW0_ID} --port-id=${VPCP0_ID}
vpc switch destroy --switch-id=${VPCSW0_ID}

vpc list
