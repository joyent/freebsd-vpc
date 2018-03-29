---
date: 2018-02-28T23:59:59Z
title: "vpc ethlink connect"
slug: vpc_ethlink_connect
url: /command/vpc_ethlink_connect
---
## vpc ethlink connect

connect a VPC EthLink interface to a physical or cloned interface

### Synopsis


`vpc ethlink connect` is used to create a VPC interface that wraps a cloned or physical interface.  The cloned or physical interface is typically the interface used for the underlay network to route between different hosts.

```
vpc ethlink connect [flags]
```

### Options

```
  -E, --ethlink-id string   Specify the EthLink ID
  -h, --help                help for connect
  -n, --l2-name string      Name of the underlay L2 interface to be wrapped by VPC EthLink
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
* [vpc ethlink](/command/vpc_ethlink)	 - VPC EthLink management

