---
date: 2018-02-28T23:59:59Z
title: "vpc mux connect"
slug: vpc_mux_connect
url: /command/vpc_mux_connect
---
## vpc mux connect

connect a VPC Mux to a VPC EthLink

### Synopsis


connect a VPC Mux to a VPC EthLink

```
vpc mux connect [flags]
```

### Options

```
  -h, --help                  help for connect
  -I, --interface-id string   Specify the VPC Interface ID
  -M, --mux-id string         Specify the VPC Mux ID
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
* [vpc mux](/command/vpc_mux)	 - VPC packet multiplexing configuration

