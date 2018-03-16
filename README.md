# `vpc` - FreeBSD Virtual Private Cloud ("VPC")

## Vagrant

A version of VPC can be run under Vagrant using VMWare Fusion or Workstation. This is distributed
across multiple boxes which makes it suitable for some types of development, as well as for creating
disposable environments for integration testing. The development environment consists of a
compilation machine, a three-node CockroachDB cluster and (optionally) compute nodes.

The CockroachDB cluster uses the TLS certificates in `vagrant/certs` both for client authentication
and encryption of intra-node traffic.

To build an environment consisting of a three-node CockroachDB cluster and a box with the compiler
and relevant development tools installed, run the following command in the root directory of the
repository:

```
vagrant up crdb1 crdb2 crdb3 compile
```

The following dependencies must be installed on the host:

- Vagrant
- VMWare Fusion 10 or VMWare Workstation 14<sup>1</sup>
- Vagrant VMWare Plugin for Fusion or Workstation respectively
- [`vagrant-winnfs`][vagrant-winnfs] if running on a Windows host

## `vpc` Command Shell Auto Completion

Assuming `shells/bash-completion` has already been installed, add the following
to your `~/.bashrc`:

```sh
[[ $PS1 && -f /usr/local/share/bash-completion/bash_completion.sh ]] && \
    source /usr/local/share/bash-completion/bash_completion.sh
[[ $PS1 && -f `go env GOPATH`/src/github.com/sean-/vpc/docs/bash.d/vpc.sh ]] && \
    source `go env GOPATH`/src/github.com/sean-/vpc/docs/bash.d/vpc.sh
```

## Footnotes

[1]: It is possible to build boxes for Virtualbox, but since nested virtualization is not supported
     under Virtualbox, it is not possible to test running instances in virtual machines.

[vagrant-winnfs]: https://github.com/winnfsd/vagrant-winnfsd
