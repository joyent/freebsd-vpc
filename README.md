# `vpc` - FreeBSD Virtual Private Cloud ("VPC")

## Vagrant

A version of VPC can be run under Vagrant using VMWare Fusion or Workstation. This is distributed
across multiple boxes which makes it suitable for some types of development, as well as for creating
disposable environments for integration testing.

To build an environment, run `vagrant up` in the root directory of this repository. The following
is required:

- Vagrant
- VMWare Fusion 10 or VMWare Workstation 14<sup>1</sup>
- Vagrant VMWare Plugin for Fusion or Workstation respectively
- [`vagrant-winnfs`][vagrant-winnfs] if running on a Windows host

## Shell Auto Completion

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
