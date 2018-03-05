---
date: 2018-02-28T23:59:59Z
title: "vpc db"
slug: vpc_db
url: /command/vpc_db
---
## vpc db

Interaction with the VPC database

### Synopsis


Interaction with the VPC database

### Options

```
      --db-host string       Database server address (default "127.0.0.1")
      --db-name string       Database name (default "triton")
      --db-password string   Database password (default "tls")
      --db-port uint         Database port (default 26257)
      --db-username string   Database username (default "root")
  -h, --help                 help for db
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
* [vpc db migrate](/command/vpc_db_migrate)	 - Migrate vpc schema
* [vpc db ping](/command/vpc_db_ping)	 - ping the database to ensure connectivity

