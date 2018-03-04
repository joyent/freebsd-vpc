---
date: 2018-02-28T23:59:59Z
title: "vpc switch port connect"
slug: vpc_switch_port_connect
url: /command/vpc_switch_port_connect
---
## vpc switch port connect

connect a VPC Interface to a VPC Switch Port

### Synopsis


connect a VPC Interface to a VPC Switch Port

```
vpc switch port connect [flags]
```

### Options

```
  -h, --help                  help for connect
  -I, --interface-id string   Specify the VPC Interface ID
      --port-id string        Specify the VPC Port ID
```

### Options inherited from parent commands

```
  -F, --log-format string   Specify the log format ("auto", "zerolog", or "human") (default "auto")
  -l, --log-level string    Change the log level being sent to stdout (default "INFO")
      --use-color           Use ASCII colors (default true)
  -P, --use-pager           Use a pager to read the output (defaults to $PAGER, less(1), or more(1)) (default true)
  -Z, --utc                 Display times in UTC
```

### SEE ALSO
* [vpc switch port](/command/vpc_switch_port)	 - VPC switch management

