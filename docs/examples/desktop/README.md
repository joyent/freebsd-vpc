# Desktop VPC

FreeBSD/VPC can, in theory, be used in place of the host network stack by making
use of a `hostlink(4)` interface.  See `setup.sh` for an example of how this
works.  The only thing not done in this example is to run `dhclient(8)` on
`hostlink0` or set a static IP on `hostlink0` using `ifconfig(8)`.
