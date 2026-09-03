#!/bin/sh
# Build and run the same three spike programs inside the linux container.
set -e
cd /w
mkdir -p /tmp/b
cp fixture.go /tmp/b/fixture_src.go
cp a9_regrace.go /tmp/b/a9_src.go
cp a10_probe.go /tmp/b/a10_src.go
cd /tmp/b
export GOFLAGS=-mod=mod GOCACHE=/tmp/gocache
go build -o /tmp/b/fixture fixture_src.go
go build -o /tmp/b/a9      a9_src.go
go build -o /tmp/b/a10     a10_src.go
/tmp/b/a9 /tmp/b/fixture
echo
/tmp/b/a10
