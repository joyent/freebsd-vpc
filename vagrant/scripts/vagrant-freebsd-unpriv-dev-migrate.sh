#!/bin/sh

set -o errexit

export GOPATH="/opt/gopath"
export PATH="/opt/gopath/bin:$PATH"

mkdir -p ~/.config/vpc
cat > ~/.config/vpc/vpc.toml <<EOF
[db]
host = "172.27.10.11"
EOF

cd /opt/gopath/src/github.com/sean-/vpc
make get-tools
make
bin/vpc db migrate
