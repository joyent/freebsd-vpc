-- VPC has a few different dimensions:
--
-- * Administrative Plane (customer of the VPC)
-- * Control Plane (customer of the VPC)
-- * Data Plane (users of a deployed service)

-- Organization ("Org") is an administrative construct.  An organization is not
-- referenceable by any control or data plane entity.
CREATE TABLE IF NOT EXISTS org (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT
);

-- Accounts are the primary unit of granularity for VPC objects.  An Account is
-- created by an Organization.  An Account is the primary object owner for
-- control plane entities.
CREATE TABLE IF NOT EXISTS account (
  id UUID NOT NULL DEFAULT gen_random_uuid(),
  org_id UUID NOT NULL,
  name TEXT,
  INDEX(org_id),
  PRIMARY KEY(org_id, id),
  UNIQUE(id),
  CONSTRAINT org_id_fk FOREIGN KEY(org_id) REFERENCES org(id)
) INTERLEAVE IN PARENT org(org_id);

-- VPC is a collection of subnets and services.  The mapping of VPCs to subnets
-- is handled by VXLAN IDs and VLANs, respectively.
CREATE TABLE IF NOT EXISTS vpc (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  account_id UUID NOT NULL,
  name TEXT,
  CONSTRAINT account_id_fk FOREIGN KEY(id) REFERENCES account(id),
  INDEX(account_id)
);

-- Region is geographic fault domain.  Quoting Google Cloud Platform:
--
-- https://cloud.google.com/docs/geography-and-regions
--
-- Regions are independent geographic areas that consist of zones. Locations
-- within regions tend to have round-trip network latencies of under 5ms on the
-- 95th percentile.
CREATE TABLE IF NOT EXISTS region (
  id TEXT NOT NULL,
  PRIMARY KEY(id)
);

-- Facility represents a data center or data module.  Customers do not interact
-- directly with facilities, only AZs.  Operators service equipment and networks
-- within a facility.  Quoting Google Cloud Platform:
--
-- https://cloud.google.com/docs/geography-and-regions
--
-- A zone is a deployment area for Cloud Platform resources within a
-- region. Zones should be considered a single failure domain within a
-- region. In order to deploy fault-tolerant applications with high
-- availability, you should deploy your applications across multiple zones in a
-- region to help protect against unexpected failures.
CREATE TABLE IF NOT EXISTS facility (
  id UUID NOT NULL DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  region_id TEXT NOT NULL,
  PRIMARY KEY(region_id, name),
  UNIQUE(name),
  UNIQUE(id),
  CONSTRAINT region_id_fk FOREIGN KEY(region_id) REFERENCES region(id)
) INTERLEAVE IN PARENT region(region_id);

-- Facility Network Trust is used to mark the transport used between two
-- different facilities.  Supported transport methods include:
--
-- * 'plain': Plain-text, unencrypted transport
-- * 'IPsec': IPsec encrypted traffic between individual CNs in each facility
CREATE TABLE IF NOT EXISTS facility_network_transit (
  src_facility_id UUID NOT NULL,
  dst_facility_id UUID NOT NULL,
  transport TEXT NOT NULL CHECK (transport IN('plain','IPsec')),
  PRIMARY KEY(src_facility_id, dst_facility_id),
  UNIQUE(dst_facility_id, src_facility_id),
  CONSTRAINT src_facility_id_fk FOREIGN KEY(src_facility_id) REFERENCES facility(id),
  CONSTRAINT dst_facility_id_fk FOREIGN KEY(dst_facility_id) REFERENCES facility(id)
);

-- Availability Zone ("AZ") is the customer designation of a cloud provider
-- designated "facility" or failure domain.
CREATE TABLE IF NOT EXISTS az (
  id UUID NOT NULL DEFAULT gen_random_uuid(),
  region_id TEXT NOT NULL,
  name STRING(1) NOT NULL CHECK (name IN('a','b','c','d','e','f','g','h','i')),
  PRIMARY KEY(region_id, name),
  UNIQUE(id),
  CONSTRAINT region_id_fk FOREIGN KEY(region_id) REFERENCES region(id),
  UNIQUE(region_id, name)
) INTERLEAVE IN PARENT region(region_id);

-- vni (VXLAN IDs) is a master table of all VNIs in use in a facility.  VNIs
-- are populated on demand for a given facility.  Ideally it would be possible
-- to create predicate indexes:
--
-- * CREATE INDEX vni_available_idx ON vnis (id) WHERE in_use = TRUE;
-- * CREATE INDEX vni_used_idx ON vnis (id) WHERE in_use = FALSE;
--
-- For now, use a more simple data model and suck up the expense of performing a
-- sequential scan to look for an available VNI in a facility.
CREATE TABLE IF NOT EXISTS vni (
  facility_id UUID NOT NULL,
  vni INT NOT NULL CHECK (vni > 0 AND vni < 2 ^ 24),

  -- When vpc_id IS NOT NULL, the vpc_id is the current owner of a given VNI in
  -- an AZ.  When vpc_id IS NULL, this VNI is available for reuse.
  vpc_id UUID,

  -- expired_at is used as a tombstone for a VNI addresses in order to prevent
  -- a VNIs reuse before the state of the system converges.
  expired_at TIMESTAMP WITH TIME ZONE,
  expire_after INTERVAL NOT NULL DEFAULT '90 days',

  PRIMARY KEY(facility_id, vni),
  UNIQUE(vni, facility_id),
  INDEX(vpc_id),
  CONSTRAINT facility_id_fk FOREIGN KEY(facility_id) REFERENCES facility(id),
  CONSTRAINT vpc_id_fk FOREIGN KEY(vpc_id) REFERENCES vpc(id)
);

-- A subnet is a network contained within a VPC.  Customers interact with
-- a subnet within a VPC, not a VXLAN or VLAN.
--
-- TODO(seanc@): convert network into a CIDR data type or INT-backed data type
-- that allows for integer CIDR math in order to provide exclusion constraints
-- and range operations.  TL;DR: PostgreSQL's CIDR data type would be nice to
-- have right about now.
CREATE TABLE IF NOT EXISTS subnet (
  id UUID NOT NULL DEFAULT gen_random_uuid(),
  vpc_id UUID NOT NULL,
  address_type TEXT NOT NULL CHECK (address_type IN('IPv4','IPv6')),
  network TEXT NOT NULL,
  prefix_len INT NOT NULL CHECK (prefix_len > 0 AND ((address_type = 'IPv4' AND prefix_len <= 32) OR (address_type = 'IPv6' AND prefix_len <= 128))),
  CONSTRAINT vpc_id_fk FOREIGN KEY(vpc_id) REFERENCES vpc(id),
  PRIMARY KEY(vpc_id, id),
  UNIQUE(id),
  UNIQUE(vpc_id, network)
) INTERLEAVE IN PARENT vpc(vpc_id);

-- Account MAC is a unique mapping of MAC addresses within an Account.  MAC
-- addresses are unique to an Account, not a VPC, because a VNIC may be moved
-- between VPCs within a given Account.  When an IP moves (i.e. a VNIC is moved
-- to a new device) the old MAC address stays on the old host is converted from
-- an interface to a forwarding interface in the VPC bridge.  The forwarding
-- interface forwards all packets destined to the old MAC address to the new MAC
-- address.  Packets destined to the old MAC address also receive a gratiutious
-- ARP ("GARP") back from the forwarding interface informing the sender of the
-- new MAC address.  Forwarding interfaces self-expire after 4hrs of inactivity
-- and have a hard expiration of 8hrs.  The 4hrs value was chosen by convention
-- as that is a common ARP cache TTL, and the 8hrs by fiat as it's 2x the
-- conventional 4hr TTL.  One consequence of this is we explicitly do not
-- support static ARP entries.  Static ARP will work until a VNIC or IP moves,
-- at which point static arp will break 4-8hrs in the future.
CREATE TABLE IF NOT EXISTS account_mac (
  id UUID NOT NULL DEFAULT gen_random_uuid(),
  account_id UUID NOT NULL,
  mac TEXT NOT NULL,
  vpc_id UUID NOT NULL,
  subnet_id UUID NOT NULL,

  -- expired_at is used as a tombstone for MAC addresses in order to allow for
  -- MAC address.
  expired_at TIMESTAMP WITH TIME ZONE,
  expire_after INTERVAL NOT NULL DEFAULT '90 days',
  PRIMARY KEY(account_id, mac),
  UNIQUE(mac, account_id),
  UNIQUE(id, subnet_id),
  CONSTRAINT account_id_fk FOREIGN KEY(account_id) REFERENCES account(id),
  CONSTRAINT vpc_subnet_id_fk FOREIGN KEY(vpc_id, subnet_id) REFERENCES subnet(vpc_id, id)
);

-- Subnet IP is the whitelist of IPs contained within a given subnet.
CREATE TABLE IF NOT EXISTS subnet_ip (
  id UUID NOT NULL DEFAULT gen_random_uuid(),
  vpc_id UUID NOT NULL,
  subnet_id UUID NOT NULL,
  ip TEXT NOT NULL,
  PRIMARY KEY(vpc_id, subnet_id, ip),
  UNIQUE(subnet_id, ip),
  UNIQUE(vpc_id, ip),
  UNIQUE(id),
  CONSTRAINT vpc_subnet_id_fk FOREIGN KEY(vpc_id, subnet_id) REFERENCES subnet(vpc_id, id)
) INTERLEAVE IN PARENT subnet(vpc_id, subnet_id);

-- Router describes a router instance.  A router instance can be connected to
-- one or more subnets (in the same or different VPCs) via subnet interfaces.
CREATE TABLE IF NOT EXISTS router (
  id UUID NOT NULL DEFAULT gen_random_uuid(),
  vpc_id UUID NOT NULL,

  PRIMARY KEY(id),
  UNIQUE(vpc_id, id),
  CONSTRAINT vpc_id_fk FOREIGN KEY(vpc_id) REFERENCES vpc(id)
);

-- Object Type is a list of all of the supported object types in the schema.  As
-- new features are added, new object types need to be added.
CREATE TABLE IF NOT EXISTS obj_type (
  id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  UNIQUE(name)
);
INSERT INTO obj_type (name) VALUES ('vm');
INSERT INTO obj_type (name) VALUES ('router');

-- VNIC is a virtual NIC.  A VNIC can be assigned to a router, VM, or service
-- within a VPC.  The MAC address for a VNIC is assigned to the interface of the
-- object receiving the VNIC.  A VNIC may also be disconnected from an object.
--
-- NOTE: mac MUST be unique to a given Account.
CREATE TABLE IF NOT EXISTS vnic (
  id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
  account_id UUID NOT NULL,

  -- obj_id is the UUID of the object currently using a VNIC and obj_type is the
  -- type.
  --
  -- NOTE(seanc@): At present we do not have any integrity checks from vnic to
  -- the various object types.  In the future we could use additional null
  -- columns to do that, however we would need to measure the expense of such a
  -- schema first (i.e. vm_id UUID, router_id UUID, CHECK(vm_id IS NOT NULL AND
  -- router_id IS NULL...).
  obj_id UUID,
  obj_type UUID,

  UNIQUE(id, obj_id),

  CONSTRAINT account_id_fk FOREIGN KEY(account_id) REFERENCES account(id),
  CONSTRAINT obj_type_fk FOREIGN KEY(obj_type) REFERENCES obj_type(id)
);

-- VNIC IP maps one or more IPs onto a VNIC.
CREATE TABLE IF NOT EXISTS vnic_ip (
  vnic_id UUID NOT NULL,
  ip_id UUID NOT NULL,
  ip_index INT NOT NULL DEFAULT 0,
  PRIMARY KEY(vnic_id, ip_id),
  UNIQUE(vnic_id, ip_index),
  CONSTRAINT vnic_id_fk FOREIGN KEY(vnic_id) REFERENCES vnic(id),
  CONSTRAINT subnet_ip_id_fk FOREIGN KEY(ip_id) REFERENCES subnet_ip(id)
) INTERLEAVE IN PARENT vnic(vnic_id);

-- Security Groups represent a set of network filters.  Security Groups are
-- applied to VNICs.  A Security Group is evaluated in context and then has its
-- rules pushed out to matching VNICs.  At the edge individual devices can
-- decide what rules match a given VNIC.
CREATE TABLE IF NOT EXISTS security_group (
  id UUID NOT NULL DEFAULT gen_random_uuid(),
  name TEXT,
  description TEXT,
  account_id UUID NOT NULL,
  CONSTRAINT account_id_fk FOREIGN KEY(account_id) REFERENCES account(id),
  PRIMARY KEY(id)
);

-- Security Group Rules.  Security Groups are exclusively permissive and
-- stateful.
CREATE TABLE IF NOT EXISTS security_group_rule (
  id UUID NOT NULL DEFAULT gen_random_uuid(),
  security_group_id UUID NOT NULL,

  direction TEXT NOT NULL DEFAULT 'in' CHECK(direction IN('in','out')),

  -- Protocol Number (See /etc/protocols)
  -- Port range
  -- ICMP type and code
  -- Source or destination addresses (CIDR notation)
  -- Security Group
  -- VPC
  -- Subnet
  -- Facility
  protocol INT CHECK(protocol IS NULL OR protocol >= 0),

  src_port_start INT, -- Reused as the ICMP type when protocol is 1 (ICMP)
  src_port_end INT,   -- Reused as the ICMP code when protocol is 1 (ICMP)
  dst_port_start INT,
  dst_port_end INT,

  src_cidr TEXT,
  dst_cidr TEXT,

  -- Hosts w/ VNICs in a Security Group that are allowed to match this rule.
  src_security_group_id UUID,
  dst_security_group_id UUID,

  -- VPCs matching this rule (i.e. routers that span VPCs) and leak traffic this
  -- way if they have the necessary routes.
  src_vpc_id UUID,
  dst_vpc_id UUID,

  -- Subnets matching this rule
  src_subnet_id UUID,
  dst_subnet_id UUID,

  -- AZs matching this rule
  src_az_id UUID,
  dst_az_id UUID,

  PRIMARY KEY(security_group_id, id),
  CONSTRAINT security_group_id_fk FOREIGN KEY(security_group_id) REFERENCES security_group(id),
  CONSTRAINT src_security_group_id_fk FOREIGN KEY(src_security_group_id) REFERENCES security_group(id),
  CONSTRAINT dst_security_group_id_fk FOREIGN KEY(dst_security_group_id) REFERENCES security_group(id),
  CONSTRAINT src_vpc_id_fk FOREIGN KEY(src_vpc_id) REFERENCES vpc(id),
  CONSTRAINT dst_vpc_id_fk FOREIGN KEY(dst_vpc_id) REFERENCES vpc(id),
  CONSTRAINT src_subnet_id_fk FOREIGN KEY(src_subnet_id) REFERENCES subnet(id),
  CONSTRAINT dst_subnet_id_fk FOREIGN KEY(dst_subnet_id) REFERENCES subnet(id),
  CONSTRAINT src_az_id_fk FOREIGN KEY(src_az_id) REFERENCES az(id),
  CONSTRAINT dst_az_id_fk FOREIGN KEY(dst_az_id) REFERENCES az(id)
) INTERLEAVE IN PARENT security_group(security_group_id);

-- Security Group VNIC maps the available security groups assigned to a given
-- VNIC.
CREATE TABLE IF NOT EXISTS security_group_vnic (
  vnic_id UUID NOT NULL,
  id UUID NOT NULL DEFAULT gen_random_uuid(),
  PRIMARY KEY(vnic_id, id),
  UNIQUE(id, vnic_id)
) INTERLEAVE IN PARENT vnic(vnic_id);

-- Router Subnet Interface maps VNICs to a router object.
CREATE TABLE IF NOT EXISTS router_subnet_interface (
  id UUID NOT NULL DEFAULT gen_random_uuid(),
  router_id UUID NOT NULL,
  subnet_id UUID NOT NULL,
  mac_id UUID NOT NULL,
  vnic_id UUID NOT NULL,
  PRIMARY KEY(router_id, subnet_id),
  UNIQUE(mac_id, subnet_id),
  UNIQUE(id),
  CONSTRAINT router_id_fk FOREIGN KEY(router_id) REFERENCES router(id),
  CONSTRAINT mac_id_fk FOREIGN KEY(mac_id, subnet_id) REFERENCES account_mac(id, subnet_id),
  CONSTRAINT vnic_id_fk FOREIGN KEY(vnic_id) REFERENCES vnic(id)
) INTERLEAVE IN PARENT router(router_id);

-- Router subnet route creates an association between two subnets in a VPC and
-- allows packets to symmetrically flow between two subnets.
CREATE TABLE IF NOT EXISTS router_subnet_route (
  id UUID NOT NULL DEFAULT gen_random_uuid(),
  router_id UUID NOT NULL,
  PRIMARY KEY(router_id, id),
  UNIQUE(id),
  src_subnet_intf_id UUID NOT NULL,
  dst_subnet_intf_id UUID NOT NULL,
  UNIQUE(src_subnet_intf_id, dst_subnet_intf_id),
  UNIQUE(dst_subnet_intf_id, src_subnet_intf_id),
  CONSTRAINT src_subnet_intf_id_fk FOREIGN KEY(src_subnet_intf_id) REFERENCES router_subnet_interface(id),
  CONSTRAINT dst_subnet_intf_id_fk FOREIGN KEY(dst_subnet_intf_id) REFERENCES router_subnet_interface(id)
) INTERLEAVE IN PARENT router(router_id);

-- Subnet VNI VLAN maps a subnet to its VNI and VLAN for a given facility.
CREATE TABLE IF NOT EXISTS subnet_vni_vlan (
  facility_id UUID NOT NULL,
  subnet_id UUID NOT NULL,
  vni INT NOT NULL,
  vlan_id INT NOT NULL CHECK (vlan_id >= 0 AND vlan_id <= 4095),
  PRIMARY KEY(facility_id, vni, vlan_id),
  UNIQUE(subnet_id, facility_id),
  CONSTRAINT vni_id_fk FOREIGN KEY(vni, facility_id) REFERENCES vni(vni, facility_id),
  CONSTRAINT subnet_id_fk FOREIGN KEY(subnet_id) REFERENCES subnet(id)
) INTERLEAVE IN PARENT vni(facility_id, vni);

-- Compute Node ("CN") is the physical server responsible for hosting VMs.
CREATE TABLE IF NOT EXISTS cn (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  facility_id UUID NOT NULL,
  CONSTRAINT facility_id_fk FOREIGN KEY(facility_id) REFERENCES facility(id)
);

-- CN Underlay IP maps the underlay IP(s) assigned to a CN.
CREATE TABLE IF NOT EXISTS cn_underlay_ip (
  cn_id UUID  NOT NULL,
  underlay_ip TEXT NOT NULL,
  PRIMARY KEY(cn_id, underlay_ip),
  UNIQUE(underlay_ip),
  CONSTRAINT cn_id_fk FOREIGN KEY(cn_id) REFERENCES cn(id)
) INTERLEAVE IN PARENT cn(cn_id);

-- VM is an instance
CREATE TABLE IF NOT EXISTS vm (
  id UUID NOT NULL DEFAULT gen_random_uuid(),
  cn_id UUID NOT NULL,
  account_id UUID NOT NULL,
  termination_protection BOOL NOT NULL DEFAULT FALSE,
  vm_type TEXT NOT NULL CHECK (vm_type IN('bhyve','kvm','jail','zone')),
  PRIMARY KEY(cn_id, id),
  UNIQUE(id),
  CONSTRAINT account_id_fk FOREIGN KEY(account_id) REFERENCES account(id),
  CONSTRAINT cn_id_fk FOREIGN KEY(cn_id) REFERENCES cn(id)
) INTERLEAVE IN PARENT cn(cn_id);
