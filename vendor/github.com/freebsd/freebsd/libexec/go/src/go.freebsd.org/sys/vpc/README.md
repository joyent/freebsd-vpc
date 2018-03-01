# Virtual Private Cloud ("VPC")

The Virtual Private Cloud ("VPC") system provides a cloud-natural experience for
bhyve guests.

## Requirements

1. `vmmnet(4)` must be loaded: `kldload /boot/modules/vmmnet.ko`
2. Add the VPC library to `GOPATH`:

    ```
% mkdir -p `go env GOPATH`/src
% ln -sf /usr/libexec/go/src/go.freebsd.org `go env GOPATH`/src

# or if developing against `src/`:

% ln -sf /usr/src/libexec/go/src/go.freebsd.org `go env GOPATH`/src

# or if you are developing the `src/` (caution though: `GOPATH` may not be
# passed as an environment variable by `doas(1)` or `sudo(1)` to the root
# and is why the symlink in the default `GOPATH/src` is preferred):

% export GOPATH /usr/src/libexec/go
```

## Testing

Run the VPC tests (once `GOPATH` is set correctly):

```
cd src/libexec/go/src/go.freebsd.org/sys/vpc
doas go test -v ./...
doas go test -bench . -benchtime 15s ./...
```

## Development

### Rapid Pull Loop

```sh
#!/bin/sh --

set -e

cd /usr/src
git pull
cd /usr/src
make -j`sysctl -n hw.ncpu` buildkernel KERNCONF=BHYVE -DNO_KERNELCLEAN -s COPTFLAGS=-O0
doas make installkernel KERNCONF=BHYVE
doas kldunload vmmnet
doas kldload vmmnet
cd /usr/src/libexec/go.freebsd.org/sys/vpc
doas go test ./...
```

### Debugging

- `doas truss -faDH go test ./...`

### Enabling `MemGuard(9)`

Add the following kernel option:

```
options         DEBUG_MEMGUARD
```

and install the new kernel.  Once restarted and up and running, enable
`MemGuard(9)` for `ifnet` or `vmmnet`.

```
$ doas sysctl vm.memguard.desc="ifnet"
$ doas sysctl vm.memguard.desc="vmmnet"
```

NOTE: It is required to enable this `sysctl(8)` *BEFORE* loading `vmmnet.ko`
below.

If a panic occurs during testing, `dump` a core and `reboot`.

#### TODO

- Build a series of scripts to trace execution
