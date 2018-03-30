---
date: 2018-02-28T23:59:59Z
title: "vpc switch port set"
slug: vpc_switch_port_set
url: /command/vpc_switch_port_set
---
## vpc switch port set

set VPC Port Information

### Synopsis


set VPC Port Information

```
vpc switch port set [flags]
```

### Options

```
  -h, --help             help for set
      --port-id string   Specify the VPC Port ID
  -n, --vni int          set the VNI of a given VPC Port (default -1)
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
* [vpc switch port](/command/vpc_switch_port)	 - VPC switch port management

