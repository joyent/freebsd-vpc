---
date: 2018-02-28T23:59:59Z
title: "vpc"
slug: vpc
url: /command/vpc
---
## vpc

vpc configures and manages VPCs

### Synopsis


vpc configures and manages VPCs

### Examples

```
# Perform a setup for a VM NIC
$ doas vpc switch create --vni=123 --switch-id=da64c3f3-095d-91e5-df13-5aabcfc52468
$ doas vpc vmnic create --vmnic-id=07f95a11-6788-2ae7-c3ce-ba95cff1db38
$ doas vpc vmnic set --vmnic-id=07f95a11-6788-2ae7-c3ce-ba95cff1db38 --num-queues=2
$ doas vpc switch port add --switch-id=da64c3f3-095d-91e5-df13-5aabcfc52468 --port-id=935cf569-17aa-11e8-a53f-507b9da3d9d0
$ doas vpc switch port connect --port-id=935cf569-17aa-11e8-a53f-507b9da3d9d0 --interface-id=07f95a11-6788-2ae7-c3ce-ba95cff1db38
$ doas vpc switch port add --switch-id=da64c3f3-095d-91e5-df13-5aabcfc52468 --port-id=ea58b648-203b-a707-cdf6-7a552c8d5295 --uplink --l2-name=em0 --ethlink-id=5c4acd32-1b8d-11e8-b4c7-0cc47a6c7d1e

$ vpc vmnic get --vmnic-id=07f95a11-6788-2ae7-c3ce-ba95cff1db38
$ vpc list

# Perform a tear down of the above
$ doas vpc ethlink destroy --ethlink-id=5c4acd32-1b8d-11e8-b4c7-0cc47a6c7d1e
$ doas vpc switch port disconnect --port-id=935cf569-17aa-11e8-a53f-507b9da3d9d0 --interface-id=07f95a11-6788-2ae7-c3ce-ba95cff1db38
$ doas vpc vmnic destroy --vmnic-id=07f95a11-6788-2ae7-c3ce-ba95cff1db38
$ doas vpc switch port remove --switch-id=da64c3f3-095d-91e5-df13-5aabcfc52468 --port-id=935cf569-17aa-11e8-a53f-507b9da3d9d0
$ doas vpc switch destroy --switch-id=da64c3f3-095d-91e5-df13-5aabcfc52468
$ vpc list

```

### Options

```
  -h, --help                help for vpc
  -F, --log-format string   Specify the log format ("auto", "zerolog", or "human") (default "auto")
  -l, --log-level string    Change the log level being sent to stdout (default "INFO")
      --use-color           Use ASCII colors
  -P, --use-pager           Use a pager to read the output (defaults to $PAGER, less(1), or more(1))
  -Z, --utc                 Display times in UTC
```

### SEE ALSO
* [vpc agent](/command/vpc_agent)	 - Run vpc
* [vpc db](/command/vpc_db)	 - Interaction with the VPC database
* [vpc doc](/command/vpc_doc)	 - Documentation for vpc
* [vpc ethlink](/command/vpc_ethlink)	 - VPC EthLink management
* [vpc hostif](/command/vpc_hostif)	 - Host network stack interface
* [vpc interface](/command/vpc_interface)	 - VPC interface management
* [vpc list](/command/vpc_list)	 - list counts of each VPC type
* [vpc shell](/command/vpc_shell)	 - shell commands
* [vpc switch](/command/vpc_switch)	 - VPC switch management
* [vpc version](/command/vpc_version)	 - Version vpc schema
* [vpc vm](/command/vpc_vm)	 - Interaction with the VM agent
* [vpc vmnic](/command/vpc_vmnic)	 - VM network interface management

