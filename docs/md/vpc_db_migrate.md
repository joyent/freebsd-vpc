---
date: 2018-02-28T23:59:59Z
title: "vpc db migrate"
slug: vpc_db_migrate
url: /command/vpc_db_migrate
---
## vpc db migrate

Migrate vpc schema

### Synopsis


Migrate vpc schema

```
vpc db migrate [flags]
```

### Options

```
  -h, --help   help for migrate
```

### Options inherited from parent commands

```
      --db-host string       Database server address (default "127.0.0.1")
      --db-name string       Database name (default "triton")
      --db-password string   Database password (default "tls")
      --db-port uint         Database port (default 26257)
      --db-username string   Database username (default "root")
  -F, --log-format string    Specify the log format ("auto", "zerolog", or "human") (default "auto")
  -l, --log-level string     Change the log level being sent to stdout (default "INFO")
      --use-color            Use ASCII colors (default true)
  -P, --use-pager            Use a pager to read the output (defaults to $PAGER, less(1), or more(1)) (default true)
  -Z, --utc                  Display times in UTC
```

### SEE ALSO
* [vpc db](/command/vpc_db)	 - Interaction with the VPC database

