# `vpc` - FreeBSD Virtual Private Cloud ("VPC")

## Vagrant

A version of VPC can be run under Vagrant using VMWare Fusion or Workstation. This is distributed
across multiple boxes which makes it suitable for some types of development, as well as for creating
disposable environments for integration testing.

To build an environment, install Vagrant and the VMWare plugin, and run `vagrant up` in the root
directory of this repository.

## Shell Auto Completion

Assuming `shells/bash-completion` has already been installed, add the following
to your `~/.bashrc`:

```sh
[[ $PS1 && -f /usr/local/share/bash-completion/bash_completion.sh ]] && \
    source /usr/local/share/bash-completion/bash_completion.sh
[[ $PS1 && -f `go env GOPATH`/src/github.com/sean-/vpc/docs/bash.d/vpc.sh ]] && \
    source `go env GOPATH`/src/github.com/sean-/vpc/docs/bash.d/vpc.sh
```
