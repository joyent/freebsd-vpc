-- VPC has a few different dimensions:
--
-- * Administrative Plane (customer of the VPC)
-- * Control Plane (customer of the VPC)
-- * Data Plane (users of a deployed service)

-- Organization ("Org") is an administrative construct.  An organization is not
-- referenceable by any control or data plane entity.
CREATE TABLE IF NOT EXISTS org (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid()
);

-- Accounts are the primary unit of granularity for VPC objects.  An Account is
-- created by an Organization.  An Account is the primary object owner for
-- control plane entities.
CREATE TABLE IF NOT EXISTS account (
  id UUID NOT NULL DEFAULT gen_random_uuid(),
  org_id UUID NOT NULL,
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
  transport TEXT NOT NULL CHECK (transport IN ('plain','IPsec')),
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
  name STRING(1) NOT NULL CHECK (name IN ('a','b','c','d','e','f','g','h','i')),
  PRIMARY KEY(region_id, name),
  UNIQUE(id),
  CONSTRAINT region_id_fk FOREIGN KEY(region_id) REFERENCES region(id),
  UNIQUE(region_id, name)
) INTERLEAVE IN PARENT region (region_id);

-- vnis (VXLAN IDs) is a master table of all VNIs in use in a facility.  VNIs
-- are populated on demand for a given facility.  Ideally it would be possible
-- to create predicate indexes:
--
-- * CREATE INDEX vnis_available_idx ON vnis (id) WHERE in_use = TRUE;
-- * CREATE INDEX vnis_used_idx ON vnis (id) WHERE in_use = FALSE;
--
-- For now, use a more simple data model and suck up the expense of performing a
-- sequential scan to look for an available VNI in a facility.
CREATE TABLE IF NOT EXISTS vnis (
  facility_id UUID NOT NULL,
  vni INT NOT NULL CHECK (vni > 0 AND vni < 2 ^ 24),

  -- vpc_id is the VPC using a given VNI in an AZ
  vpc_id UUID,
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

-- VPC MAC is a unique mapping of MAC addresses within a VPC
CREATE TABLE IF NOT EXISTS vpc_macs (
  id UUID NOT NULL DEFAULT gen_random_uuid(),
  vpc_id UUID NOT NULL,
  mac TEXT NOT NULL,
  subnet_id UUID, -- May be NULL once a MAC has been expired from the VPC.

  -- expired_at is used as a tombstone for MAC addresses in order to prevent
  -- their reuse before the state of the system converges.
  expired_at TIMESTAMP WITH TIME ZONE,
  expire_after INTERVAL NOT NULL DEFAULT '90 days',
  PRIMARY KEY(vpc_id, mac),
  UNIQUE(mac, vpc_id),
  UNIQUE(id, subnet_id),
  UNIQUE(subnet_id, mac), -- Redundant constraint, but useful lookup
  CONSTRAINT vpc_id_fk FOREIGN KEY(vpc_id) REFERENCES vpc(id),
  CONSTRAINT subnet_id_fk FOREIGN KEY(subnet_id) REFERENCES subnet(id)
); -- INTERLEAVE IN PARENT subnet(vpc_id, subnet_id);

CREATE TABLE IF NOT EXISTS subnet_ip (
  id UUID NOT NULL DEFAULT gen_random_uuid(),
  vpc_id UUID NOT NULL,
  subnet_id UUID NOT NULL,
  ip TEXT NOT NULL,
  PRIMARY KEY(vpc_id, subnet_id, ip),
  UNIQUE(ip, subnet_id),
  UNIQUE(vpc_id, ip),
  UNIQUE(id),
  CONSTRAINT vpc_id_fk FOREIGN KEY(vpc_id) REFERENCES vpc(id),
  CONSTRAINT subnet_id_fk FOREIGN KEY(subnet_id) REFERENCES subnet(id)
) INTERLEAVE IN PARENT subnet(vpc_id, subnet_id);

-- Router describes a router instance.  A router instance can be connected to
-- one or more subnets via subnet interfaces.
CREATE TABLE IF NOT EXISTS router (
  id UUID NOT NULL DEFAULT gen_random_uuid(),
  vpc_id UUID NOT NULL,
  ip_id UUID NOT NULL,
  PRIMARY KEY(id),
  UNIQUE(vpc_id, id),
  UNIQUE(vpc_id, ip_id),
  CONSTRAINT vpc_id_fk FOREIGN KEY(vpc_id) REFERENCES vpc(id),
  CONSTRAINT ip_id_fk FOREIGN KEY(ip_id) REFERENCES subnet_ip(id)
);

-- Router Subnet Interface connects a Router Interface to a given subnet.
CREATE TABLE IF NOT EXISTS router_subnet_interface (
  id UUID NOT NULL DEFAULT gen_random_uuid(),
  router_id UUID NOT NULL,
  subnet_id UUID NOT NULL,
  ip_id UUID NOT NULL,
  mac_id UUID NOT NULL,
  PRIMARY KEY(router_id, subnet_id),
  UNIQUE(id),
  UNIQUE(ip_id),
  CONSTRAINT router_id_fk FOREIGN KEY(router_id) REFERENCES router(id),
  CONSTRAINT ip_id_fk FOREIGN KEY(ip_id) REFERENCES subnet_ip(id),
  CONSTRAINT mac_id_fk FOREIGN KEY(mac_id, subnet_id) REFERENCES vpc_macs(id, subnet_id)
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
  CONSTRAINT vni_id_fk FOREIGN KEY(vni, facility_id) REFERENCES vnis(vni, facility_id),
  CONSTRAINT subnet_id_fk FOREIGN KEY(subnet_id) REFERENCES subnet(id)
) INTERLEAVE IN PARENT vnis(facility_id, vni);

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
  vpc_id UUID NOT NULL,
  vm_type TEXT NOT NULL CHECK (vm_type IN('bhyve','kvm','jail','zone')),
  PRIMARY KEY(cn_id, id),
  UNIQUE(id),
  CONSTRAINT vpc_id_fk FOREIGN KEY(vpc_id) REFERENCES vpc(id),
  CONSTRAINT cn_id_fk FOREIGN KEY(cn_id) REFERENCES cn(id)
) INTERLEAVE IN PARENT cn(cn_id);

-- VNIC is a virtual NIC assigned to a VM.
--
-- NOTE: mac MUST be unique to a given VPC.
CREATE TABLE IF NOT EXISTS vnic (
  id UUID NOT NULL DEFAULT gen_random_uuid(),
  vm_id UUID NOT NULL,
  subnet_id UUID NOT NULL,
  vpc_id UUID NOT NULL,
  mac_id UUID NOT NULL,
  PRIMARY KEY(vm_id, subnet_id, id),
  UNIQUE(id),

  -- TODO(seanc@): we may relax this in the future if necessary.  We need to add
  -- an ordering or priority attribute.
  UNIQUE(vm_id, subnet_id),
  CONSTRAINT vm_id_fk FOREIGN KEY(vm_id) REFERENCES vm(id),
  CONSTRAINT vpc_id_fk FOREIGN KEY(vpc_id) REFERENCES vpc(id),
  CONSTRAINT mac_id_fk FOREIGN KEY(mac_id, subnet_id) REFERENCES vpc_macs(id, subnet_id)
);

-- VNIC IP maps an IP onto a VNIC.
--
-- TODO(seanc@): Confirm that when an IP on a VNIC is removed the MAC address of
-- the VNIC must also change.  Or the loosing CN needs to install a transient
-- forwarder interface?
CREATE TABLE IF NOT EXISTS vnic_ip (
  vnic_id UUID NOT NULL,
  ip_id UUID NOT NULL,
  subnet_id UUID NOT NULL,
  PRIMARY KEY(vnic_id, subnet_id),
  CONSTRAINT vnic_id_fk FOREIGN KEY(vnic_id) REFERENCES vnic(id),
  CONSTRAINT subnet_id_fk FOREIGN KEY(subnet_id) REFERENCES subnet(id),
  CONSTRAINT subnet_ip_id_fk FOREIGN KEY(ip_id) REFERENCES subnet_ip(id)
);
