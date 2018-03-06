#!/bin/sh --

set -x
set -e
set -u

. 2vm_input.sh

vpc ethlink destroy --ethlink-id=${ETHLINK0_ID}
vpc switch port disconnect --port-id=${VPCP0_ID} --interface-id=${VMNIC0_ID}
vpc switch port disconnect --port-id=${VPCP1_ID} --interface-id=${VMNIC1_ID}
vpc vmnic destroy --vmnic-id=${VMNIC0_ID}
vpc vmnic destroy --vmnic-id=${VMNIC1_ID}
vpc switch port remove --switch-id=${VPCSW0_ID} --port-id=${VPCP0_ID}
vpc switch port remove --switch-id=${VPCSW1_ID} --port-id=${VPCP1_ID}
#vpc switch port remove --switch-id=${VPCSW0_ID} --port-id=${UPLINK_PORT_ID}
vpc switch destroy --switch-id=${VPCSW0_ID}

vpc list
