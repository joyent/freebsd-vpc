---
date: 2018-02-28T23:59:59Z
title: "vpc list"
slug: vpc_list
url: /command/vpc_list
---
## vpc list

list counts of each VPC type

### Synopsis


The list operation of vpc(8) is used to display all VPC objects in the system
and their respective VPC IDs.

```
vpc list [flags]
```

### Examples

```
% vpc list
 TYPE     ID                                    UNIT NAME
 ethlink  5c4acd32-1b8d-11e8-b408-0cc47a6c7d1e  ethlink0
 vmnic    07f95a11-6788-2ae7-c306-ba95cff1db38  vmnic0
 vmnic    a774ba3a-1f77-11e8-8006-0cc47a6c7d1e  vmnic1
 vpcp     0ebf50e1-1f79-11e8-8002-0cc47a6c7d1e  vpcp1
 vpcp     ea58b648-203b-a707-cd02-7a552c8d5295  vpcp2
 vpcp     fd436f9c-1f77-11e8-8002-0cc47a6c7d1e  vpcp0
 vpcsw    da64c3f3-095d-91e5-df01-5aabcfc52468  vpcsw0

   TOTAL                    7
```

### Options

```
  -h, --help              help for list
  -c, --obj-counts        list the number of objects per type
  -t, --obj-type string   List objects of a given type. Valid types: ethlink, hostlink, mgmt, vmnic, vpcmux, vpcnat, vpcp, vpcrtr, vpcsw (default "all")
  -s, --sort-by string    Change the sort order within a given type: id, name (default "id")
```

### Options inherited from parent commands

```
  -F, --log-format string   Specify the log format ("auto", "zerolog", or "human") (default "auto")
  -l, --log-level string    Change the log level being sent to stdout (default "INFO")
      --use-color           Use ASCII colors
  -P, --use-pager           Use a pager to read the output (defaults to $PAGER, less(1), or more(1))
  -Z, --utc                 Display times in UTC
```

### SEE ALSO
* [vpc](/command/vpc)	 - vpc configures and manages VPCs

