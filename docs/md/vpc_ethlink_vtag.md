---
date: 2018-02-28T23:59:59Z
title: "vpc ethlink vtag"
slug: vpc_ethlink_vtag
url: /command/vpc_ethlink_vtag
---
## vpc ethlink vtag

Get or set the VTag on a VPC EthLink

### Synopsis


Get or set the VTag on a VPC EthLink

```
vpc ethlink vtag [flags]
```

### Options

```
  -E, --ethlink-id string   Specify the EthLink ID
  -g, --get-vtag            get the VTag for a given VPC EthLink NIC (default true)
  -h, --help                help for vtag
  -s, --set-vtag int        set the VTag for a given VPC EthLink NIC (default -1)
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

