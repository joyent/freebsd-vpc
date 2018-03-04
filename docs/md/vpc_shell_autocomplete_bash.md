---
date: 2018-02-28T23:59:59Z
title: "vpc shell autocomplete bash"
slug: vpc_shell_autocomplete_bash
url: /command/vpc_shell_autocomplete_bash
---
## vpc shell autocomplete bash

Generates and install vpc bash autocompletion script

### Synopsis


Generates a bash autocompletion script for vpc

By default, the file is written directly to /usr/local/share/bash-completion/completions/
for convenience, and the command may need superuser rights, e.g.:

	$ sudo vpc shell autocomplete bash

Add `--bash-autocomplete-dir=/path/to/file`. The default file name
is vpc.sh.

Logout and in again to reload the completion scripts,
or just source them in directly:

	$ . /bash_completion.d

```
vpc shell autocomplete bash [flags]
```

### Options

```
  -d, --dir string   autocompletion directory (default "/usr/local/share/bash-completion/completions/")
  -h, --help         help for bash
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
* [vpc shell autocomplete](/command/vpc_shell_autocomplete)	 - Autocompletion generation

