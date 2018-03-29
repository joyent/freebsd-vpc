---
date: 2018-02-28T23:59:59Z
title: "vpc mux listen"
slug: vpc_mux_listen
url: /command/vpc_mux_listen
---
## vpc mux listen

listen address to use when sending/receiving muxed VPC traffic

### Synopsis


listen address to use when sending/receiving muxed VPC traffic

```
vpc mux listen [flags]
```

### Options

```
  -h, --help                 help for listen
      --listen-addr string   Address and port the VPC Mux will use to listen for traffic on the underlay network
  -M, --mux-id string        Specify the VPC Mux ID
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

