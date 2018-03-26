#!/bin/sh --

set -x
set -e
set -u

. input.sh

vpc switch create --switch-id=${VPCSW0_ID} --vni=${VNI}

vpc hostlink create --hostlink-id=${HOSTLINK_ID}
vpc switch port add --switch-id=${VPCSW0_ID} --port-id=${VPCP0_ID}
vpc switch port connect --port-id=${VPCP0_ID} --interface-id=${HOSTLINK_ID}

vpc switch port add --switch-id=${VPCSW0_ID} --port-id=${UPLINK_PORT_ID} --uplink --l2-name=${UPLINK_IF} --ethlink-id=${ETHLINK0_ID}

vpc vmnic get --vmnic-id=${VMNIC0_ID}

vpc list
