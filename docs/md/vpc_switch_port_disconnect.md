---
date: 2018-02-28T23:59:59Z
title: "vpc switch port disconnect"
slug: vpc_switch_port_disconnect
url: /command/vpc_switch_port_disconnect
---
## vpc switch port disconnect

disconnect a VPC Interface from a VPC Switch Port

### Synopsis


disconnect a VPC Interface from a VPC Switch Port

```
vpc switch port disconnect [flags]
```

### Options

```
  -h, --help                  help for disconnect
  -I, --interface-id string   Specify the VPC Interface ID
      --port-id string        Specify the VPC Port ID
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

