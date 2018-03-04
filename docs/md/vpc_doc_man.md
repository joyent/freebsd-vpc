---
date: 2018-02-28T23:59:59Z
title: "vpc doc man"
slug: vpc_doc_man
url: /command/vpc_doc_man
---
## vpc doc man

Generates and install vpc man(1) pages

### Synopsis


This command automatically generates up-to-date man(1) pages of vpc(8)
command-line interface.  By default, it creates the man page files
in the "docs/man" directory under the current directory.

```
vpc doc man [flags]
```

### Options

```
  -h, --help             help for man
  -m, --man-dir string   Specify the MANDIR to use (default "./docs/man")
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
* [vpc doc](/command/vpc_doc)	 - Documentation for vpc

