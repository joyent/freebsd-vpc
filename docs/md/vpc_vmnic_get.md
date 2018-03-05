---
date: 2018-02-28T23:59:59Z
title: "vpc vmnic get"
slug: vpc_vmnic_get
url: /command/vpc_vmnic_get
---
## vpc vmnic get

get VMNIC information

### Synopsis


get VMNIC information

```
vpc vmnic get [flags]
```

### Options

```
  -h, --help              help for get
  -n, --num-queues        get the number of hardware queues for a given VM NIC (default true)
  -N, --vmnic-id string   Specify the VM NIC ID
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
* [vpc vmnic](/command/vpc_vmnic)	 - VM network interface management

