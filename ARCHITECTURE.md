# Virtual Private Cloud (VPC) Architecture

## `syscall(2)` interface

```c
typedef vpc_handle_t int32_t;
typedef vpc_id_t     struct uuid;
typedef struct {
  uint64_t version:4;
  uint64_t pad1:4;
  uint64_t obj_type:8;
  uint64_t pad2:48;
} vpc_handle_type_t;
typedef vpc_flags_t  uint64_t;
typedef vpc_op_t     uint32_t;
typedef vpc_txn_t    int32_t;
typedef struct {
  char srcmac[6];
  char dstmac[6];
  int ttl_ms;
} vpc_sw_cam_t;
typedef vpc_swp_flags_t uint64_t;

const vpc_type_t VPC_TYPE_UNSET = 0;
const vpc_type_t VPC_TYPE_NIC    = 1;
const vpc_type_t VPC_TYPE_SWITCH = 2;
const vpc_type_t VPC_TYPE_SWPORT = 3;
const vpc_type_t VPC_TYPE_ROUTER = 4;
const vpc_type_t VPC_TYPE_NAT    = 5;
const vpc_type_t VPC_TYPE_LINK   = 6;

const vpc_flag_t VPC_FLAG_DEFAULT = 1 << 0;
const vpc_flag_t VPC_FLAG_CREATE  = 1 << 1;
const vpc_flag_t VPC_FLAG_OPEN    = 1 << 2;

const vpc_op_t VPC_OP_NOP = 0;
const vpc_op_t VPC_OP_GET_TXN = 1;
const vpc_op_t VPC_OP_NIC_SET_MAC = 2;

typedef struct {
  vpc_txn_t txn_t;
  u_char    mac[ETHER_ADDR_LEN];
} vpc_vmnic_t;

typedef struct {
  vpc_txn_t txn_t;
  int num_ports;
  vpc_handle_t vpcpd;
} vpc_switch_t;

const vpc_swp_flags_t VPC_FLAG_SWP_DEFAULT    = 1 << 0;
const vpc_swp_flags_t VPC_FLAG_SWP_UPLINK     = 1 << 1;
const vpc_swp_flags_t VPC_FLAG_SWP_ROUTER     = 1 << 2;
const vpc_swp_flags_t VPC_FLAG_SWP_VMNIC      = 1 << 3;
const vpc_swp_flags_t VPC_FLAG_SWP_VLAN_ID    = 1 << 4;
const vpc_swp_flags_t VPC_FLAG_SWP_VNI        = 1 << 5;
const vpc_swp_flags_t VPC_FLAG_SWP_FW_EGRESS  = 1 << 6;
const vpc_swp_flags_t VPC_FLAG_SWP_FW_INGRESS = 1 << 7;
const vpc_swp_flags_t VPC_FLAG_SWP_DESCRIPTOR = 1 << 8;

typedef struct {
  vpc_txn_t txn_t;
  u_char smac[ETHER_ADDR_LEN];
  int32_t vni;
  int16_t vlan_id;
  void *descriptor_or_pointer_to_ipfw_state_table;
  void *descriptor_or_pointer_to_ipfw_ruleset_for_egress_flows;
  vpc_swp_flags_t flags;
  vpc_handle_t descriptor;
} vpc_swport_t;

vpc_handle_t
vpc_open(const vpc_id_t *vpc_id,
         vpc_type_t obj_type, vpc_flags_t flags);

int
vpc_ctl(vpc_handle_t vpcd, vpc_op_t op,
        size_t keylen, const void *key,
        size_t *vallen, void *buf);
```

## Types

1. `vpcnic`
2. `vpcsw`
3. `vpcp`
4. `vpcr`
5. `vpcnat`
6. `vpclink`

### `vmnic`

* 1x IP on underlay per NUMA domain (i.e. 1x `vcc(4)` per NUMA domain, however
  initially only one IP for the entire system).
* 1x IP per customer account for pulic EIP.
* 1x IP per facility

```c
vpc_handle_t nicd;
{
  id := UUID{"a-new-nic-uuid"}
  obj_type := vpc_type_t{
      version: 1,
      obj_type: VPC_TYPE_NIC,
  };
  nicd = vpc_open(id, obj_type, VPC_FLAG_CREATE);
  if (nicd < 0) err(...);
}

{
  vpc_txn_t txn;
  int err;

  err = vpc_ctl(nicd, VPC_OP_GET_TXN, NULL,
      0, NULL, sizeof(txn), &txn);

  vpc_vmnic_mac_t nic_mac;
  bzero(&nic_mac, sizeof(nic_mac));
  nic_mac.mac = 0x0;
  nic_mac.txn = txn;

  err = vpc_ctl(nicd, VPC_OP_NIC_SET_MAC,
          sizeof(mac), &nic_mac, 0, NULL);
  if (err < 0) err(...);
}
```

Returns for `vpc_open()`:

* `0` if no error
* `EEXIST` if the UUID on `VPC_FLAG_CREATE` is inuse.
* `ENOENT` if the UUID on `VPC_FLAG_OPEN` does not exist.

### `vpcsw`

* `vpcsw` is a per-VPC artifact.  Each `vpcsw` can be configured with one and
  only one VNI.
* `vpcsw`s configured with a VNI of `0` will never send VXLAN-encapsulated
  traffic but may send VLAN-tagged traffic the wire.
* `vpcsw` creates its own switch ports, `vpcp`
* `vpcsw` has different port types, one for each of the different types of
  objects that plug into the switch.
* Each `vpcsw` can have one port designated as an `uplink` port.  An `uplink`
  port is used to drain `mvec`s to either the physical device (or a `vpclink`).
  `mvec`s who have a resolved `dmac` can egress the `vpcsw` using the `uplink`
  port.

  If the `vpcsw`:

  1. receives an `mvec` from a `vpcp` handling a leaf; AND
  2. does not have a forwarding entry for a given `dmac`.

  the `vpcsw` will perform an upcall to resolve the IP address of the
  destination CN on the underlay network.
* Each `vpcsw` can have one port designated as the `router` port.  The `router`
  port is the implicit next hop for `mvec`s that have an overlay destination IP
  that is not part of the same subnet as the overlay source IP.


```c
vpc_handle_t swd;
{
  id := UUID{"a-new-sw-uuid"}
  obj_type := vpc_type_t{
    version: 1,
    obj_type: VPC_TYPE_SWITCH,
  };
  swd = vpc_open(id, obj_type, VPC_FLAG_CREATE);
  if (swd < 0) err(...);
}

retry = 1;
do {
  vpc_txn_t txn;
  int err;

  err = vpc_ctl(swd, VPC_OP_GET_TXN, NULL,
      0, NULL, sizeof(txn), &txn);

  vpc_switch_t sw;
  bzero(&sw, sizeof(sw));
  sw.txn = txn;

  err = vpc_ctl(nicd, VPC_OP_SW_GET_PORT_COUNT,
      0, NULL, sizeof(sw), &sw);
  if (err < 0) {
    err(...);
    continue;
  }

  if (sw.num_ports == 0) {
    /* Add an uplink port */
    vpc_swport_t uplink_port_req;
    vpc_swport_t uplink_port_res;

    bzero(&uplink_port_req, sizeof(uplink_port_req));
    uplink_port_req.txn = txn;
    uplink_port_req.flags = VPC_FLAG_SWP_UPLINK;

    err = vpc_ctl(swd, VPC_OP_SW_ADD_PORT,
        sizeof(uplink_port_req), &uplink_port_req,
        sizeof(uplink_port_res), &uplink_port_res);
    if (err < 0) {
      err(...);
      continue;
    }

    /* Add a router port, assuming there's a router */
    vpc_swport_t router_port_req;
    vpc_swport_t router_port_res;

    bzero(&router_port_req, sizeof(router_port_req));
    router_port_req.txn = txn;
    router_port_req.flags = VPC_FLAG_SWP_ROUTER;

    err = vpc_ctl(swd, VPC_OP_SW_ADD_PORT,
        sizeof(router_port_req), &router_port_req,
        sizeof(router_port_res), &router_port_res);
    if (err < 0) {
      err(...);
      continue;
    }
  }

  /* Add a `vmnic` port */
  vpc_swport_t vmnic_port_req;
  vpc_swport_t vmnic_port_res;

  bzero(&vmnic_port_req, sizeof(vmnic_port_req));
  vmnic_port_req.txn = txn;
  vmnic_port_req.flags = VPC_FLAG_SWP_VMNIC;
  vmnic_port_req.vlan_id = 123;
  vmnic_port_req.flags |= VPC_FLAG_SWP_VLAN_ID;
  vmnic_port_req.vni = 456;
  vmnic_port_req.flags |= VPC_FLAG_SWP_VNI;

  swpd = vpc_ctl(swd, VPC_OP_SW_ADD_PORT,
      sizeof(vmnic_port_req), &vmnic_port_req,
      sizeof(vmnic_port_res), &vmnic_port_res);
  if (err < 0) {
    err(...);
    continue;
  }

  vpc_handle_t vmnicd = 0;
  if (vmnic_port_res.flags & VPC_FLAG_SWP_DESCRIPTOR != 0) {
    vmnicd = vmnic_port_res.descriptor;
  }

  retry = 0;
  break
} while(retry)
```
