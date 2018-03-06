#!/bin/sh --

set -x
set -e
set -u

. 2vm_input.sh

vpc switch create --switch-id=${VPCSW0_ID} --vni=${VNI}

vpc vmnic create --vmnic-id=${VMNIC1_ID}
vpc switch port add --switch-id=${VPCSW0_ID} --port-id=${VPCP0_ID}
vpc switch port connect --port-id=${VPCP0_ID} --interface-id=${VMNIC0_ID}

vpc vmnic create --vmnic-id=${VMNIC0_ID}
vpc switch port add --switch-id=${VPCSW0_ID} --port-id=${VPCP1_ID}
vpc switch port connect --port-id=${VPCP1_ID} --interface-id=${VMNIC1_ID}

vpc switch port add --switch-id=${VPCSW0_ID} --port-id=${UPLINK_PORT_ID} --uplink --l2-name=${UPLINK_IF} --ethlink-id=${ETHLINK0_ID}

vpc vmnic get --vmnic-id=${VMNIC0_ID}
vpc vmnic get --vmnic-id=${VMNIC1_ID}

vpc list
